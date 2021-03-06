package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"
)

type ec2InstanceStateResponse struct {
  Id   string       `json:"id"`
  Data instanceData `json:"data"`
}

type instanceData struct {
  ImageId []string `json:"image_id"`
}

func newEc2InstanceStateResponse(id string) ec2InstanceStateResponse {
  d := instanceData{}
  return ec2InstanceStateResponse{Id: id, Data: d}
}

func serveAwsEc2Instance(mux *http.ServeMux) {
  mux.HandleFunc("/accounts/aws/providers/ec2/collections/instance/state", func(w http.ResponseWriter, r *http.Request) {
    body := r.Body
    defer r.Body.Close()

    buf := new(bytes.Buffer)
    buf.ReadFrom(body)

    sr := newStateRequest()
    err := json.Unmarshal(buf.Bytes(), &sr)
    if err != nil {
      http.Error(w, newErrorResponse(err).String(), 400)
    }
    resp, err := processEc2InstanceRequest(sr)
    if err != nil {
      http.Error(w, newErrorResponse(err).String(), 400)
    }
    fmt.Fprintf(w, resp)
  })
}

func processEc2InstanceRequest(sr stateRequest) (string, error) {
  ec2Conn, err := connectEc2(sr)
  if err != nil {
    return "", err
  }

  instances, err := ec2Conn.Instances([]string{}, nil)
  if err != nil {
    return "", err
  }

  res := []ec2InstanceStateResponse{}
  for _, reservation := range instances.Reservations {
    for _, i := range reservation.Instances {
      sr := newEc2InstanceStateResponse(i.InstanceId)
      sr.Data.ImageId = []string{i.ImageId}
      res = append(res, sr)
    }
  }

  data, err := json.Marshal(res)
  if err != nil {
    return "", err
  }
  return string(data), nil
}
