package target

import (
  log "github.com/cihub/seelog"
  "github.com/siftproject/sift/sift/state"
)

type Target struct {
  states []state.State
}

func NewTarget(accountName string, providerName string, collectionName string) *Target {
  t := Target{}
  return &t
}

func (t Target) FilterState(include []string, exclude []string, attrs map[string]map[string][]string, states []state.State) ([]state.State, error) {
  fs := filterIncludeExclude(include, exclude, states)

  filteredState, err := filterAttributes(attrs, fs)
  if err != nil {
    return []state.State{}, err
  }
  return filteredState, nil
}

func (t Target) Verify(name string, comparison string, value []string, states []state.State) (results map[string]bool, err error) {
  if comparison == "report" {
    log.Debugf("Running report '%s'.", name)
    return runReport(name, value, states)
  } else {
    log.Debugf("Running attribute verification '%s' against '%s'.", comparison, name)
    return runComparison(name, comparison, value, states)
  }
}
