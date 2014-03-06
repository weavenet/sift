package state

import (
  "fmt"
  "strings"
)

type State struct {
  Id       string              `json:"id"`
  ParentId string              `json:"parent_id"`
  Data     map[string][]string `json:"data"`
}

func NewState(i string, d map[string][]string) *State {
  return &State{Id: i, Data: d}
}

func (s State) String() string {
  ss := fmt.Sprintf("%s :", s.Id)
  for i, values := range s.Data {
    ss += fmt.Sprintf(" %s=%s", i, strings.Join(values, ","))
  }
  return ss
}

func (s State) Get(id string) ([]string, error) {
  if val, ok := s.Data[id]; ok {
    return val, nil
  }
  return []string{}, fmt.Errorf("Value for '%s' not set.", id)
}
