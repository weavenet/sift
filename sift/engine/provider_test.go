package engine

import (
  "testing"
)

type providerTest struct {
  name      string
  arguments map[string]string
  id        string
}

var args1 = map[string]string{"test1": "cred1"}
var args2 = map[string]string{"test1": "cred2"}
var args3 = map[string]string{"test2": "cred2"}

var providerTestCases = []providerTest{
  {"prov1", args1, "469f129f76332a4334718960c99f739b2362408b"},
  {"prov2", args1, "469f129f76332a4334718960c99f739b2362408b"},
  {"prov1", args2, "f6b18e0e420cbe39d2b1f1b9a7d51bb2cc29da28"},
  {"prov1", args3, "e286d560c972c6d86f6c63ae4f92fb079eaee27a"},
}

func TestProviderId(t *testing.T) {
  for _, tc := range providerTestCases {
    a := newProvider(tc.name, tc.arguments)
    if a.id() != tc.id {
      t.Errorf("ID %s does not equal %s.", a.id(), tc.id)
    }
  }
}

func TestProviderString(t *testing.T) {
  if newProvider("test123", args1).String() != "test123" {
    t.Errorf("Error setting name.")
  }
  if newProvider("", args1).String() != "unspecified" {
    t.Errorf("Error setting name.")
  }
}
