package util

import "testing"

func TestCreateUUID(t *testing.T) {
  uuid, _ := CreateUUID()
  if len(uuid) != 32 {
    t.Errorf("Error creating UUID.")
  }
}
