package engine

import (
  "strings"
)

type collection struct {
  Name string `json:"name"`
}

func newCollection(name string) collection {
  return collection{Name: name}
}

func (c collection) hasParent() bool {
  if len(c.lineage()) > 1 {
    return true
  }
  return false
}

func (c collection) parent() collection {
  return newCollection(c.lineage()[len(c.lineage())-2])
}

func (c collection) lineage() []string {
  l := make([]string, 0)
  for _, i := range strings.Split(c.Name, "-") {
    if len(l) > 0 {
      l = append(l, (l[len(l)-1] + "-" + i))
    } else {
      l = append(l, i)
    }
  }
  return l
}
