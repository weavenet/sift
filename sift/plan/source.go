package plan

type source struct {
  Arguments map[string]map[string][]string `json:"arguments"`
}

func newSource() *source {
  return &source{Arguments: map[string]map[string][]string{}}
}
