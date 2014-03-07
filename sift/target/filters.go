package target

import (
  "github.com/brettweavnet/sift/sift/state"
)

type attributeList map[string]map[string][]string

func newAttributeList(name string, comparison string, values []string) attributeList {
  return attributeList{name: map[string][]string{comparison: values}}
}

func filterIncludeExclude(i []string, e []string, states []state.State) (rs []state.State) {
  rs = states
  if len(i) > 0 {
    rs = include(i, rs)
  }
  if len(e) > 0 {
    rs = exclude(e, rs)
  }
  return rs
}

func filterAttributes(attrs attributeList, states []state.State) (rs []state.State, err error) {
  if len(attrs) > 0 {
    filtered, err := attributes(attrs, states)
    if err != nil {
      return []state.State{}, err
    }
    return filtered, nil
  }
  return states, nil
}

func exclude(ids []string, states []state.State) []state.State {
  fs := []state.State{}
  for _, s := range states {
    addState := true
    for _, id := range ids {
      if s.Id == id {
        addState = false
      }
    }
    if addState {
      fs = append(fs, s)
    }
  }
  return fs
}

func include(ids []string, states []state.State) []state.State {
  fs := []state.State{}
  for _, id := range ids {
    for _, s := range states {
      if s.Id == id {
        fs = append(fs, s)
      }
    }
  }
  return fs
}

func attributes(attrs attributeList, states []state.State) ([]state.State, error) {
  fs := make([]state.State, 0)
  for name, data := range attrs {
    for comparison, value := range data {
      r, err := runComparison(name, comparison, value, states)
      if err != nil {
        return fs, err
      }
      for _, s := range states {
        if r[s.Id] == true {
          fs = append(fs, s)
        }
      }
    }
  }

  return fs, nil
}
