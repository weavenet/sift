package engine

import (
  "crypto/sha1"
  "fmt"
  "sort"
  "strings"
)

type account struct {
  Name        string
  credentials map[string]string
}

func newAccount(name string, credentials map[string]string) account {
  return account{Name: name, credentials: credentials}
}

func (a account) id() (id string) {
  ordered := []string{}

  for key, value := range a.credentials {
    ordered = append(ordered, fmt.Sprintf("%s=%s", key, value))
  }

  sort.Strings(ordered)
  c := strings.Join(ordered, "-")
  hasher := sha1.New()
  hasher.Write([]byte(c))
  id = fmt.Sprintf("%x", hasher.Sum(nil))
  return
}
