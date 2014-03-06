package engine

import (
  "fmt"
  log "github.com/cihub/seelog"
  "net/http"
  "net/http/httptest"
  "reflect"
  "testing"
)

var verificationEvaluationTCs = []struct {
  evaluation string
  results    []verificationResult
}{
  {
    `
    [
      {
        "name" : "Basic evaluation",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "ec2",
          "region": "us-west-2"
        },
        "collection": {
          "name": "instance"
        },
        "filters" : {
          "instance": {
            "include" : [],
            "exclude" : [],
            "attributes" : {}
          }
        },
        "verifications": {
          "image_id" : {
            "value" : ["img-01"]
          }
        },
        "reports" : {}
      }
    ]
    `,
    []verificationResult{newVerificationResult("1", "", true), newVerificationResult("2", "", false)},
  },
  {
    `
    [
      {
        "name" : "Validate equal against multiple value attribute",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "ec2",
          "region": "us-west-2"
        },
        "collection": {
          "name": "instance"
        },
        "verifications": {
          "groups" : {
            "comparison" : ["equals"],
            "value" : ["web","admin"]
          }
        },
        "reports" : {}
      }
    ]
    `,
    []verificationResult{newVerificationResult("1", "", true), newVerificationResult("2", "", false)},
  },
  {
    `
    [
      {
        "name" : "Test include filter.",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "ec2",
          "region": "us-west-2"
        },
        "collection": {
          "name": "instance"
        },
        "filters" : {
          "instance" : {
            "include" : ["1"]
          }
        },
        "verifications": {
          "groups" : {
            "comparison" : ["include"],
            "value" : ["admin"]
          }
        }
      }
    ]
    `,
    []verificationResult{newVerificationResult("1", "", true)},
  },
  {
    `
    [
      {
        "name" : "Test exclude filters.",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "ec2",
          "region": "us-west-2"
        },
        "collection": {
          "name": "instance"
        },
        "filters" : {
          "instance" : {
            "exclude" : ["1"]
          }
        },
        "verifications": {
          "groups" : {
            "comparison" : ["include"],
            "value" : ["admin"]
          }
        }
      }
    ]
    `,
    []verificationResult{newVerificationResult("2", "", true)},
  },
  {
    `
    [
      {
        "name" : "Test filtering equals attributes.",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "ec2",
          "region": "us-west-2"
        },
        "collection": {
          "name": "instance"
        },
        "filters" : {
          "instance" : {
            "attributes" : {
              "groups" : {
                "equals" : ["db", "admin"]
              }
            }
          }
        },
        "verifications": {
          "image_id" : {
            "comparison" : ["equals"],
            "value" : ["img-01"]
          }
        }
      }
    ]
    `,
    []verificationResult{newVerificationResult("2", "", false)},
  },
  {
    `
    [
      {
        "name" : "Test filtering within attributes.",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "ec2",
          "region": "us-west-2"
        },
        "collection": {
          "name": "instance"
        },
        "filters" : {
          "instance" : {
            "attributes" : {
              "groups" : {
                "within" : ["bar", "db", "admin", "replication", "foo"]
              }
            }
          }
        },
        "verifications": {
          "image_id" : {
            "comparison" : ["equals"],
            "value" : ["img-01"]
          }
        }
      }
    ]
    `,
    []verificationResult{newVerificationResult("2", "", false)},
  },
}

var reportEvaluationTCs = []struct {
  evaluation string
  results    []reportResult
}{
  {
    `
    [
      {
        "name" : "Test equals report.",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "ec2",
          "region": "us-west-2"
        },
        "collection": {
          "name": "instance"
        },
        "filters" : {
          "instance" : {
            "exclude" : ["1"]
          }
        },
        "reports": {
          "equals" : [ "2" ]
        }
      }
    ]
    `,
    []reportResult{newReportResult(true)},
  },
  {
    `
    [
      {
        "name" : "Test greater_than report.",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "ec2",
          "region": "us-west-2"
        },
        "collection": {
          "name": "instance"
        },
        "filters" : {
          "instance" : {
            "exclude" : ["1"]
          }
        },
        "reports": {
          "greater_than" : [ "1" ]
        }
      }
    ]
    `,
    []reportResult{newReportResult(false)},
  },
  {
    `
    [
      {
        "name" : "Test less_than report.",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "ec2",
          "region": "us-west-2"
        },
        "collection": {
          "name": "instance"
        },
        "filters" : {
          "instance" : {
            "exclude" : ["1"]
          }
        },
        "reports": {
          "less_than" : [ "10" ]
        }
      }
    ]
    `,
    []reportResult{newReportResult(true)},
  },
}

var recursiveEvaluationTCs = []struct {
  evaluation string
  results    []reportResult
}{
  {
    `
    [
      {
        "name" : "Test equals for child resource.",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "s3",
          "region": "us-west-2"
        },
        "collection": {
          "name": "bucket-object"
        },
        "reports": {
          "equals" : [ "x", "y", "z" ]
        }
      }
    ]
    `,
    []reportResult{newReportResult(true)},
  },
  {
    `
    [
      {
        "name" : "Filter by parent include.",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "s3",
          "region": "us-west-2"
        },
        "filters" : {
          "bucket" : {
            "include" : ["a"]
          }
        },
        "collection": {
          "name": "bucket-object"
        },
        "reports": {
          "equals" : [ "x", "y" ]
        }
      }
    ]
    `,
    []reportResult{newReportResult(true)},
  },
  {
    `
    [
      {
        "name" : "Filter by parent exclude.",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "s3",
          "region": "us-west-2"
        },
        "filters" : {
          "bucket" : {
            "exclude" : ["a"]
          }
        },
        "collection": {
          "name": "bucket-object"
        },
        "reports": {
          "equals" : [ "z" ]
        }
      }
    ]
    `,
    []reportResult{newReportResult(true)},
  },
  {
    `
    [
      {
        "name" : "Filter by parent attributes.",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "s3",
          "region": "us-west-2"
        },
        "filters" : {
          "bucket" : {
            "attributes" : { 
              "versioning_enabled" : {
                "equals" : ["true"]
              }
            }
          }
        },
        "collection": {
          "name": "bucket-object"
        },
        "reports": {
          "equals" : [ "x", "y" ]
        }
      }
    ]
    `,
    []reportResult{newReportResult(true)},
  },
  {
    `
    [
      {
        "name" : "Filter by parent and collection attributes.",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "s3",
          "region": "us-west-2"
        },
        "filters" : {
          "bucket" : {
            "attributes" : { 
              "versioning_enabled" : {
                "equals" : ["true"]
              }
            }
          },
          "bucket-object" : {
            "attributes" : { 
              "public" : {
                "equals" : ["true"]
              }
            }
          }
        },
        "collection": {
          "name": "bucket-object"
        },
        "reports": {
          "equals" : [ "y" ]
        }
      }
    ]
    `,
    []reportResult{newReportResult(true)},
  },
  {
    `
    [
      {
        "name" : "Filter by parent and collection exclude.",
        "account": {
          "name" : "aws",
          "access_key": "XXX",
          "secret_key": "YYY"
        },
        "provider": {
          "name": "s3",
          "region": "us-west-2"
        },
        "filters" : {
          "bucket" : {
            "attributes" : {
              "versioning_enabled" : {
                "equals" : ["true"]
              }
            }
          },
          "bucket-object" : {
            "exclude" : ["y"]
          }
        },
        "collection": {
          "name": "bucket-object"
        },
        "reports": {
          "equals" : [ "x" ]
        }
      }
    ]
    `,
    []reportResult{newReportResult(true)},
  },
}

var argData string = `
[
    "region"
]
`

var credData string = `
[
    "access_key",
    "secret_key"
]
`

var instanceStateData string = `
[
    {
        "id": "1",
        "data": {
            "image_id": [
                "img-01"
            ],
            "groups": [
                "admin",
                "web"
            ]
        }
    },
    {
        "id": "2",
        "data": {
            "image_id": [
                "img-02"
            ],
            "groups": [
                "admin",
                "db"
            ]
        }
    }
]
`

var bucketStateData string = `
[
    {
        "id": "a",
        "data": {
            "versioning_enabled": [
                "true"
            ]
        }
    },
    {
        "id": "b",
        "data": {
            "versioning_enabled": [
                "false"
            ]
        }
    }
]
`

var bucketObjectStateData string = `
[
    {
        "id": "x",
        "parent_id": "a",
        "data": {
            "public": [
                "false"
            ]
        }
    },
    {
        "id": "y",
        "parent_id": "a",
        "data": {
            "public": [
                "true"
            ]
        }
    },
    {
        "id": "z",
        "parent_id": "b",
        "data": {
            "public": [
                "false"
            ]
        }
    }
]
`

var server *httptest.Server

func setup() {
  mux := http.NewServeMux()
  mux.HandleFunc("/accounts/aws/credentials", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, credData)
  })
  mux.HandleFunc("/accounts/aws/providers/ec2/arguments", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, argData)
  })
  mux.HandleFunc("/accounts/aws/providers/s3/arguments", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, argData)
  })
  mux.HandleFunc("/accounts/aws/providers/ec2/collections/instance/state", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, instanceStateData)
  })
  mux.HandleFunc("/accounts/aws/providers/s3/collections/bucket/state", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, bucketStateData)
  })
  mux.HandleFunc("/accounts/aws/providers/s3/collections/bucket-object/state", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, bucketObjectStateData)
  })
  server = httptest.NewServer(mux)
}

func teardown() {
  server.Close()
  log.Flush()
}

func TestVerificationExecute(t *testing.T) {
  setup()
  defer teardown()

  for count, tc := range verificationEvaluationTCs {
    t.Logf("Test case %d.", count)
    e := NewEngine()
    if err := e.LoadEvaluationsFromJSON([]byte(tc.evaluation)); err != nil {
      t.Fatalf(err.Error())
    }
    rc := NewRunConfig(120, 0)
    rc.SetOverrideEndpoint(server.URL)
    if err := e.Execute(rc); err != nil {
      t.Fatalf(err.Error())
    }
    evaluation := e.Evaluations()[0]
    t.Logf("Processing evaluation '%s'.", evaluation.Name)

    if !reflect.DeepEqual(evaluation.Verifications[0].Results, tc.results) {
      t.Fatalf("%+v \n does not equal \n %+v", evaluation.Verifications[0].Results, tc.results)
    }
  }
}

func TestReportExecute(t *testing.T) {
  setup()
  defer teardown()

  for _, tc := range reportEvaluationTCs {
    e := NewEngine()
    if err := e.LoadEvaluationsFromJSON([]byte(tc.evaluation)); err != nil {
      t.Fatalf(err.Error())
    }
    rc := NewRunConfig(120, 0)
    rc.SetOverrideEndpoint(server.URL)
    if err := e.Execute(rc); err != nil {
      t.Fatalf(err.Error())
    }
    evaluation := e.Evaluations()[0]
    t.Logf("Processing evaluation '%s'.", evaluation.Name)

    if !reflect.DeepEqual(evaluation.Reports[0].Results, tc.results) {
      t.Fatalf("%+v \n does not equal \n %+v", evaluation.Reports[0].Results, tc.results)
    }
  }
}

func TestRecursiveExecute(t *testing.T) {
  setup()
  defer teardown()

  for _, tc := range recursiveEvaluationTCs {
    e := NewEngine()
    if err := e.LoadEvaluationsFromJSON([]byte(tc.evaluation)); err != nil {
      t.Fatalf(err.Error())
    }
    rc := NewRunConfig(120, 0)
    rc.SetOverrideEndpoint(server.URL)
    if err := e.Execute(rc); err != nil {
      t.Fatalf(err.Error())
    }
    evaluation := e.Evaluations()[0]
    t.Logf("Processing evaluation '%s'.", evaluation.Name)

    if !reflect.DeepEqual(evaluation.Reports[0].Results, tc.results) {
      t.Fatalf("%+v \n does not equal \n %+v", evaluation.Reports[0].Results, tc.results)
    }
  }
}
