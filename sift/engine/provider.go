package engine

import (
  "crypto/sha1"
  "fmt"
  "sort"
  "strings"
)

type provider struct {
  Name      string            `json:"name"`
  Arguments map[string]string `json:"arguments"`
}

func (p provider) String() (output string) {
  if p.Name == "" {
    output = "unspecified"
  } else {
    output = p.Name
  }
  return
}

func newProvider(name string, arguments map[string]string) provider {
  return provider{Name: name, Arguments: arguments}
}

func (p provider) id() (id string) {
  ordered := []string{}

  for key, value := range p.Arguments {
    ordered = append(ordered, fmt.Sprintf("%s=%s", key, value))
  }

  sort.Strings(ordered)
  c := strings.Join(ordered, "-")
  hasher := sha1.New()
  hasher.Write([]byte(c))
  id = fmt.Sprintf("%x", hasher.Sum(nil))
  return
}
