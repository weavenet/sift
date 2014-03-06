package engine

import (
  log "github.com/cihub/seelog"
)

type run struct {
  Cache       *cache       `json:"cache"`
  Evaluations []evaluation `json:"evaluations"`
  runConfig   runConfig
}

func NewRun(c *cache, e []evaluation, rc runConfig) *run {
  return &run{Cache: c,
    Evaluations: e,
    runConfig:   rc}
}

func (r *run) Execute() error {
  log.Debugf("Starting source endpoints.")
  sum := newEndpointManager(r.Cache, r.runConfig)
  defer sum.stopEndpointProcesses()
  if err := sum.setContextsEndpoints(); err != nil {
    return err
  }
  log.Infof("Validating evaluations.")
  if err := r.validateContexts(); err != nil {
    return err
  }
  log.Infof("Loading current state.")
  csl := newCacheStateLoader(r.Cache, r.runConfig)
  if err := csl.loadStates(); err != nil {
    return err
  }
  log.Infof("Processing evaluations.")
  if err := r.performVerifications(); err != nil {
    return err
  }
  if err := r.performReports(); err != nil {
    return err
  }
  log.Infof("Processing evaluations completed.")
  return nil
}

func (r *run) validateContexts() (err error) {
  for _, context := range r.Cache.Contexts() {
    if err := context.validate(); err != nil {
      return err
    }
  }
  return nil
}

func (r *run) performVerifications() error {
  for _, e := range r.Evaluations {
    if err := e.performVerifications(); err != nil {
      return err
    }
  }
  return nil
}

func (r *run) performReports() error {
  for _, e := range r.Evaluations {
    if err := e.performReports(); err != nil {
      return err
    }
  }
  return nil
}
