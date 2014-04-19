package plan

import (
  "encoding/json"
  "fmt"
  log "github.com/cihub/seelog"
)

type Plan struct {
  Accounts map[string]map[string]account `json:"accounts"`
  Policies []policy                      `json:"policies"`
  Filters  map[string]planFilter         `json:"filters"`
  Sources  map[string]source             `json:"sources"`
  Lists    map[string]list               `json:"lists"`
}

func NewPlan() *Plan {
  return &Plan{
    Lists:    map[string]list{},
    Accounts: map[string]map[string]account{},
    Filters:  map[string]planFilter{},
    Policies: []policy{},
    Sources:  map[string]source{},
  }
}

func (p Plan) Evaluations() ([]evaluation, error) {
  eos, err := p.createEvaluations()
  if err != nil {
    return []evaluation{}, err
  }
  data, err := json.MarshalIndent(eos, "", "  ")
  if err != nil {
    return []evaluation{}, err
  }
  log.Tracef("Created evaluation JSON: \n%s\n", data)
  return eos, nil
}

func (p *Plan) LoadJSON(data []byte) error {
  log.Tracef("Loading plan from JSON data '%s'.", string(data))
  if err := json.Unmarshal(data, p); err != nil {
    return err
  }
  log.Debugf("Successfully loaded plan.")
  return nil
}

func (p *Plan) LoadRepo(path string) error {
  log.Debugf("Loading repo from '%s'.", path)
  if !dirExists(path) {
    return fmt.Errorf("Repo directory does not exist.")
  }
  r := newRepo(path)
  if err := r.loadAccounts(p); err != nil {
    return err
  }
  if err := r.loadLists(p); err != nil {
    return err
  }
  if err := r.loadSources(p); err != nil {
    return err
  }
  if err := r.loadFilters(p); err != nil {
    return err
  }
  if err := r.loadPolicies(p); err != nil {
    return err
  }
  data, err := json.MarshalIndent(p, "", "  ")
  if err != nil {
    return err
  }
  log.Debugf("Loaded repo with plan '%s'.", data)
  return nil
}

func (p Plan) createEvaluations() ([]evaluation, error) {
  evaluations := []evaluation{}

  log.Debugf("Plan has %d policies to convert to evaluations.", len(p.Policies))
  for _, c := range p.Policies {
    es, err := c.toEvaluations(p)
    if err != nil {
      return []evaluation{}, err
    }
    log.Tracef("Creating %d evaluations from policiy.", len(evaluations))
    evaluations = append(evaluations, es...)
  }
  log.Infof("Plan includes %d evaluations.", len(evaluations))
  log.Infof("Evaluations created succesfully.")
  return evaluations, nil
}

func (p Plan) scopeToAccounts(name string, scope string) []account {
  accounts := []account{}
  for _, a := range p.Accounts[name] {
    for _, s := range a.Scope {
      switch scope {
      case s:
        {
          log.Debugf("Adding account '%s'. In scope '%s'.", a, scope)
          accounts = append(accounts, a)
        }
      case "":
        {
          log.Debugf("Adding account '%s' (no scope set).", a)
          accounts = append(accounts, a)
        }
      default:
        {
          log.Debugf("Account '%s' not in scope.", a)
        }
      }
    }
  }
  return accounts
}

func (p Plan) sourceArguments(source string, name string) (args map[string][]string) {
  for k, v := range p.Sources[source].Arguments {
    if k == name {
      log.Debugf("Arguments '%s' found for '%s'.", name, source)
      return v
    }
  }
  return args
}

func (p Plan) convertPlanFiltersToEvaluationFilters(filters map[string]string) (map[string]filter, error) {
  evaluationFilters := map[string]filter{}

  for key, value := range filters {
    if planFilter, ok := p.Filters[value]; ok {
      evaluationFilter := newFilter()
      log.Tracef("Building filter '%s' for '%s'.", value, key)
      inc, err := p.setDynamicValue(planFilter.Include)
      if err != nil {
        return evaluationFilters, err
      }
      if inc == nil {
        inc = []string{}
      }
      log.Tracef("Setting include to '%s'.", inc)
      evaluationFilter.Include = inc

      exc, err := p.setDynamicValue(planFilter.Exclude)
      if err != nil {
        return evaluationFilters, err
      }
      if exc == nil {
        exc = []string{}
      }
      log.Tracef("Setting exclude to '%s'.", exc)
      evaluationFilter.Exclude = exc

      attrs, err := p.setDynamicValues(planFilter.Attributes)
      if err != nil {
        return evaluationFilters, err
      }
      log.Tracef("Setting attributes to '%s'.", attrs)
      evaluationFilter.Attributes = attrs
      evaluationFilters[key] = *evaluationFilter
    }
  }

  log.Tracef("Created %d filters.", len(evaluationFilters))
  return evaluationFilters, nil
}

func (p Plan) setDynamicValues(data map[string]interface{}) (map[string][]string, error) {
  r := map[string][]string{}
  for name, value := range data {
    log.Tracef("Setting dynamic item '%s'.", name)
    values, err := p.setDynamicValue(value)
    if err != nil {
      return r, err
    }
    r[name] = values
  }
  return r, nil
}

func (p Plan) setDynamicValue(value interface{}) (values []string, err error) {
  if value == nil {
    return values, nil
  }
  if w, ok := value.([]string); ok {
    log.Tracef("Value '%v' is an array.", value)
    for _, v := range w {
      values = append(values, v)
    }
    return values, nil
  }

  if w, ok := value.([]interface{}); ok {
    log.Tracef("Value '%v' is an array.", value)
    for _, v := range w {
      values = append(values, v.(string))
    }
    return values, nil
  }

  if w, ok := value.(map[string]interface{}); ok {
    log.Tracef("Value '%v' is a map.", value)
    log.Tracef("Executing function '%s' with value '%v'", w, p.Lists)
    values, err := executeFn(value.(map[string]interface{}), p)
    if err != nil {
      return values, err
    }
    return values, nil
  }

  log.Tracef("Value does not appear to be a fn or slice. Converting from string to single value slice of strings.")
  values = append(values, value.(string))
  return values, nil
}
