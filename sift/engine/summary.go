package engine

import (
  "fmt"
  "strings"

  log "github.com/cihub/seelog"
)

type summary struct {
  run *run
}

func NewSummary(r *run) *summary {
  return &summary{run: r}
}

func (s *summary) Output() {
  bd := s.breakdown()
  for name, metric := range bd {
    total := float64(metric["total"])
    if total > 0 {
      pass := float64(metric["pass"])
      percent := ((pass / total) * 100)
      log.Infof("%3.2f %% passed (%d / %d) %s passed.", percent, metric["pass"], metric["total"], name)
    } else {
      log.Infof("No %s evaluated.", name)
    }
  }
  for _, fr := range s.failedReports() {
    log.Infof(fr)
  }
}

func (s *summary) breakdown() map[string]map[string]int {
  bd := map[string]map[string]int{}
  bd["evaluations"] = s.evaluationsBreakdown()
  bd["verifications"] = s.verificationsBreakdown()
  bd["resources"] = s.resultsBreakdown()
  bd["reports"] = s.reportsBreakdown()
  return bd
}

func (s *summary) evaluationsBreakdown() map[string]int {
  var pass, fail, total int
  total = len(s.run.Evaluations)
  log.Tracef("Summarizing evaluation results.")
  for _, e := range s.run.Evaluations {
    if e.pass() {
      pass = pass + 1
    } else {
      fail = fail + 1
    }
  }
  return map[string]int{"total": total, "pass": pass, "fail": fail}
}

func (s *summary) verificationsBreakdown() map[string]int {
  var pass, fail, total int
  log.Tracef("Summarizing verification results.")
  for _, e := range s.run.Evaluations {
    for _, v := range e.Verifications {
      total = total + 1
      if v.pass() {
        pass = pass + 1
      } else {
        fail = fail + 1
      }
    }
  }
  return map[string]int{"total": total, "pass": pass, "fail": fail}
}

func (s *summary) resultsBreakdown() map[string]int {
  var pass, fail, total int
  log.Tracef("Summarizing individual results.")
  for _, e := range s.run.Evaluations {
    for _, v := range e.Verifications {
      for _, r := range v.Results {
        total = total + 1
        if r.Pass() {
          pass = pass + 1
        } else {
          fail = fail + 1
        }
      }
    }
  }
  return map[string]int{"total": total, "pass": pass, "fail": fail}
}

func (s *summary) reportsBreakdown() map[string]int {
  var pass, fail, total int
  log.Tracef("Summarizing report results.")
  for _, e := range s.run.Evaluations {
    for _, r := range e.Reports {
      total = total + 1
      if r.pass() {
        pass = pass + 1
      } else {
        fail = fail + 1
      }
    }
  }
  return map[string]int{"total": total, "pass": pass, "fail": fail}
}

func (s *summary) failedReports() (failed []string) {
  for _, e := range s.run.Evaluations {
    a := e.Account
    c := e.Collection
    p := e.Provider

    for _, r := range e.Reports {
      n := r.Name
      v := strings.Join(r.Value, ", ")
      if !r.pass() {
        failed = append(failed, fmt.Sprintf("Report '%s' value '%s' failed for '%s %s %s'.", n, v, a, p, c))
      }
    }
  }
  return failed
}

func (s *summary) failedVerifications() (failed []string) {
  for _, e := range s.run.Evaluations {
    a := e.Account
    c := e.Collection
    p := e.Provider

    for _, r := range e.Verifications {
      comp := r.Comparison
      n := r.Name
      v := strings.Join(r.Value, ", ")
      if !r.pass() {
        failed = append(failed, fmt.Sprintf("Verification '%s' is '%s' value '%s' failed for '%s %s %s'.", n, comp, v, a, p, c))
      }
    }
  }
  return failed
}
