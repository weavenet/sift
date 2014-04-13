package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/mitchellh/goamz/aws"
  "github.com/mitchellh/goamz/ec2"
)

func serveAwsEc2SecurityGroup(mux *http.ServeMux) {
  mux.HandleFunc("/accounts/aws/providers/ec2/collections/securitygroup/state", func(w http.ResponseWriter, r *http.Request) {
    body := r.Body
    defer r.Body.Close()

    buf := new(bytes.Buffer)
    buf.ReadFrom(body)

    sr := stateRequest{}
    err := json.Unmarshal(buf.Bytes(), &sr)
    if err != nil {
      http.Error(w, err.Error(), 400)
    }
    resp, err := processEc2SecurityGroupRequest(sr)
    if err != nil {
      http.Error(w, err.Error(), 400)
    }
    fmt.Fprintf(w, resp)
  })
}

func processEc2SecurityGroupRequest(sr stateRequest) (string, error) {
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

  securityGroups, err := ec2Conn.SecurityGroups([]ec2.SecurityGroup{}, nil)

  res := []ec2SecurityGroupStateResponse{}

  for _, sg := range securityGroups.Groups {
    sr := newEc2SecurityGroupStateResponse(sg.Id)
    sr.Data.VpcId = sg.VpcId
    res = append(res, sr)
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

func newEc2SecurityGroupStateResponse(id string) ec2SecurityGroupStateResponse {
  d := securityGroupData{}
  return ec2SecurityGroupStateResponse{Id: id, Data: d}
}

type ec2SecurityGroupStateResponse struct {
  Id   string            `json:"id"`
  Data securityGroupData `json:"data"`
}

type securityGroupData struct {
  VpcId string `json:"vpc_id"`
}
