package engine

import (
  "fmt"
  log "github.com/cihub/seelog"
)

func createContextsFromJSON(data []interface{}, cache *cache) error {
  for _, d := range data {
    p, err := newParser(d)
    if err != nil {
      return err
    }
    if err := createContextLineage(p, cache); err != nil {
      return err
    }
  }
  return nil
}

func createContextLineage(p parser, cache *cache) error {
  cntxt := newContext(p)

  log.Debugf("Context has ID '%s'.", cntxt.Id())
  if cntxt.hasParent() {
    parentParser := p
    parentCollection := cntxt.parent()

    parentParser.Collection = parentCollection
    log.Debugf("Context '%s' has parent collection '%s'", cntxt.Id(), parentCollection.Name)

    parentParser.setId()
    log.Debugf("Setting context '%s' parent to '%s'", cntxt.Id(), parentParser.Id)
    cntxt.ParentId = parentParser.Id

    createContextLineage(parentParser, cache)
  } else {
    log.Debugf("Context '%s' does not have a parent.", cntxt.Id())
  }
  log.Debugf("Adding context '%s' to cache.", cntxt.Id())
  cache.addContextToCache(cntxt)
  return nil
}

func createEvaluationsFromJSON(data []interface{}, cache *cache) (evaluations []evaluation, err error) {
  log.Debugf("Validating evaluations JSON.")
  if err := validateJSON(data); err != nil {
    return evaluations, err
  }
  log.Debugf("Validation completed succesfully.")
  for _, d := range data {
    parser, err := newParser(d)
    if err != nil {
      return evaluations, err
    }
    evaluation := newEvaluation(cache)
    if err = evaluation.loadFromParser(parser); err != nil {
      return evaluations, err
    }
    evaluations = append(evaluations, evaluation)
  }
  return
}

func validateJSON(data []interface{}) error {
  names := make(map[string]bool)

  for _, d := range data {
    parser, err := newParser(d)
    if err != nil {
      return err
    }
    if _, exists := names[parser.Name]; exists {
      return fmt.Errorf("Evaluations '%s' already declared.", parser.Name)
    }
    if parser.Name != "" {
      names[parser.Name] = true
    }
  }
  return nil
}
