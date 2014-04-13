package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
)

var server *http.Server

var credData string = `["access_key_id", "secret_access_key"]`

type stateRequest struct {
  Credentials map[string]string `json:"credentials"`
  Arguments   map[string]string `json:"arguments"`
  ParentIds   []string          `json:"parent_ids"`
}

func main() {
  port := flag.String("port", "32786", "port to listen on")
  mux := http.NewServeMux()
  flag.Parse()
  serveAws(mux)
  serveAwsEc2(mux)
  serveAwsEc2Instance(mux)
  serveAwsEc2SecurityGroup(mux)
  log.Fatal(http.ListenAndServe(":"+*port, mux))
}

func serveAws(mux *http.ServeMux) {
  mux.HandleFunc("/accounts/aws/credentials", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, credData)
  })
}
