package main

import (
  "bytes"
  "encoding/json"
  "flag"
  "fmt"
  "log"
  "net/http"

  "github.com/mitchellh/goamz/aws"
  "github.com/mitchellh/goamz/ec2"
)

var server *http.Server

var credData string = `["access_key_id", "secret_access_key"]`
var ec2ArgData string = `["region"]`

type stateRequest struct {
  Credentials map[string]string `json:"credentials"`
  Arguments   map[string]string `json:"arguments"`
  ParentIds   []string          `json:"parent_ids"`
}

func main() {
  port := flag.String("port", "32786", "port to listen on")
  flag.Parse()
  setup(*port)
}

func setup(port string) {
  mux := http.NewServeMux()

  mux.HandleFunc("/accounts/aws/credentials", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, credData)
  })
  mux.HandleFunc("/accounts/aws/providers/ec2/arguments", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, ec2ArgData)
  })
  mux.HandleFunc("/accounts/aws/providers/ec2/collections/instance/state", func(w http.ResponseWriter, r *http.Request) {
    body := r.Body
    defer r.Body.Close()

    buf := new(bytes.Buffer)
    buf.ReadFrom(body)

    sr := stateRequest{}
    err := json.Unmarshal(buf.Bytes(), &sr)
    if err != nil {
      http.Error(w, err.Error(), 400)
    }
    resp, err := processRequest(sr)
    if err != nil {
      http.Error(w, err.Error(), 400)
    }
    fmt.Fprintf(w, resp)
  })
  log.Fatal(http.ListenAndServe(":"+port, mux))
}

func processRequest(sr stateRequest) (string, error) {
  accessKey := sr.Credentials["access_key_id"]
  secretKey := sr.Credentials["secret_access_key"]
  region := sr.Arguments["region"]

  if accessKey == "" || secretKey == "" {
    return "", fmt.Errorf("access_key_id or secret_access_key not specified")
  }
  auth := aws.Auth{accessKey, secretKey, ""}

  if _, ok := aws.Regions[region]; !ok {
    return "", fmt.Errorf("invalid or unspecified region")
  }
  awsRegion := aws.Regions[region]
  ec2Conn := ec2.New(auth, awsRegion)

  instances, err := ec2Conn.Instances([]string{}, nil)

  res := []stateResponse{}

  for _, reservation := range instances.Reservations {
    for _, i := range reservation.Instances {
      sResp := stateResponse{Id: i.InstanceId}
      res = append(res, sResp)
    }
  }

  if err != nil {
    return "", err
  }

  data, err := json.Marshal(res)
  if err != nil {
    return "", err
  }
  return string(data), err
}

type stateResponse struct {
  Id   string       `json:"id"`
  data instanceData `json:"data"`
}

type instanceData struct {
  ImageId []string `json:"image_id"`
}
