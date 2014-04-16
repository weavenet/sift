package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"
  "strconv"

  "github.com/mitchellh/goamz/ec2"
)

type ec2SecurityGroupIpPermissionStateResponse struct {
  Id        string                        `json:"id"`
  PartentId string                        `json:"parent_id"`
  Data      securityGroupIpPermissionData `json:"data"`
}

type securityGroupIpPermissionData struct {
  Protocol  []string `json:"protocol"`
  FromPort  []string `json:"from_port"`
  ToPort    []string `json:"to_port"`
  SourceIPs []string `json:"source_ips"`
}

func newEc2SecurityGroupIpPermissionStateResponse(id string, pid string, d securityGroupIpPermissionData) ec2SecurityGroupIpPermissionStateResponse {
  return ec2SecurityGroupIpPermissionStateResponse{Id: id, ParentId: pid, Data: d}
}

func newEc2SecurityGroupIpPermissionStateResponseData(ipperm ec2.IPPerm) securityGroupIpPermissionData {
  d := securityGroupIpPermissionData{
    Protocol:  []string{ipperm.Protocol},
    FromPort:  []string{strconv.Itoa(ipperm.FromPort)},
    ToPort:    []string{strconv.Itoa(ipperm.ToPort)},
    SourceIPs: ipperm.SourceIPs,
  }
  return d
}

func serveAwsEc2SecurityGroupIpPermission(mux *http.ServeMux) {
  mux.HandleFunc("/accounts/aws/providers/ec2/collections/securitygroup-ippermission/state", func(w http.ResponseWriter, r *http.Request) {
    body := r.Body
    defer r.Body.Close()

    buf := new(bytes.Buffer)
    buf.ReadFrom(body)

    sr := newStateRequest()
    err := json.Unmarshal(buf.Bytes(), &sr)
    if err != nil {
      http.Error(w, newErrorResponse(err).String(), 400)
    }
    resp, err := processEc2SecurityGroupIpPermissionRequest(sr)
    if err != nil {
      http.Error(w, newErrorResponse(err).String(), 400)
    }
    fmt.Fprintf(w, resp)
  })
}

func processEc2SecurityGroupIpPermissionRequest(sr stateRequest) (string, error) {
  ec2Conn, err := connectEc2(sr)
  if err != nil {
    return "", err
  }

  sgs := []ec2.SecurityGroup{}
  for _, id := range sr.ParentIds {
    sg := ec2.SecurityGroup{Id: id}
    sgs = append(sgs, sg)
  }

  securityGroups, err := ec2Conn.SecurityGroups(sgs, nil)
  if err != nil {
    return "", err
  }

  res := []ec2SecurityGroupIpPermissionStateResponse{}

  for _, sg := range securityGroups.Groups {
    for count, ipperm := range sg.IPPerms {
      data := newEc2SecurityGroupIpPermissionStateResponseData(ipperm)
      sr := newEc2SecurityGroupIpPermissionStateResponse(count, sg.Id, data)
      res = append(res, sr)
    }
  }

  data, err := json.Marshal(res)
  if err != nil {
    return "", err
  }
  return string(data), nil
}
