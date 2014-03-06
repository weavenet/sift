package source

import (
  "bytes"
  "encoding/json"
  "fmt"
  log "github.com/cihub/seelog"
  "github.com/siftproject/sift/sift/state"
  "io/ioutil"
  "net/http"
)

type Source struct {
  accountName    string
  providerName   string
  collectionName string
  url            string
}

func NewSource(accountName string, providerName string, collectionName string, url string) *Source {
  s := Source{
    accountName:    accountName,
    providerName:   providerName,
    collectionName: collectionName,
    url:            url,
  }
  return &s
}

func (s *Source) SetURL(u string) {
  s.url = u
}

func (s Source) Credentials() (creds []string, err error) {
  log.Debugf("Loading required credentials for '%s' '%s' %s'.", s.accountName, s.providerName, s.collectionName)
  log.Tracef("Connecting to '%s'.", s.credentialsURL())
  data, err := s.query(s.credentialsURL(), "get", []byte{})
  if err != nil {
    return creds, err
  }
  if err := json.Unmarshal(data, &creds); err != nil {
    return creds, err
  }
  return creds, nil
}

func (s Source) ProviderArguments() (args []string, err error) {
  log.Debugf("Loading required arguments for '%s' '%s' %s'.", s.accountName, s.providerName, s.collectionName)
  log.Tracef("Connecting to '%s'.", s.providerArgumentsURL())
  data, err := s.query(s.providerArgumentsURL(), "get", []byte{})
  if err != nil {
    return args, err
  }
  if err := json.Unmarshal(data, &args); err != nil {
    return args, err
  }
  return args, nil
}

func (s Source) State(creds map[string]string, args map[string]string, parentIds []string) (states []state.State, err error) {
  log.Tracef("Connecting to '%s'.", s.stateURL())
  requestBody := newStateRequestBody(creds, args, parentIds)
  requestBodyJson, err := json.Marshal(requestBody)
  if err != nil {
    return states, err
  }
  log.Tracef("Posting JSON '%s'.", string(requestBodyJson))
  data, err := s.query(s.stateURL(), "post", requestBodyJson)
  if err != nil {
    return states, err
  }
  if err := json.Unmarshal(data, &states); err != nil {
    return states, err
  }
  log.Tracef("Loaded '%d' states from source.", len(states))

  return states, nil
}

func (s Source) baseURL() string {
  return s.url
}

func (s Source) credentialsURL() string {
  return s.baseURL() + "/accounts/" + s.accountName + "/credentials"
}

func (s Source) providerArgumentsURL() string {
  return s.baseURL() + "/accounts/" + s.accountName + "/providers/" + s.providerName + "/arguments"
}

func (s Source) stateURL() string {
  return s.baseURL() + "/accounts/" + s.accountName + "/providers/" + s.providerName + "/collections/" + s.collectionName + "/state"
}

func (s Source) query(url string, method string, body []byte) (data []byte, err error) {
  var resp *http.Response
  switch method {
  case "get":
    {
      resp, err = http.Get(url)
    }
  case "post":
    {
      resp, err = http.Post(url, "applicaiton/json", bytes.NewReader(body))
    }
  }
  if err != nil {
    return data, err
  }
  log.Tracef("Recieved response from external source '%+v'.", resp)

  if resp.StatusCode != 200 {
    log.Debugf("Received error code '%s' connecting to '%s'.", resp.Status, s.stateURL())
    return data, fmt.Errorf("Unable to connect to external source '%s'.", s.stateURL())
  }
  defer resp.Body.Close()

  data, err = ioutil.ReadAll(resp.Body)
  if err != nil {
    return data, err
  }
  log.Tracef("Received body '%s'", data)
  return data, err
}
