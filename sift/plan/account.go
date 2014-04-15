package plan

import (
  "fmt"
  log "github.com/cihub/seelog"
  "os"
  "strings"
)

type account struct {
  Credentials map[string]string `json:"credentials"`
  Scope       []string          `json:"scope"`
}

func newAccount() *account {
  return &account{}
}

func (a *account) convertEnvVarsToCredentials() error {
  for _, cred := range a.Credentials {
    if strings.HasPrefix(cred, "$") {
      credEnv := os.Getenv(strings.Trim(cred, "$"))
      if credEnv != "" {
        log.Debugf("Loading credentials '%s' from env variable '%s'.", credEnv, cred)
        a.Credentials[cred] = credEnv
      } else {
        return fmt.Errorf("Env variable '%s' not set.", cred)
      }
    }
  }
  return nil
}
