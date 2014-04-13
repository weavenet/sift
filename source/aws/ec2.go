package main

import (
  "fmt"
  "github.com/mitchellh/goamz/aws"
  "github.com/mitchellh/goamz/ec2"
  "net/http"
)

var ec2ArgData string = `["region"]`

func serveAwsEc2(mux *http.ServeMux) {
  mux.HandleFunc("/accounts/aws/providers/ec2/arguments", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, ec2ArgData)
  })
}

func connectEc2(sr stateRequest) (*ec2.EC2, error) {
  accessKey := sr.Credentials["access_key_id"]
  secretKey := sr.Credentials["secret_access_key"]
  region := sr.Arguments["region"]

  if accessKey == "" || secretKey == "" {
    return nil, fmt.Errorf("access_key_id or secret_access_key not specified")
  }
  auth := aws.Auth{accessKey, secretKey, ""}

  if _, ok := aws.Regions[region]; !ok {
    return nil, fmt.Errorf("invalid or unspecified region")
  }
  awsRegion := aws.Regions[region]

  return ec2.New(auth, awsRegion), nil
}
