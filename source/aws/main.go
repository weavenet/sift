package main

import (
  "encoding/json"
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

type errorResponse struct {
  Message string `json:"message"`
}

func newStateRequest() stateRequest {
  return stateRequest{}
}

func newErrorResponse(e error) errorResponse {
  return errorResponse{Message: e.Error()}
}

func (e errorResponse) String() string {
  data, _ := json.Marshal(e)
  return string(data)
}

func main() {
  port := flag.String("port", "", "port to listen on")
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
