package source

import (
  "encoding/json"
  "github.com/brettweavnet/sift/sift/state"
  "io"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "reflect"
  "testing"
)

var credTCs = []struct {
  path, body, account, provider, resource string
  credentials                             []string
  e                                       error
}{
  {"/accounts/aws/credentials", `["region"]`, "aws", "ec2", "stub", []string{"region"}, nil},
}

func TestCredentials(t *testing.T) {
  for _, tc := range credTCs {
    handler := func(w http.ResponseWriter, r *http.Request) {
      if r.URL.Path == tc.path {
        io.WriteString(w, tc.body)
      }
    }
    server := httptest.NewServer(http.HandlerFunc(handler))
    defer server.Close()

    s := NewSource(tc.account, tc.provider, tc.resource, server.URL)

    r, err := s.Credentials()
    if err != nil {
      t.Errorf("Received error loading credentials '%s'.", err)
    }

    if !reflect.DeepEqual(tc.credentials, r) {
      t.Errorf("Incorrect credentials loaded. Expect '%s' Received '%s'.", tc.credentials, r)
    }
  }
}

var argsTCs = []struct {
  path, body, account, provider, resource string
  arguments                               []string
  e                                       error
}{
  {"/accounts/aws/providers/ec2/arguments", `["arg1","arg2"]`, "aws", "ec2", "stub", []string{"arg1", "arg2"}, nil},
}

func TestArguments(t *testing.T) {
  for _, tc := range argsTCs {
    handler := func(w http.ResponseWriter, r *http.Request) {
      if r.URL.Path == tc.path {
        io.WriteString(w, tc.body)
      }
    }
    server := httptest.NewServer(http.HandlerFunc(handler))
    defer server.Close()

    s := NewSource(tc.account, tc.provider, tc.resource, server.URL)

    r, err := s.ProviderArguments()
    if err != nil {
      t.Errorf("Received error loading arguments '%s'.", err)
    }

    if !reflect.DeepEqual(tc.arguments, r) {
      t.Errorf("Incorrect arguments loaded. Expect '%s' Received '%s'.", tc.arguments, r)
    }
  }
}

var testState1 = state.State{
  Id:   "1",
  Data: map[string][]string{"attr_id": []string{"a1", "a2"}},
}
var testState2 = state.State{
  Id:   "2",
  Data: map[string][]string{"attr_id": []string{"a2", "a3"}},
}

var stateTCs = []struct {
  path, body, account, provider, resource string
  args                                    map[string]string
  creds                                   map[string]string
  requestBody                             stateRequestBody
  states                                  []state.State
  e                                       error
}{
  {
    "/accounts/aws/providers/ec2/collections/stub/state",
    `[{"id":"1", "data":{"attr_id":["a1","a2"]}},{"id":"2","data":{"attr_id":["a2","a3"]}}]`,
    "aws",
    "ec2",
    "stub",
    map[string]string{"region": "us-west-1"},
    map[string]string{"key": "123", "secret": "321"},
    newStateRequestBody(map[string]string{"key": "123", "secret": "321"},
      map[string]string{"region": "us-west-1"},
      []string{"i-1"}),
    []state.State{testState1, testState2},
    nil,
  },
}

func TestState(t *testing.T) {
  for _, tc := range stateTCs {
    handler := func(w http.ResponseWriter, r *http.Request) {
      if r.URL.Path == tc.path {
        requestBody := stateRequestBody{}
        body, _ := ioutil.ReadAll(r.Body)
        t.Logf("%s", body)
        if err := json.Unmarshal(body, &requestBody); err != nil {
          t.Fatalf("Error reading body.")
        }
        if !reflect.DeepEqual(requestBody, tc.requestBody) {
          t.Fatalf("Incorrect body received. Expected '%+v' Received '%+v'.", tc.requestBody, requestBody)
        }
        io.WriteString(w, tc.body)
      }
    }
    server := httptest.NewServer(http.HandlerFunc(handler))
    defer server.Close()

    s := NewSource(tc.account, tc.provider, tc.resource, server.URL)

    parentIds := []string{"i-1"}
    r, err := s.State(tc.creds, tc.args, parentIds)
    if err != nil {
      t.Errorf("Received error loading arguments '%s'.", err)
    }

    if !reflect.DeepEqual(tc.states, r) {
      t.Errorf("Incorrect state loaded. Expect '%s' Received '%s'.", tc.states, r)
    }
  }
}
