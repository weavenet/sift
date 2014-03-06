package plan

import (
  "testing"
)

var validFileTCs = []struct {
  file   string
  result bool
}{
  {"123.json", true},
  {"ABC123.json", true},
  {"abc123.json", true},
  {"a-B_c-1.json", true},
  {".lksjdf.json", false},
  {".lksjdf", false},
  {"lksjdf", false},
  {"lksjdf", false},
}

func TestValidFile(t *testing.T) {
  for _, tc := range validFileTCs {
    if validFile(tc.file) != tc.result {
      t.Fatalf("Error testin file validation for %s", tc.file)
    }
  }
}
