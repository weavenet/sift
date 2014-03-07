package target

import (
  "fmt"
  "github.com/brettweavnet/sift/sift/state"
  "reflect"
  "testing"
)

var includeExcludeTCs = []struct {
  include []string
  exclude []string
  states  []state.State
  result  []state.State
}{
  {[]string{"1", "2"}, []string{}, states, []state.State{*s1, *s2}},
  {[]string{"1", "5"}, []string{}, states, []state.State{*s1}},
  {[]string{"3"}, []string{}, states, []state.State{}},
  {[]string{}, []string{"1"}, states, []state.State{*s2}},
  {[]string{}, []string{"8"}, states, []state.State{*s1, *s2}},
  {[]string{}, []string{"1", "2", "3"}, states, []state.State{}},
}

func TestFilterIncludeExclude(t *testing.T) {
  for count, tc := range includeExcludeTCs {
    t.Logf("Executing tc %d", count)
    results := filterIncludeExclude(tc.include, tc.exclude, tc.states)
    if !reflect.DeepEqual(tc.result, results) {
      t.Fatalf("Error testing include / exclude.")
    }
  }
}

var attributeTCs = []struct {
  attributes attributeList
  states     []state.State
  results    []state.State
  err        error
}{
  {newAttributeList("image_id", "equals", []string{"img-01"}), states, []state.State{*s1}, nil},
  {newAttributeList("image_id", "equals", []string{"img-01"}), states, []state.State{*s1}, nil},
  {newAttributeList("image_id", "within", []string{"img-01", "img-03"}), states, []state.State{*s1}, nil},
  {newAttributeList("image_id", "within", []string{"boom"}), states, []state.State{}, nil},
  {newAttributeList("groups", "includes", []string{"admin"}), states, []state.State{*s1, *s2}, nil},
  {newAttributeList("groups", "includes", []string{"not_here"}), states, []state.State{}, nil},
  {newAttributeList("image_id", "set", []string{"true"}), states, []state.State{*s1, *s2}, nil},
  {newAttributeList("image_id", "set", []string{"false"}), states, []state.State{}, nil},
  {newAttributeList("not_here", "set", []string{"false"}), states, []state.State{}, fmt.Errorf("Value for 'not_here' not set.")},
}

func TestFilterAttributes(t *testing.T) {
  for count, tc := range attributeTCs {
    t.Logf("Executing tc %d", count)
    results, err := filterAttributes(tc.attributes, tc.states)
    if !reflect.DeepEqual(err, tc.err) {
      t.Fatalf("Unexpected err '%s' expected '%s'", err, tc.err)
    }

    if !reflect.DeepEqual(tc.results, results) {
      t.Fatalf("Error testing attribute. Expected '%+v' got '%+v'", tc.results, results)
    }
  }
}
