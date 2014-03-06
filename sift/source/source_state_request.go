package source

type stateRequestBody struct {
  Credentials map[string]string `json:"credentials"`
  Arguments   map[string]string `json:"arguments"`
  ParentIds   []string          `json:"parent_ids"`
}

func newStateRequestBody(creds map[string]string, args map[string]string, parentIds []string) stateRequestBody {
  return stateRequestBody{Credentials: creds, Arguments: args, ParentIds: parentIds}
}
