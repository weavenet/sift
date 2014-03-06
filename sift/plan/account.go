package plan

type account struct {
  Credentials map[string]string `json:"credentials"`
  Scope       []string          `json:"scope"`
}

func newAccount() *account {
  return &account{}
}
