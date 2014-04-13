package engine

import (
  "fmt"
  "github.com/brettweavnet/sift/sift/source"
  "github.com/brettweavnet/sift/sift/state"
  log "github.com/cihub/seelog"
  "reflect"
  "sort"
)

type context struct {
  id         string
  ParentId   string        `json:"parent_id"`
  Account    account       `json:"account"`
  Provider   provider      `json:"provider"`
  Collection collection    `json:"collection"`
  Source     source.Source `json:"source"`
  State      []state.State `json:"state"`
}

func (c context) Id() string {
  return c.id
}

func (c context) String() string {
  return c.targetAccount()
}

func (c context) targetAccount() string {
  return c.Account.Name
}

func newContext(p parser) *context {
  c := context{}
  c.Account = p.Account
  c.Collection = p.Collection
  c.Provider = p.Provider
  c.id = p.Id
  log.Tracef("Creating new context with account '%s' provider '%s' collection '%s' from '%s'.", c.Account.Name, c.Provider.Name, c.Collection.Name)
  return &c
}

func (c context) hasParent() bool {
  return c.Collection.hasParent()
}

func (c context) parent() collection {
  return c.Collection.parent()
}

func (c *context) setSource(sourceURL string) {
  contextSource := source.NewSource(c.Account.Name, c.Provider.Name, c.Collection.Name, sourceURL)
  c.Source = *contextSource
}

func (c *context) loadState(cash *cache) error {
  // Recursively ensure parent and it's parents loaded
  if c.hasParent() {
    log.Tracef("Context '%s' has parent '%s', loading.", c.id, c.ParentId)
    parentContext, ok := cash.GetContext(c.ParentId)
    if !ok {
      return fmt.Errorf("Parent context not found.")
    }
    parentContext.loadState(cash)
  }

  parentIds, err := c.parentIds(cash)
  if err != nil {
    return err
  }

  return c.setStateFromSource(parentIds)
}

func (c *context) setStateFromSource(parentIds []string) error {
  ts, err := c.Source.State(c.Account.credentials, c.Provider.Arguments, parentIds)
  if err != nil {
    return err
  }
  c.State = ts
  return nil
}

func (c context) parentIds(cash *cache) ([]string, error) {
  ids := make([]string, 0)
  if c.hasParent() {
    if parentContext, ok := cash.GetContext(c.ParentId); ok {
      for _, state := range parentContext.State {
        ids = append(ids, state.Id)
      }
    } else {
      return ids, fmt.Errorf("Parent '%s' not found in cache.", c.ParentId)
    }
  }
  return ids, nil
}

func (c *context) validate() error {
  if err := c.validateProvider(); err != nil {
    return err
  }
  if err := c.validateAccount(); err != nil {
    return err
  }
  return nil
}

func (c *context) validateProvider() error {
  log.Debugf("Validating context '%s'.", c.id)
  providerArguments := []string{}
  for k, _ := range c.Provider.Arguments {
    providerArguments = append(providerArguments, k)
  }
  requiredArgs, err := c.Source.ProviderArguments()
  if err != nil {
    return err
  }

  sort.Strings(providerArguments)
  sort.Strings(requiredArgs)
  if !reflect.DeepEqual(requiredArgs, providerArguments) {
    return fmt.Errorf("Provider '%s' required arguments '%s' do not match given arguments '%s'.", c.Provider.Name, requiredArgs, providerArguments)
  }
  return nil
}

func (c *context) validateAccount() error {
  accountCredentials := []string{}
  for k, _ := range c.Account.credentials {
    accountCredentials = append(accountCredentials, k)
  }
  sort.Strings(accountCredentials)
  requiredCreds, err := c.Source.Credentials()
  if err != nil {
    return err
  }
  sort.Strings(requiredCreds)
  if !reflect.DeepEqual(requiredCreds, accountCredentials) {
    return fmt.Errorf("Account '%s' required credentials '%s' do not match given credentials '%s'.", c.Account.Name, requiredCreds, accountCredentials)
  }
  return nil
}
