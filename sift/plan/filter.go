package plan

type planFilter struct {
  Include    interface{}            `json:"include"`
  Exclude    interface{}            `json:"exclude"`
  Attributes map[string]interface{} `json:"attributes"`
}

type filter struct {
  Include    []string            `json:"include"`
  Exclude    []string            `json:"exclude"`
  Attributes map[string][]string `json:"attributes"`
}

func newFilter() *filter {
  return &filter{
    Include:    []string{},
    Exclude:    []string{},
    Attributes: map[string][]string{},
  }
}
