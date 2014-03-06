package util

import (
  "crypto/rand"
  "encoding/hex"
)

func CreateUUID() (string, error) {
  uuid := make([]byte, 16)
  n, err := rand.Read(uuid)
  if n != len(uuid) || err != nil {
    return "", err
  }
  return hex.EncodeToString(uuid), nil
}
