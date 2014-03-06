package target

import (
  "github.com/siftproject/sift/sift/state"
)

var data1 = map[string][]string{"image_id": []string{"img-01"}, "groups": []string{"web", "admin"}}
var data2 = map[string][]string{"image_id": []string{"img-02"}, "groups": []string{"db", "admin"}}
var s1 = state.NewState("1", data1)
var s2 = state.NewState("2", data2)
var states = []state.State{*s1, *s2}
