package plan

type list struct {
  Entries map[string]entry `json:"entries"`
}

type entry struct {
  Id   string            `json:"id"`
  Tags map[string]string `json:"tags"`
}

func newList() *list {
  return &list{Entries: map[string]entry{}}
}

func (l list) all() (r []string) {
  for _, e := range l.Entries {
    r = append(r, e.Id)
  }
  return r
}

func (l list) sub(tag string) (r []string) {
  for _, e := range l.Entries {
    if val, ok := e.Tags[tag]; ok {
      r = append(r, val)
    } else {
      r = append(r, e.Id)
    }
  }
  return r
}

func (l list) only(tag string) (r []string) {
  for _, e := range l.Entries {
    if val, ok := e.Tags[tag]; ok {
      r = append(r, val)
    }
  }
  return r
}
