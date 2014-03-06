package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
)

var server *http.Server

var credData string = `["access_key_id", "secret_access_key"]`

var iamArgData string = `[]`
var iamUserStateData string = `
[
  {
    "id" : "user1",
    "data" : { }
  },
  {
    "id" : "user2",
    "data" : { }
  }
]
`

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

var s3ArgData string = `["region"]`
var s3BucketStateData string = `
[
  {
    "id" : "bucket1",
    "data" : {}
  },
  {
    "id" : "bucket2",
    "data" : {}
  }
]`
var s3BucketObjectStateData string = `
[
  {
    "id" : "object1",
    "parent_id" : "bucket1",
    "data" : { "public" : ["false"] }
  },
  {
    "id" : "object2",
    "parent_id" : "bucket1",
    "data" : { "public" : ["true"] }
  },
  {
    "id" : "object3",
    "parent_id" : "bucket1",
    "data" : { "public" : ["false"] }
  },
  {
    "id" : "object4",
    "parent_id" : "bucket2",
    "data" : { "public" : ["true"] }
  }
]`

func setup(port string) {
  mux := http.NewServeMux()
  mux.HandleFunc("/accounts/test-aws/credentials", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, credData)
  })
  mux.HandleFunc("/accounts/test-aws/providers/ec2/arguments", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, ec2ArgData)
  })
  mux.HandleFunc("/accounts/test-aws/providers/iam/arguments", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, iamArgData)
  })
  mux.HandleFunc("/accounts/test-aws/providers/s3/arguments", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, s3ArgData)
  })
  mux.HandleFunc("/accounts/test-aws/providers/ec2/collections/instance/state", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, ec2InstanceStateData)
  })
  mux.HandleFunc("/accounts/test-aws/providers/iam/collections/user/state", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, iamUserStateData)
  })
  mux.HandleFunc("/accounts/test-aws/providers/s3/collections/bucket/state", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, s3BucketStateData)
  })
  mux.HandleFunc("/accounts/test-aws/providers/s3/collections/bucket-object/state", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, s3BucketObjectStateData)
  })
  log.Fatal(http.ListenAndServe(":"+port, mux))
}

func main() {
  port := flag.String("port", "32786", "port to listen on")
  flag.Parse()
  setup(*port)
}
