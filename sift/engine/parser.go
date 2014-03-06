package engine

import (
  "fmt"
  log "github.com/cihub/seelog"
  "github.com/mitchellh/mapstructure"
)

type parser struct {
  data          map[string]interface{}
  Name          string
  Account       account
  Provider      provider
  Collection    collection
  Filters       map[string]filter
  Id            string
  Verifications []*verification
  Reports       []*report
}

func newParser(data interface{}) (parser, error) {
  dm := data.(map[string]interface{})
  p := parser{data: dm}
  if err := p.load(); err != nil {
    return p, err
  }
  return p, nil
}

func (p *parser) load() error {
  p.loadName()
  if err := p.loadAccount(p.data["account"].(map[string]interface{})); err != nil {
    return err
  }
  if err := p.loadProvider(p.data["provider"].(map[string]interface{})); err != nil {
    return err
  }
  if err := p.loadCollection(p.data["collection"].(map[string]interface{})); err != nil {
    return err
  }
  if err := p.loadFilter(p.data["filters"]); err != nil {
    return err
  }
  if err := p.loadVerifications(p.data["verifications"]); err != nil {
    return err
  }
  if err := p.loadReports(p.data["reports"]); err != nil {
    return err
  }
  p.setId()
  return nil
}

func (p *parser) setId() {
  hash := fmt.Sprintf("%s-%s", p.Account.id(), p.Provider.id())
  p.Id = fmt.Sprintf("%s-%s", p.source(), hash)
}

func (p *parser) loadName() {
  if _, ok := p.data["name"]; ok {
    p.Name = p.data["name"].(string)
  } else {
    p.Name = ""
  }
}

func (p *parser) source() string {
  return fmt.Sprintf("%s-%s-%s", p.Account.Name, p.Provider.Name, p.Collection.Name)
}

func (p *parser) loadCollection(collection map[string]interface{}) (err error) {
  var name string

  if val, ok := collection["name"]; ok {
    name = val.(string)
  } else {
    return fmt.Errorf("Collection name not set.")
  }

  p.Collection = newCollection(name)
  return nil
}

func (p *parser) loadAccount(account map[string]interface{}) (err error) {
  credentials := make(map[string]string)
  var name string

  for k, v := range account {
    if k == "name" {
      name = v.(string)
    } else {
      credentials[k] = v.(string)
    }
  }
  if name == "" {
    return fmt.Errorf("Account name not set.")
  }
  p.Account = newAccount(name, credentials)
  return nil
}

func (p *parser) loadProvider(provider map[string]interface{}) (err error) {
  config := make(map[string]string)
  var name string

  for k, v := range provider {
    if k == "name" {
      name = v.(string)
    } else {
      config[k] = v.(string)
    }
  }
  if name == "" {
    return fmt.Errorf("Provider name not set.")
  }
  p.Provider = newProvider(name, config)
  return nil
}

func (p *parser) loadFilter(f interface{}) (err error) {
  var filters map[string]filter
  mapstructure.Decode(f, &filters)
  log.Tracef("Decoded raw filter data '%s' to '%+v'.", f, filters)
  for _, f := range filters {
    if (len(f.Include) > 0) && (len(f.Exclude) > 0) {
      return fmt.Errorf("Both include and exclude filters cannot be set in the same evaluation.")
    }
  }
  p.Filters = filters
  return nil
}

func (p *parser) loadVerifications(rawVerificationData interface{}) (err error) {
  var verificationData map[string]map[string][]string
  mapstructure.Decode(rawVerificationData, &verificationData)
  if len(verificationData) == 0 {
    log.Tracef("No verifications found, continuing.")
    return nil
  }
  var verifications []*verification
  for verificationName, data := range verificationData {
    verification := newVerification(verificationName)
    verification.Comparison = "equals"
    log.Tracef("Parsing verification '%s' with raw data '%s'.", verificationName, data)
    if err := loadVerificationData(&verification, data); err != nil {
      return err
    }
    log.Tracef("Verification value set to '%s'.", verification.Value)
    log.Tracef("Verification comparison set to '%s'.", verification.Comparison)
    verifications = append(verifications, &verification)
  }
  p.Verifications = verifications
  return nil
}

func loadVerificationData(v *verification, data map[string][]string) error {
  for key, values := range data {
    switch key {
    case "value":
      {
        v.Value = append(v.Value, values...)
      }
    case "comparison":
      {
        if len(values) > 1 {
          return fmt.Errorf("Comparison only accepts single value. Recieved '%s'.", values)
        }
        v.Comparison = values[0]
      }
    default:
      {
        return fmt.Errorf("Unknown verification key '%s'.", key)
      }
    }
  }
  return nil
}

func (p *parser) loadReports(reportData interface{}) (err error) {
  if reportData == nil {
    log.Tracef("No reports found, continuing.")
    return nil
  }
  var reports []*report
  log.Tracef("Parsing raw report data '%s'.", reportData)
  for name, value := range reportData.(map[string]interface{}) {
    report := newReport()
    report.Name = name
    for _, r := range value.([]interface{}) {
      report.Value = append(report.Value, r.(string))
    }
    reports = append(reports, &report)
  }
  p.Reports = reports
  return nil
}
