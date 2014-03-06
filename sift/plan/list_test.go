package plan

import (
  "reflect"
  "testing"
)

func TestList(t *testing.T) {
  l := newList()
  e1 := entry{Id: "user1", Tags: map[string]string{"preprod": "something_else"}}
  e2 := entry{Id: "user2", Tags: map[string]string{"prod": "admin1"}}
  l.Entries = map[string]entry{"user1": e1, "user2": e2}

  if !reflect.DeepEqual(l.all(), []string{"user1", "user2"}) {
    t.Errorf("List all failed.")
  }
  if !reflect.DeepEqual(l.sub("preprod"), []string{"something_else", "user2"}) {
    t.Errorf("List sub failed.")
  }
  if !reflect.DeepEqual(l.sub("blah"), []string{"user1", "user2"}) {
    t.Errorf("List sub failed.")
  }
  if !reflect.DeepEqual(l.only("prod"), []string{"admin1"}) {
    t.Errorf("List only failed.")
  }
}
