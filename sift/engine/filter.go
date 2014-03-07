package engine

import (
  "fmt"
  log "github.com/cihub/seelog"
  "github.com/brettweavnet/sift/sift/state"
  "github.com/brettweavnet/sift/sift/target"
)

type filter struct {
  Include    []string      `json:"include"`
  Exclude    []string      `json:"exclude"`
  Attributes attributeList `json:"attributes"`
}

type attributeList map[string]map[string][]string

func newFilter() filter {
  return filter{
    Include:    []string{},
    Exclude:    []string{},
    Attributes: attributeList{},
  }
}

func (f filter) apply(t target.Target, states []state.State) ([]state.State, error) {
  for _, v := range f.Include {
    log.Debugf("Including resource '%s'", v)
  }
  for _, v := range f.Exclude {
    log.Debugf("Excluding resource '%s'", v)
  }
  for name, value := range f.Attributes {
    log.Debugf("Applying attribute filter '%s' = '%s'", name, value)
  }
  log.Debugf("Pre filtered list includes %d states.", len(states))
  log.Debugf("Filtering state with include: '%s' exclude: '%s' attributes: '%'.", f.Include, f.Exclude, f.Attributes)
  fs, err := t.FilterState(f.Include, f.Exclude, f.Attributes, states)
  if err != nil {
    return []state.State{}, err
  }
  log.Debugf("Filtered list includes %d states.", len(fs))
  return fs, nil
}

func applyFiltersToContextLineage(cntxt *context, filters map[string]filter, cash *cache) ([]state.State, error) {
  t := *target.NewTarget(cntxt.Account.Name, cntxt.Provider.Name, cntxt.Collection.Name)

  var f filter
  var filteredState []state.State
  var err error

  log.Debugf("Applying filters to context '%s'.", cntxt.id)

  filteredState, err = recursiveFilterCollectionByParent(cntxt, filters, cash)
  if err != nil {
    return filteredState, err
  }

  if val, ok := filters[cntxt.Collection.Name]; ok {
    log.Tracef("Found filter for collection '%s'.", cntxt.Collection.Name)
    f = val
  } else {
    log.Tracef("No filter set for collection '%s'.", cntxt.Collection.Name)
    f = newFilter()
  }

  return f.apply(t, filteredState)
}

func recursiveFilterCollectionByParent(cntxt *context, filters map[string]filter, cash *cache) ([]state.State, error) {
  col := cntxt.Collection
  acct := cntxt.Account
  prov := cntxt.Provider
  tgt := *target.NewTarget(acct.Name, prov.Name, col.Name)

  var filteredState []state.State
  var err error
  var f filter

  if col.hasParent() {
    log.Debugf("Context '%s' has parent, applying parent filters.", cntxt.id)
    filteredState, err = filterCollectionByParent(cntxt, filters, cash)
    if err != nil {
      return filteredState, err
    }
  } else {
    log.Debugf("Context '%s' has no parent, continuing.", cntxt.id)
    filteredState = cntxt.State
  }
  log.Debugf("Post parent filtered state for context '%s' has %d states.", cntxt.id, len(filteredState))

  if val, ok := filters[col.Name]; ok {
    log.Tracef("Found filter for collection '%s'.", col.Name)
    f = val
  } else {
    log.Tracef("No filter found for collection '%s'.", col.Name)
    f = newFilter()
  }
  return f.apply(tgt, filteredState)
}

func filterCollectionByParent(cntxt *context, filters map[string]filter, cash *cache) ([]state.State, error) {
  col := cntxt.Collection

  var filteredState []state.State
  var filteredParentState []state.State
  var err error

  log.Debugf("Apply parent filter to collection '%s' parent '%s'.", col.Name, col.parent().Name)
  parentContext, ok := cash.GetContext(cntxt.ParentId)
  if !ok {
    return filteredState, fmt.Errorf("Error getting parent id '%s' from cache.", cntxt.ParentId)
  }

  filteredParentState, err = recursiveFilterCollectionByParent(parentContext, filters, cash)
  if err != nil {
    return []state.State{}, err
  }
  log.Debugf("Parent '%s' filtered to include %d states.", col.parent().Name, len(filteredParentState))

  for _, s := range cntxt.State {
    includeState := false
    for _, p := range filteredParentState {
      if p.Id == s.ParentId {
        includeState = true
      }
    }
    if includeState {
      filteredState = append(filteredState, s)
    }
  }
  return filteredState, nil
}
