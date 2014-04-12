package main

import (
  "flag"
  "fmt"
  //"github.com/mitchellh/goamz/aws"
  //"github.com/mitchellh/goamz/ec2"
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

func setup(port string) {
  mux := http.NewServeMux()
  mux.HandleFunc("/accounts/aws/credentials", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, credData)
  })
  mux.HandleFunc("/accounts/aws/providers/ec2/arguments", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, ec2ArgData)
  })
  mux.HandleFunc("/accounts/aws/providers/ec2/collections/instance/state", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, ec2InstanceStateData)
  })
  log.Fatal(http.ListenAndServe(":"+port, mux))
}

func main() {
  port := flag.String("port", "32786", "port to listen on")
  flag.Parse()
  setup(*port)
}
