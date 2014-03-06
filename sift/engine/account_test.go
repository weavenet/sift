package engine

import (
  "testing"
)

type accountTest struct {
  name        string
  credentials map[string]string
  id          string
}

var creds1 = map[string]string{"test1": "cred1"}
var creds2 = map[string]string{"test1": "cred2"}
var creds3 = map[string]string{"test2": "cred2"}

var accountTestCases = []accountTest{
  {"acct1", creds1, "469f129f76332a4334718960c99f739b2362408b"},
  {"acct2", creds1, "469f129f76332a4334718960c99f739b2362408b"},
  {"acct1", creds2, "f6b18e0e420cbe39d2b1f1b9a7d51bb2cc29da28"},
  {"acct1", creds3, "e286d560c972c6d86f6c63ae4f92fb079eaee27a"},
}

func TestAccountId(t *testing.T) {
  for _, tc := range accountTestCases {
    a := newAccount(tc.name, tc.credentials)
    if a.id() != tc.id {
      t.Errorf("ID %s does not equal %s.", a.id(), tc.id)
    }
  }
}
