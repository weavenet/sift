package target

import (
  "fmt"
  "github.com/brettweavnet/sift/sift/state"
  "reflect"
  "testing"
)

var reportTCs = []struct {
  name   string
  states []state.State
  value  []string
  pass   bool
}{
  {"equals", []state.State{*s1, *s2}, []string{"1", "2"}, true},
  {"equals", []state.State{*s1, *s2}, []string{"2", "1"}, true},
  {"equals", []state.State{*s1, *s2}, []string{"2", "1", "3"}, false},
  {"less_than", []state.State{*s1, *s2}, []string{"3"}, true},
  {"less_than", []state.State{*s1, *s2}, []string{"2"}, false},
  {"greater_than", []state.State{*s1, *s2}, []string{"0"}, true},
  {"greater_than", []state.State{*s1, *s2}, []string{"10"}, false},
}

func TestReports(t *testing.T) {
  for _, tc := range reportTCs {
    result, err := runReport(tc.name, tc.value, tc.states)
    if err != nil {
      t.Fatalf("Error evaluating equal states.")
    }
    if result[tc.name] != tc.pass {
      t.Fatalf("Error evaluating '%s' states.", tc.name)
    }
  }
}

var valueTCs = []struct {
  name       string
  comparison string
  value      []string
  states     []state.State
  passing    map[string]bool
  err        error
}{
  {"image_id", "equals", []string{"img-01"}, states, map[string]bool{"1": true, "2": false}, nil},
  {"groups", "equals", []string{"db", "admin"}, states, map[string]bool{"1": false, "2": true}, nil},
  {"not_here", "equals", []string{"db", "admin"}, states, map[string]bool{}, fmt.Errorf("Value for 'not_here' not set.")},
  {"groups", "includes", []string{"admin"}, states, map[string]bool{"1": true, "2": true}, nil},
  {"groups", "includes", []string{"db"}, states, map[string]bool{"1": false, "2": true}, nil},
  {"groups", "includes", []string{"blah"}, states, map[string]bool{"1": false, "2": false}, nil},
  {"not_here", "includes", []string{"db", "admin"}, states, map[string]bool{}, fmt.Errorf("Value for 'not_here' not set.")},
  {"groups", "excludes", []string{"admin"}, states, map[string]bool{"1": false, "2": false}, nil},
  {"groups", "excludes", []string{"db"}, states, map[string]bool{"1": true, "2": false}, nil},
  {"groups", "excludes", []string{"web"}, states, map[string]bool{"1": false, "2": true}, nil},
  {"not_here", "excludes", []string{"true"}, states, map[string]bool{}, fmt.Errorf("Value for 'not_here' not set.")},
  {"image_id", "set", []string{"true"}, states, map[string]bool{"1": true, "2": true}, nil},
  {"image_id", "set", []string{"false"}, states, map[string]bool{"1": false, "2": false}, nil},
  {"not_here", "set", []string{"true"}, states, map[string]bool{}, fmt.Errorf("Value for 'not_here' not set.")},
  {"groups", "within", []string{"one", "db", "admin", "web", "another"}, states, map[string]bool{"1": true, "2": true}, nil},
  {"groups", "within", []string{"db", "admin", "other"}, states, map[string]bool{"1": false, "2": true}, nil},
  {"not_here", "within", []string{"db", "admin"}, states, map[string]bool{}, fmt.Errorf("Value for 'not_here' not set.")},
}

func TestValues(t *testing.T) {
  for _, tc := range valueTCs {
    result, err := runComparison(tc.name, tc.comparison, tc.value, tc.states)
    if !reflect.DeepEqual(err, tc.err) {
      t.Fatalf("Evaluating equal states got wrong error '%s' expect '%+v'.", err.Error(), tc.err)
    }
    if !reflect.DeepEqual(result, tc.passing) {
      t.Fatalf("Error evaluating '%s' states '%s'. Expect '%+v' Got '%+v'", tc.name, tc.comparison, tc.passing, result)
    }
  }
}
