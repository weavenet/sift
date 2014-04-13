package main

import (
  "flag"
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
  flag.Parse()
  serveEc2(*port)
}
