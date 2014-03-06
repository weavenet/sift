package state

import (
  "fmt"
  "reflect"
  "testing"
)

func TestState(t *testing.T) {
  data := map[string][]string{}
  data["foo"] = []string{"i-01", "i-02"}
  s := NewState("test1", data)

  if s.String() != "test1 : foo=i-01,i-02" {
    t.Fatalf("Error converting to string, received '%s'.", s.String())
  }

  val, err := s.Get("foo")
  if !reflect.DeepEqual(val, []string{"i-01", "i-02"}) {
    t.Fatalf("Error getting state data.")
  }
  if !reflect.DeepEqual(err, nil) {
    t.Fatalf("Error getting err %s", err.Error())
  }

  _, err = s.Get("not_here")
  if !reflect.DeepEqual(err, fmt.Errorf("Value for 'not_here' not set.")) {
    t.Fatalf("Error getting err %s", err.Error())
  }
}
