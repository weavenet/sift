package plan

type evaluation struct {
  Account       map[string]string              `json:"account"`
  Provider      map[string]string              `json:"provider"`
  Collection    map[string]string              `json:"collection"`
  Filters       map[string]filter              `json:"filters"`
  Verifications map[string]map[string][]string `json:"verifications"`
  Reports       map[string][]string            `json:"reports"`
}

func newEvaluation() *evaluation {
  return &evaluation{
    Account:       map[string]string{},
    Provider:      map[string]string{},
    Collection:    map[string]string{},
    Filters:       map[string]filter{},
    Verifications: map[string]map[string][]string{},
    Reports:       map[string][]string{},
  }
}
