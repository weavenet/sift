package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/mitchellh/goamz/ec2"
)

type ec2SecurityGroupStateResponse struct {
  Id   string            `json:"id"`
  Data securityGroupData `json:"data"`
}

type securityGroupData struct {
  VpcId []string `json:"vpc_id"`
}

func newEc2SecurityGroupStateResponse(id string) ec2SecurityGroupStateResponse {
  d := securityGroupData{}
  return ec2SecurityGroupStateResponse{Id: id, Data: d}
}

func serveAwsEc2SecurityGroup(mux *http.ServeMux) {
  mux.HandleFunc("/accounts/aws/providers/ec2/collections/securitygroup/state", func(w http.ResponseWriter, r *http.Request) {
    body := r.Body
    defer r.Body.Close()

    buf := new(bytes.Buffer)
    buf.ReadFrom(body)

    sr := newStateRequest()
    err := json.Unmarshal(buf.Bytes(), &sr)
    if err != nil {
      http.Error(w, newErrorResponse(err).String(), 400)
    }
    resp, err := processEc2SecurityGroupRequest(sr)
    if err != nil {
      http.Error(w, newErrorResponse(err).String(), 400)
    }
    fmt.Fprintf(w, resp)
  })
}

func processEc2SecurityGroupRequest(sr stateRequest) (string, error) {
  ec2Conn, err := connectEc2(sr)
  if err != nil {
    return "", err
  }

  securityGroups, err := ec2Conn.SecurityGroups([]ec2.SecurityGroup{}, nil)
  if err != nil {
    return "", err
  }

  res := []ec2SecurityGroupStateResponse{}
  for _, sg := range securityGroups.Groups {
    sr := newEc2SecurityGroupStateResponse(sg.Id)
    sr.Data.VpcId = []string{sg.VpcId}
    res = append(res, sr)
  }

  data, err := json.Marshal(res)
  if err != nil {
    return "", err
  }
  return string(data), nil
}
