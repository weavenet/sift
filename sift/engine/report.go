package engine

import (
  "github.com/brettweavnet/sift/sift/target"
  log "github.com/cihub/seelog"
)

type report struct {
  Name    string         `json:"name"`
  Value   []string       `json:"value"`
  Results []reportResult `json:"results"`
}

func newReport() report {
  return report{}
}

func (r report) pass() bool {
  // reports only have a single result
  return r.Results[0].Pass()
}

func (r report) fail() bool {
  return !r.pass()
}

func (r *report) perform(cntxt *context, cash *cache, filters map[string]filter) error {
  t := *target.NewTarget(cntxt.Account.Name, cntxt.Provider.Name, cntxt.Collection.Name)

  filteredState, err := applyFiltersToContextLineage(cntxt, filters, cash)
  if err != nil {
    return err
  }

  log.Debugf("Performing verification '%s' with value '%s'.", r.Name, r.Value)
  results, err := t.Verify(r.Name, "report", r.Value, filteredState)
  if err != nil {
    return err
  }

  for _, pass := range results {
    rslt := newReportResult(pass)
    r.Results = append(r.Results, rslt)
  }

  if r.pass() {
    log.Debugf("Report passed.")
  } else {
    log.Debugf("Report failed.")
  }

  return nil
}
