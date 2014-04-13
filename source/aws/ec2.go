package main

import (
  "fmt"
  "net/http"
)

var ec2ArgData string = `["region"]`

func serveAwsEc2(mux *http.ServeMux) {
  mux.HandleFunc("/accounts/aws/providers/ec2/arguments", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, ec2ArgData)
  })
}
