package engine

import (
  "encoding/json"
  log "github.com/cihub/seelog"
  "io/ioutil"
)

type engine struct {
  cache       *cache
  run         *run
  evaluations []evaluation
}

func NewEngine() *engine {
  c := newCache()
  return &engine{cache: c}
}

func (e engine) Evaluations() []evaluation {
  return e.evaluations
}

func (e *engine) LoadEvaluationsFromFile(path string) error {
  log.Infof("Loading evaluations from JSON file '%s'", path)
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return err
  }
  return e.LoadEvaluationsFromJSON(data)
}

func (e *engine) LoadEvaluationsFromJSON(rawEvaluationData []byte) error {
  var data []interface{}
  if err := json.Unmarshal(rawEvaluationData, &data); err != nil {
    return err
  }
  if err := createContextsFromJSON(data, e.cache); err != nil {
    return err
  }
  evaluations, err := createEvaluationsFromJSON(data, e.cache)
  if err != nil {
    return err
  }
  e.evaluations = evaluations
  return nil
}

func (e engine) Run() *run {
  return e.run
}

func (e *engine) Execute(rc runConfig) error {
  r := NewRun(e.cache, e.evaluations, rc)
  log.Debugf("Starting run.")
  if err := r.Execute(); err != nil {
    return err
  }
  e.run = r
  return nil
}
