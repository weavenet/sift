package target

import (
  "fmt"
  "github.com/brettweavnet/sift/sift/state"
  log "github.com/cihub/seelog"
  "reflect"
  "sort"
  "strconv"
)

var reportFuncs = map[string]func([]state.State, []string) (map[string]bool, error){
  "equals":       reportEquals,
  "greater_than": reportGreaterThan,
  "less_than":    reportLessThan,
}

var filterFuncs = map[string]func(string, []string, []state.State) (map[string]bool, error){
  "equals":   valueEquals,
  "include":  valueIncludes,
  "includes": valueIncludes,
  "exclude":  valueExcludes,
  "excludes": valueExcludes,
  "within":   valueWithin,
  "set":      valueSet,
}

func runReport(name string, value []string, states []state.State) (map[string]bool, error) {
  if val, ok := reportFuncs[name]; ok {
    return val(states, value)
  }
  return map[string]bool{}, fmt.Errorf("Unknown report '%s'.", name)
}

func runComparison(name string, comparison string, value []string, states []state.State) (map[string]bool, error) {
  if val, ok := filterFuncs[comparison]; ok {
    return val(name, value, states)
  }
  return map[string]bool{}, fmt.Errorf("Unknown comparison '%s'.", comparison)
}

func valueEquals(name string, value []string, states []state.State) (map[string]bool, error) {
  r := map[string]bool{}
  sort.Strings(value)
  for _, s := range states {
    stateValue, err := s.Get(name)
    if err != nil {
      return r, err
    }
    sort.Strings(stateValue)
    if reflect.DeepEqual(stateValue, value) {
      r[s.Id] = true
    } else {
      r[s.Id] = false
    }
  }
  return r, nil
}

func valueIncludes(name string, value []string, states []state.State) (map[string]bool, error) {
  r := map[string]bool{}
  for _, s := range states {
    stateValue, err := s.Get(name)
    if err != nil {
      return r, err
    }
    r[s.Id] = true
    log.Tracef("Validating '%s' includes '%s'.", stateValue, value)
    v := false
    for _, sv := range stateValue {
      for _, cv := range value {
        if cv == sv {
          v = true
        }
      }
    }
    r[s.Id] = v
  }
  return r, nil
}

func valueExcludes(name string, value []string, states []state.State) (map[string]bool, error) {
  r := map[string]bool{}
  for _, s := range states {
    stateValue, err := s.Get(name)
    if err != nil {
      return r, err
    }
    r[s.Id] = true
    v := true
    for _, sv := range stateValue {
      for _, cv := range value {
        if cv == sv {
          v = false
        }
      }
    }
    log.Tracef("Validating '%s' does not include '%s' result '%t'.", stateValue, value, v)
    r[s.Id] = v
  }
  return r, nil
}

func valueSet(name string, value []string, states []state.State) (map[string]bool, error) {
  r := map[string]bool{}
  for _, s := range states {
    stateValue, err := s.Get(name)
    if err != nil {
      return r, err
    }
    if !reflect.DeepEqual(stateValue, []string{""}) {
      if value[0] == "true" {
        r[s.Id] = true
      } else {
        r[s.Id] = false
      }
    } else {
      if value[0] == "true" {
        r[s.Id] = false
      } else {
        r[s.Id] = true
      }
    }
  }
  return r, nil
}

func valueWithin(name string, value []string, states []state.State) (map[string]bool, error) {
  r := map[string]bool{}

  for _, s := range states {
    stateValue, err := s.Get(name)
    if err != nil {
      return r, err
    }
    r[s.Id] = false
    log.Tracef("Validating '%s' is represented in '%s'.", stateValue, value)
    matchCount := 0
    for _, sv := range stateValue {
      for _, cv := range value {
        if cv == sv {
          matchCount = matchCount + 1
        }
      }
    }
    valCount := len(stateValue)
    if valCount == matchCount {
      r[s.Id] = true
    } else {
      r[s.Id] = false
    }
  }

  return r, nil
}

func reportEquals(states []state.State, value []string) (map[string]bool, error) {
  ids := []string{}
  for _, s := range states {
    ids = append(ids, s.Id)
  }
  sort.Strings(ids)
  sort.Strings(value)
  eq := reflect.DeepEqual(ids, value)
  return map[string]bool{"equals": eq}, nil
}

func reportLessThan(states []state.State, value []string) (map[string]bool, error) {
  i, err := strconv.Atoi(value[0])
  if err != nil {
    return map[string]bool{}, fmt.Errorf("Error running less_than report '%s'.", err)
  }
  res := len(states) < i
  return map[string]bool{"less_than": res}, nil
}

func reportGreaterThan(states []state.State, value []string) (map[string]bool, error) {
  i, err := strconv.Atoi(value[0])
  if err != nil {
    return map[string]bool{}, fmt.Errorf("Error running greater_than report '%s'.", err)
  }
  res := len(states) > i
  return map[string]bool{"greater_than": res}, nil
}
