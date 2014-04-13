package main

import (
  "bytes"
  "encoding/json"
  "flag"
  "fmt"
  "github.com/mitchellh/goamz/aws"
  "log"
  "net/http"
)

var server *http.Server

var credData string = `["access_key_id", "secret_access_key"]`
var ec2ArgData string = `["region"]`
var ec2InstanceStateData string = `
[
  {
    "id" : "i-12345678",
    "data" : { "image_id" : ["ami-12345678"] }
  },
  {
    "id" : "i-9876abcd",
    "data" : { "image_id" : ["ami-87654321"] }
  }
]
`

type stateRequest struct {
  Credentials map[string]string `json:"credentials"`
  Arguments   map[string]string `json:"arguments"`
  ParentIds   []string          `json:"parent_ids"`
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

func main() {
  port := flag.String("port", "32786", "port to listen on")
  flag.Parse()
  setup(*port)
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
  return fmt.Sprintf("%s %s", auth, awsRegion), nil
}
