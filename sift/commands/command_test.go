package command

import (
  "github.com/codegangsta/cli"
  "testing"
)

type flagTestCase struct {
  Flag cli.Flag
  Name string
}

var flagTestCases = []flagTestCase{
  {cli.StringFlag{"file, f", "", "file to upload"}, "file"},
  {cli.StringFlag{"name, n", "", "name of file"}, "name"},
  {cli.StringFlag{"with-dash, w", "", "yes a dash"}, "with-dash"},
}

func TestFlagName(t *testing.T) {
  for _, tc := range flagTestCases {
    pn := parseStringFlagName(tc.Flag)
    if pn != tc.Name {
      t.Errorf("Flag name '%s' does not parse to '%s'.", pn, tc.Name)
    }
  }
}
