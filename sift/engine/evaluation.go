package engine

import (
  "fmt"
  "github.com/brettweavnet/sift/sift/util"
  log "github.com/cihub/seelog"
)

type evaluation struct {
  cache         *cache
  ContextId     string            `json:"context_id"`
  Name          string            `json:"name"`
  Filters       map[string]filter `json:"filters"`
  Verifications []*verification   `json:"verifications"`
  Reports       []*report         `json:"reports"`
}

func (e evaluation) String() (output string) {
  if e.Name != "" {
    output = fmt.Sprintf("%s", e.Name)
  } else {
    output = "unspecified"
  }
  return
}

func newEvaluation(cache *cache) evaluation {
  return evaluation{
    cache:         cache,
    Filters:       map[string]filter{},
    Verifications: []*verification{},
    Reports:       []*report{},
  }
}

func (e evaluation) fail() bool {
  return !e.pass()
}

func (e evaluation) pass() bool {
  for _, v := range e.Verifications {
    if v.fail() {
      return false
    }
  }
  for _, r := range e.Reports {
    if r.fail() {
      return false
    }
  }
  return true
}

func (e *evaluation) loadFromParser(p parser) (err error) {
  if err := e.setName(p); err != nil {
    return err
  }
  if len(p.Verifications) == 0 && len(p.Reports) == 0 {
    return fmt.Errorf("No reports of verifications specified for evaluation '%s'.", e.Name)
  }
  log.Tracef("Adding %d verifications to evaluation '%s'.", len(p.Verifications), e.Name)
  for _, v := range p.Verifications {
    log.Tracef("Adding verification '%s' to evaluation '%s'.", v.Name, e.Name)
    v.setEvaluation(e)
    e.Verifications = append(e.Verifications, v)
  }
  log.Tracef("Adding %d reports to evaluation '%s'.", len(p.Reports), e.Name)
  for _, r := range p.Reports {
    log.Tracef("Adding report '%s' to evaluation '%s'.", r.Name, e.Name)
    e.Reports = append(e.Reports, r)
  }
  e.setContextId(p.Id)

  e.setFilters(p.Filters)
  return nil
}

func (e *evaluation) performVerifications() error {
  log.Debugf("Performing verifications for evaluation '%s'", e.Name)
  for _, v := range e.Verifications {
    if err := v.perform(e.context(), e.cache, e.Filters); err != nil {
      return err
    }
  }
  return nil
}

func (e *evaluation) performReports() error {
  log.Debugf("Running reports for evaluation '%s'", e.Name)
  for _, r := range e.Reports {
    if err := r.perform(e.context(), e.cache, e.Filters); err != nil {
      return err
    }
  }
  return nil
}

func (e evaluation) context() *context {
  context, _ := e.cache.GetContext(e.ContextId)
  return context
}

func (e *evaluation) setContextId(id string) {
  log.Tracef("Loading context '%s' from cache for evaluation '%s'.", id, e.Name)
  e.ContextId = id
}

func (e *evaluation) setFilters(filters map[string]filter) {
  e.Filters = filters
}

func (e *evaluation) setName(p parser) error {
  if p.Name != "" {
    e.Name = p.Name
  } else {
    log.Debugf("Evaluation name not specified.")
    id, err := util.CreateUUID()
    if err != nil {
      return err
    }
    e.Name = "unspecified-" + id
  }
  log.Debugf("Setting evaluation '%s' from parser.", e.Name)
  return nil
}
