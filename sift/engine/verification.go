package engine

import (
  "fmt"
  log "github.com/cihub/seelog"
  "github.com/brettweavnet/sift/sift/target"
)

type verification struct {
  Name       string               `json:"name"`
  Value      []string             `json:"value"`
  Comparison string               `json:"comparison"`
  Results    []verificationResult `json:"results"`
  evaluation *evaluation
}

func (v verification) String() string {
  return fmt.Sprintf("Resource verification %s", v.Name)
}

func newVerification(name string) verification {
  return verification{Name: name}
}

func (v *verification) setEvaluation(e *evaluation) {
  v.evaluation = e
}

func (v verification) fail() (t bool) {
  return !v.pass()
}

func (v verification) pass() (t bool) {
  for _, r := range v.Results {
    if r.Pass() == false {
      return false
    }
  }
  return true
}

func (v verification) resultsFail() (results []verificationResult) {
  for _, r := range v.Results {
    if r.Fail() {
      results = append(results, r)
    }
  }
  return results
}

func (v *verification) perform(cntxt *context, cash *cache, filters map[string]filter) error {
  t := *target.NewTarget(cntxt.Account.Name, cntxt.Provider.Name, cntxt.Collection.Name)

  filteredState, err := applyFiltersToContextLineage(cntxt, filters, cash)
  if err != nil {
    return err
  }

  log.Debugf("Performing verification '%s' with value '%s'.", v.Name, v.Value)
  results, err := t.Verify(v.Name, v.Comparison, v.Value, filteredState)
  if err != nil {
    return err
  }

  for id, pass := range results {
    rslt := newVerificationResult(id, cntxt.ParentId, pass)
    v.Results = append(v.Results, rslt)
  }
  if v.pass() {
    log.Debugf("Verification passed.")
  } else {
    log.Debugf("Verification failed.")
  }
  return nil
}
