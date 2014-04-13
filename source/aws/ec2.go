package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "log"
  "net/http"

  "github.com/mitchellh/goamz/aws"
  "github.com/mitchellh/goamz/ec2"
)

var ec2ArgData string = `["region"]`

func serveEc2(port string) {
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
    resp, err := processInstanceEc2Request(sr)
    if err != nil {
      http.Error(w, err.Error(), 400)
    }
    fmt.Fprintf(w, resp)
  })
  log.Fatal(http.ListenAndServe(":"+port, mux))
}

func processInstanceEc2Request(sr stateRequest) (string, error) {
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

  res := []ec2InstanceStateResponse{}

  for _, reservation := range instances.Reservations {
    for _, i := range reservation.Instances {
      sr := newEc2InstanceStateResponse(i.InstanceId)
      sr.Data.ImageId = []string{i.ImageId}
      res = append(res, sr)
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

func newEc2InstanceStateResponse(id string) ec2InstanceStateResponse {
  d := instanceData{}
  return ec2InstanceStateResponse{Id: id, Data: d}
}

type ec2InstanceStateResponse struct {
  Id   string       `json:"id"`
  Data instanceData `json:"data"`
}

type instanceData struct {
  ImageId []string `json:"image_id"`
}
