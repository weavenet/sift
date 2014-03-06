package engine

import (
  "time"
)

type runConfig struct {
  timeout          time.Duration
  statusFrequency  time.Duration
  waitForParent    time.Duration
  waitForEndpoint  time.Duration
  overrideEndpoint string
}

func NewRunConfig(timeout int, wait int) runConfig {
  rc := runConfig{}
  rc.setTimeout(timeout)
  rc.setWaitForParent(wait)
  rc.setWaitForEndpoint(1000)
  rc.setStatusFrequency(timeout / 30)
  rc.overrideEndpoint = ""
  return rc
}

// Send all sources to a single endpoint
// Primarily used for tests
func (r *runConfig) SetOverrideEndpoint(oep string) {
  r.overrideEndpoint = oep
}

func (r *runConfig) setTimeout(seconds int) {
  r.timeout = time.Duration(seconds) * time.Second
}

func (r *runConfig) setStatusFrequency(seconds int) {
  if seconds < 1 {
    seconds = 1
  }
  r.statusFrequency = time.Duration(seconds) * time.Second
}

func (r *runConfig) setWaitForParent(milliseconds int) {
  r.waitForParent = time.Duration(milliseconds) * (time.Nanosecond * 1e6)
}

func (r *runConfig) setWaitForEndpoint(milliseconds int) {
  r.waitForEndpoint = time.Duration(milliseconds) * (time.Nanosecond * 1e6)
}

func (r runConfig) overrideEndpointSet() bool {
  if r.overrideEndpoint == "" {
    return false
  }
  return true
}
