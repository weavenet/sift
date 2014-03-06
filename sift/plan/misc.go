package plan

import (
  "os"
  "strings"
)

func dirExists(path string) bool {
  _, err := os.Stat(path)
  if err == nil {
    return true
  }
  if os.IsNotExist(err) {
    return false
  }
  return false
}

func fileToName(file string) string {
  return strings.Split(file, ".json")[0]
}
