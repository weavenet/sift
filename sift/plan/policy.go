package plan

import (
  "fmt"
  log "github.com/cihub/seelog"
  "strings"
)

type policy struct {
  Scope         string                            `json:"scope"`
  Arguments     string                            `json:"arguments"`
  Source        string                            `json:"source"`
  Filters       map[string]string                 `json:"filters"`
  Verifications map[string]map[string]interface{} `json:"verifications"`
  Reports       map[string]interface{}            `json:"reports"`
}

func (pol policy) toEvaluations(p Plan) ([]evaluation, error) {
  evaluations := []evaluation{}

  log.Debugf("Evaluation source '%s'.", pol.Source)
  accounts, err := pol.loadAccountsFromPlan(p)
  if err != nil {
    return evaluations, err
  }

  argumentsName := pol.Arguments
  if argumentsName == "" {
    log.Debugf("Arguments name to load not set, loading default.")
    argumentsName = "default"
  }

  log.Tracef("Loading arguments '%s' for source '%s'.", argumentsName, pol.Source)
  arguments := p.sourceArguments(pol.Source, argumentsName)

  log.Tracef("Arguments '%+v' loaded from plan.", arguments)

  return pol.buildEvaluations(accounts, arguments, p)
}

func (pol policy) providerName() string {
  return strings.Split(pol.Source, "_")[1]
}

func (pol policy) accountName() string {
  return strings.Split(pol.Source, "_")[0]
}

func (pol policy) collectionName() string {
  return strings.Split(pol.Source, "_")[len(strings.Split(pol.Source, "_"))-1]
}

func (pol policy) loadAccountsFromPlan(p Plan) ([]account, error) {
  scope := pol.Scope
  log.Tracef("Loading accounts from '%s' scoped to '%s'.", pol.accountName(), scope)

  evaluationAccounts := p.scopeToAccounts(pol.accountName(), scope)
  if len(evaluationAccounts) == 0 {
    return []account{}, fmt.Errorf("No accounts in scope '%s'.", scope)
  }
  return evaluationAccounts, nil
}

func (pol policy) buildEvaluations(accounts []account, arguments map[string][]string, p Plan) ([]evaluation, error) {
  evaluations := []evaluation{}

  log.Tracef("Loading filters '%s'.", pol.Filters)
  fltrs, err := p.convertPlanFiltersToEvaluationFilters(pol.Filters)
  if err != nil {
    return evaluations, err
  }

  for _, acct := range accounts {
    if len(arguments) == 0 {
      log.Tracef("No arguments found for '%s' '%s'.", pol.Source, arguments)
      e, err := pol.buildEvaluation(acct, fltrs, p)
      if err != nil {
        return evaluations, err
      }
      evaluations = append(evaluations, e)
    } else {
      argumentMatrix := buildArgumentMatrix(arguments)

      for _, matrixCase := range argumentMatrix {
        // Create a new evaluation witch each argument case in the matrix
        e, err := pol.buildEvaluation(acct, fltrs, p)
        if err != nil {
          return evaluations, err
        }
        e.Provider["name"] = pol.providerName()
        for argName, argValue := range matrixCase {
          e.Provider[argName] = argValue
        }
        evaluations = append(evaluations, e)
      }
    }
  }
  return evaluations, nil
}

func (pol policy) buildEvaluation(a account, filters map[string]filter, p Plan) (evaluation, error) {
  log.Tracef("Building evaluation for account '%s' with no arguments.", a)
  e := newEvaluation()
  for k, v := range a.Credentials {
    e.Account[k] = v
  }
  e.Account["name"] = pol.accountName()
  e.Provider["name"] = pol.providerName()
  e.Collection["name"] = pol.collectionName()
  e.Filters = filters

  // Create empty filter if one is not specified
  if _, ok := e.Filters[pol.collectionName()]; !ok {
    nf := newFilter()
    e.Filters[pol.collectionName()] = *nf
  }

  verifications, err := pol.parseVerifications(p)
  if err != nil {
    return *e, err
  }
  e.Verifications = verifications

  reports, err := p.setDynamicValues(pol.Reports)
  if err != nil {
    return *e, err
  }
  e.Reports = reports

  log.Tracef("Created evaluation '%s'", e)
  return *e, nil
}

func (pol policy) parseVerifications(p Plan) (r map[string]map[string][]string, err error) {
  r = map[string]map[string][]string{}
  for attribute, data := range pol.Verifications {
    r[attribute] = map[string][]string{}

    //for key, value := range data {
    updatedVerifications, err := p.setDynamicValues(data)
    if err != nil {
      return r, err
    }
    for k, v := range updatedVerifications {
      r[attribute][k] = v
    }
    //}
  }
  return r, err
}

func buildArgumentMatrix(args map[string][]string) []map[string]string {
  matrix := make([]map[string]string, 0)

  log.Debugf("Generating argument matrix from '%+v'.", args)
  for n, v := range args {
    matrix = expandMatrix(matrix, n, v)
  }
  log.Debugf("Built argument matrix '%+v'.", matrix)

  return matrix
}

func expandMatrix(matrix []map[string]string, entryName string, entryValues []string) []map[string]string {
  expandedMatrix := make([]map[string]string, 0)

  if len(matrix) == 0 {
    for _, v := range entryValues {
      row := make(map[string]string)
      row[entryName] = v
      expandedMatrix = append(expandedMatrix, row)
    }
  } else {
    for _, row := range matrix {
      for _, v := range entryValues {
        // Not sure why, but creating a new entry is necessary
        // as appending row will cause it to be changed as
        // entryValues is iterated over
        newEntry := make(map[string]string)
        for k, v := range row {
          newEntry[k] = v
        }
        newEntry[entryName] = v
        expandedMatrix = append(expandedMatrix, newEntry)
      }
    }
  }
  return expandedMatrix
}
