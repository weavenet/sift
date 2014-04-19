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
  Name        string            `json:"name"`
}

func newAccount() *account {
  return &account{}
}

func (a *account) convertEnvVarsToCredentials() error {
  for name, value := range a.Credentials {
    if strings.HasPrefix(value, "$") {
      envVar := strings.Trim(value, "$")
      credEnv := os.Getenv(envVar)
      if credEnv != "" {
        log.Debugf("Loading credentials '%s' from env variable '%s'.", name, envVar)
        a.Credentials[name] = credEnv
      } else {
        return fmt.Errorf("Env variable '%s' not set.", envVar)
      }
    }
  }
  return nil
}
