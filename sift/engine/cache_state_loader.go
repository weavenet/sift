package engine

import (
  "fmt"
  log "github.com/cihub/seelog"
  "math/rand"
  "time"
)

type cacheStateLoader struct {
  cache           *cache
  timeout         time.Duration
  statusFrequency time.Duration
  waitForParent   time.Duration
}

func newCacheStateLoader(c *cache, rc runConfig) *cacheStateLoader {
  csl := cacheStateLoader{cache: c}
  csl.timeout = rc.timeout
  csl.statusFrequency = rc.statusFrequency
  csl.waitForParent = rc.waitForParent
  return &csl
}

var random *rand.Rand = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func (c cacheStateLoader) loadStates() (err error) {
  // Channel of contexts to be processed
  // contexts waiting on parents are put back on the end of the queue
  queueChan := make(chan *context, len(c.cache.Contexts()))
  go func() {
    for _, c := range c.cache.Contexts() {
      queueChan <- c
    }
  }()

  // Timeout will exit after duration
  log.Debugf("Timing out after %d seconds if loading states does not complete.", int(c.timeout.Seconds()))
  timeoutChan := make(chan bool, 1)
  go func() {
    time.Sleep(c.timeout)
    timeoutChan <- true
  }()

  // Send periodic status message
  statusChan := make(chan bool, 1)
  go func() {
    for {
      time.Sleep(c.statusFrequency)
      statusChan <- true
    }
  }()

  return c.processQueue(queueChan, timeoutChan, statusChan)
}

func (c cacheStateLoader) processQueue(queueChan chan *context, timeoutChan chan bool, statusChan chan bool) error {
  contextCount := len(c.cache.Contexts())
  resultChan := make(chan error, contextCount)
  resultCount := 0

  for {
    select {
    case cntxt := <-queueChan:
      go c.loadContextState(cntxt, c.cache, resultChan, queueChan)
    case err := <-resultChan:
      if err != nil {
        return err
      } else {
        resultCount = resultCount + 1
      }

      if resultCount == contextCount {
        log.Infof("Loading states completed successfully.")
        return nil
      }
    case <-statusChan:
      w := c.cache.NumberContextStatusWaitingForParent()
      i := c.cache.NumberContextStatusLoadInProgress()
      l := c.cache.NumberContextStatusLoaded()
      log.Infof("Waiting for cache to load. %d waiting, %d in-progress, %d loaded.", w, i, l)
    case <-timeoutChan:
      return fmt.Errorf("Timeout loading state.")
    }
  }
  return nil
}

func (c cacheStateLoader) loadContextState(cntxt *context, cash *cache, resultChan chan error, queueChan chan *context) {
  log.Debugf("Loading state for context with ID '%s'.", cntxt.id)

  if cntxt.hasParent() {
    log.Tracef("Context '%s' has parent '%s'.", cntxt.id, cntxt.ParentId)

    if cash.ContextStatusLoaded(cntxt.ParentId) {
      log.Tracef("Parent '%s' has been loaded.", cntxt.ParentId)
    } else {
      cash.SetContextStatusWaitingForParent(cntxt.id)
      log.Debugf("Parent '%s' has not been loaded, '%s' waiting.", cntxt.ParentId, cntxt.id)
      time.Sleep(durationJitter(c.waitForParent, random))
      queueChan <- cntxt
      return
    }
  }

  log.Debugf("Loading state '%s'.", cntxt.id)
  cash.SetContextStatusLoadInProgress(cntxt.id)
  if err := cntxt.loadState(cash); err != nil {
    log.Debugf("Error loading context '%s'.", cntxt.id)
    resultChan <- err
    return
  }
  log.Debugf("Succesfully loaded state '%s'.", cntxt.id)
  cash.SetContext(cntxt.id, cntxt)
  cash.SetContextStatusLoaded(cntxt.id)
  resultChan <- nil
}

func durationJitter(d time.Duration, r *rand.Rand) time.Duration {
  if d == 0 {
    return 0
  }
  log.Debugf("Waiting %s for parent to load.", d+time.Duration(r.Int63n(2*int64(d))))
  return d + time.Duration(r.Int63n(2*int64(d)))
}
