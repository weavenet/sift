package engine

import (
  log "github.com/cihub/seelog"
)

const (
  CONTEXT_PENDING = iota
  CONTEXT_WAITING_FOR_PARENT
  CONTEXT_LOAD_IN_PROGRESS
  CONTEXT_LOADED
)

type cache struct {
  contexts        map[string]*context `json:"context"`
  status          map[string]int
  sourceEndpoints map[string]endpoint
}

func newCache() *cache {
  return &cache{contexts: map[string]*context{},
    status:          map[string]int{},
    sourceEndpoints: map[string]endpoint{}}
}

func (c cache) Contexts() map[string]*context {
  return c.contexts
}

func (c cache) GetContext(id string) (*context, bool) {
  val, ok := c.contexts[id]
  return val, ok
}

func (c cache) GetContextStatus(id string) int {
  return c.status[id]
}

func (c cache) SetContext(id string, cntxt *context) {
  c.contexts[id] = cntxt
}

func (c cache) SetContextStatus(id string, status int) {
  c.status[id] = status
}

func (c cache) SetContextStatusPending(id string) {
  c.status[id] = CONTEXT_PENDING
}

func (c cache) SetContextStatusWaitingForParent(id string) {
  c.status[id] = CONTEXT_WAITING_FOR_PARENT
}

func (c cache) SetContextStatusLoadInProgress(id string) {
  c.status[id] = CONTEXT_LOAD_IN_PROGRESS
}

func (c cache) SetContextStatusLoaded(id string) {
  c.status[id] = CONTEXT_LOADED
}

func (c cache) NumberContextStatusPending() int {
  return c.numberContextStatus(CONTEXT_PENDING)
}

func (c cache) NumberContextStatusWaitingForParent() int {
  return c.numberContextStatus(CONTEXT_WAITING_FOR_PARENT)
}

func (c cache) NumberContextStatusLoadInProgress() int {
  return c.numberContextStatus(CONTEXT_LOAD_IN_PROGRESS)
}

func (c cache) NumberContextStatusLoaded() int {
  return c.numberContextStatus(CONTEXT_LOADED)
}

func (c cache) numberContextStatus(s int) int {
  count := 0
  for _, v := range c.status {
    if v == s {
      count = count + 1
    }
  }
  return count
}

func (c cache) ContextStatusLoaded(id string) bool {
  if c.GetContextStatus(id) == CONTEXT_LOADED {
    return true
  }
  return false
}

func (c *cache) addContextToCache(cntxt *context) {
  if c.ContextStatusLoaded(cntxt.id) {
    log.Tracef("Context in cache, not re-loading.")
  } else {
    log.Tracef("Context not in cache, loading.")
    c.SetContext(cntxt.id, cntxt)
    c.SetContextStatusPending(cntxt.id)
  }
}

func (c *cache) setSourceEndpoint(name string, ep endpoint) {
  c.sourceEndpoints[name] = ep
}

func (c cache) getSourceEndpoint(name string) endpoint {
  return c.sourceEndpoints[name]
}

func (c cache) sourceEndpointSet(name string) bool {
  _, ok := c.sourceEndpoints[name]
  return ok
}
