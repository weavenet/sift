package plan

import (
  "encoding/json"
  "fmt"
  log "github.com/cihub/seelog"
  "io/ioutil"
  "regexp"
)

type repo struct {
  path string
}

func newRepo(path string) *repo {
  return &repo{path: path}
}

func (r *repo) loadAccounts(p *Plan) error {
  path := r.accountsPath()
  if !dirExists(path) {
    return fmt.Errorf("Accounts directory does not exist in repo.")
  }
  accountDirs, err := ioutil.ReadDir(path)
  if err != nil {
    return err
  }
  for _, accountDir := range accountDirs {
    if accountDir.IsDir() {
      err := r.loadAccountsFromDir(accountDir.Name(), p)
      if err != nil {
        return err
      }
    }
  }
  if len(p.Accounts) == 0 {
    return fmt.Errorf("No accounts found in repo.")
  }
  return nil
}

func (r *repo) loadAccountsFromDir(path string, p *Plan) error {
  fp := r.accountsPath() + "/" + path
  log.Debugf("Loading accounts from directory '%s'.", fp)
  accountFiles, err := ioutil.ReadDir(fp)
  if err != nil {
    return err
  }
  p.Accounts[path] = map[string]account{}
  for _, accountFile := range accountFiles {
    fp := fp + "/" + accountFile.Name()

    if !validFile(accountFile.Name()) {
      log.Warnf("Not processing invalid account file '%s'.", accountFile.Name())
      continue
    }

    log.Debugf("Loading account file '%s'.", fp)
    file, err := ioutil.ReadFile(fp)
    if err != nil {
      return err
    }
    a := newAccount()
    if err := json.Unmarshal(file, &a); err != nil {
      return err
    }
    if err := a.convertEnvVarsToCredentials(); err != nil {
      return err
    }
    if len(a.Scope) == 0 {
      log.Tracef("Scope not set, using 'default'.")
      a.Scope = []string{"default"}
    }
    name := fileToName(accountFile.Name())
    log.Debugf("Loaded account '%s'.", name)
    p.Accounts[path][name] = *a
  }
  return nil
}

func (r *repo) loadLists(p *Plan) error {
  if !dirExists(r.listsPath()) {
    log.Debugf("Lists directory does not exist in repo, skipping.")
    return nil
  }
  listDirs, err := ioutil.ReadDir(r.listsPath())
  if err != nil {
    return err
  }
  for _, listDir := range listDirs {
    if listDir.IsDir() {
      err := r.loadListsFromDir(listDir.Name(), p)
      if err != nil {
        return err
      }
    }
  }
  return nil
}

func (r *repo) loadListsFromDir(path string, p *Plan) error {
  lp := r.listsPath() + "/" + path
  log.Debugf("Loading list from '%s'.", lp)
  listFiles, err := ioutil.ReadDir(lp)
  if err != nil {
    return err
  }
  lst := newList()

  for _, entryFile := range listFiles {
    entryFilePath := lp + "/" + entryFile.Name()

    if !validFile(entryFile.Name()) {
      log.Warnf("Not processing invalid list entry file '%s'.", entryFile.Name())
      continue
    }

    log.Debugf("Loading entry from file '%s'.", entryFilePath)
    file, err := ioutil.ReadFile(entryFilePath)
    if err != nil {
      return err
    }
    e := entry{}
    json.Unmarshal(file, &e.Tags)

    // id is a special tag
    // It is used to over-ride the entry id if presnet
    // It is then removed from the tags list
    if val, ok := e.Tags["id"]; ok {
      e.Id = val
      delete(e.Tags, "id")
    } else {
      e.Id = fileToName(entryFile.Name())
    }
    log.Debugf("Loaded entry '%s'.", e.Id)
    lst.Entries[e.Id] = e
  }
  p.Lists[path] = *lst
  return nil
}

func (r *repo) loadSources(p *Plan) error {
  log.Debugf("Loading sources from '%s'.", r.sourcesPath())
  sourceDirs, err := ioutil.ReadDir(r.sourcesPath())
  if err != nil {
    return err
  }

  for _, sourceFile := range sourceDirs {
    sfp := r.sourcesPath() + "/" + sourceFile.Name()

    if !validFile(sourceFile.Name()) {
      log.Warnf("Not processing invalid source file '%s'.", sourceFile.Name())
      continue
    }

    log.Debugf("Loading source from '%s'.", sfp)
    file, err := ioutil.ReadFile(sfp)
    if err != nil {
      return err
    }
    s := newSource()
    if err := json.Unmarshal(file, &s.Arguments); err != nil {
      return err
    }
    name := fileToName(sourceFile.Name())
    log.Debugf("Loaded source '%s'.", name)
    p.Sources[name] = *s
  }
  return nil
}

func (r *repo) loadFilters(p *Plan) error {
  if !dirExists(r.filtersPath()) {
    log.Debugf("Filters directory does not exist in repo, skipping.")
    return nil
  }

  log.Debugf("Loading filters from '%s'.", r.filtersPath())
  filterFiles, err := ioutil.ReadDir(r.filtersPath())
  if err != nil {
    return err
  }

  for _, filterFile := range filterFiles {
    log.Debugf("Loading filters from file '%s'.", filterFile.Name())

    ffp := r.filtersPath() + "/" + filterFile.Name()

    if !validFile(filterFile.Name()) {
      log.Warnf("Not processing invalid filter file '%s'.", filterFile.Name())
      continue
    }

    log.Debugf("Loading filter from '%s'.", ffp)
    filters, err := ioutil.ReadFile(ffp)
    if err != nil {
      return err
    }

    fltrs := map[string]planFilter{}
    if err := json.Unmarshal(filters, &fltrs); err != nil {
      return err
    }
    for key, value := range fltrs {
      log.Debugf("Creating filter '%s'.", key)
      p.Filters[key] = value
    }
  }
  return nil
}

func (r *repo) loadPolicies(p *Plan) error {
  path := r.policiesPath()
  log.Debugf("Loading policies from '%s'.", path)
  if !dirExists(path) {
    return fmt.Errorf("Policies directory does not exist in repo.")
  }
  policyFiles, err := ioutil.ReadDir(r.policiesPath())
  if err != nil {
    return err
  }
  for _, policyFile := range policyFiles {
    if !validFile(policyFile.Name()) {
      log.Warnf("Not processing invalid policy file '%s'.", policyFile.Name())
      continue
    }

    name := fileToName(policyFile.Name())
    log.Debugf("Loading policies from file '%s'.", name)

    cfp := r.policiesPath() + "/" + policyFile.Name()
    file, err := ioutil.ReadFile(cfp)
    if err != nil {
      return err
    }
    policies := []policy{}
    if err := json.Unmarshal(file, &policies); err != nil {
      return err
    }
    p.Policies = append(p.Policies, policies...)
  }

  if len(p.Policies) == 0 {
    return fmt.Errorf("No policies found in repo.")
  }
  return nil
}

func validFile(fileName string) bool {
  r := regexp.MustCompile(`^[0-9a-zA-Z-_]+\.json$`)
  if r.MatchString(fileName) {
    return true
  }
  return false
}

func (r repo) accountsPath() string {
  return fmt.Sprintf("%s/accounts", r.path)
}

func (r repo) listsPath() string {
  return fmt.Sprintf("%s/lists", r.path)
}

func (r repo) sourcesPath() string {
  return fmt.Sprintf("%s/sources", r.path)
}

func (r repo) filtersPath() string {
  return fmt.Sprintf("%s/filters", r.path)
}

func (r repo) policiesPath() string {
  return fmt.Sprintf("%s/policies", r.path)
}
