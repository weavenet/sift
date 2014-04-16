package engine

import (
  log "github.com/cihub/seelog"
)

type endpointManager struct {
  cache     *cache
  runConfig runConfig
}

func newEndpointManager(c *cache, rc runConfig) endpointManager {
  return endpointManager{cache: c, runConfig: rc}
}

func (s endpointManager) setContextsEndpoints() error {
  cash := s.cache

  log.Infof("Starting source process manager.")

  for _, cntxt := range cash.Contexts {
    tgt := cntxt.targetAccount()

    if s.runConfig.overrideEndpointSet() {
      ep := s.runConfig.overrideEndpoint
      log.Debugf("Override set, using '%s'.", ep)
      cntxt.setSource(ep)
      continue
    }

    if !s.cache.sourceEndpointSet(tgt) {
      ep := newEndpoint(tgt, s.runConfig.waitForEndpoint)
      if err := ep.start(); err != nil {
        return err
      }

      cash.setSourceEndpoint(tgt, ep)

      if err := ep.waitForEndpoint(); err != nil {
        return err
      }
    }
    url := cash.getSourceEndpoint(tgt).url()
    cntxt.setSource(url)
  }
  return nil
}

func (s endpointManager) stopEndpointProcesses() error {
  log.Debugf("Stopping all endpoint processes.")
  for _, ep := range s.cache.sourceEndpoints {
    if err := ep.kill(); err != nil {
      return err
    }
  }
  return nil
}
