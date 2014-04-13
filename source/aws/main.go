package main

import (
  "flag"
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
  serveEc2Instance(mux)
  log.Fatal(http.ListenAndServe(":"+*port, mux))
}
