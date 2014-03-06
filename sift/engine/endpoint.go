package engine

import (
  "bytes"
  "fmt"
  log "github.com/cihub/seelog"
  "net"
  "os"
  "os/exec"
  "path/filepath"
  "strconv"
  "time"
)

type endpoint struct {
  name string
  port string
  wait time.Duration
  cmd  *exec.Cmd
}

func newEndpoint(name string, wait time.Duration) endpoint {
  return endpoint{name: name, wait: wait}
}

func (e endpoint) url() string {
  return "http://localhost:" + e.port
}

func (e *endpoint) start() error {
  log.Debugf("Starting endpoint '%s'.", e.name)

  if err := e.setPort(); err != nil {
    return err
  }
  if err := e.executeProcess(); err != nil {
    return err
  }

  return nil
}

func (e *endpoint) setPort() error {
  for p := 32768; p <= 65523; p++ {
    port := strconv.Itoa(p)
    log.Tracef("Testing port '%s'.", port)
    _, err := net.Dial("tcp", "localhost:"+port)

    // if port listening continue
    if err == nil {
      log.Tracef("Port '%s' listening, continuing.", port)
      continue
    }

    log.Tracef("Port '%s' available.", port)
    e.port = port
    return nil
  }
  return fmt.Errorf("Unable to find open port.")
}

func (e *endpoint) executeProcess() error {
  workingDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
  if err != nil {
    return err
  }
  cmdName := workingDir + "/sift-source-" + e.name
  cmd := exec.Command(cmdName, "-port="+e.port)
  var errBuf, outBuf bytes.Buffer
  cmd.Stdout = &outBuf
  cmd.Stderr = &errBuf

  log.Debugf("Starting '%s' endpoint server with command '%s -p %s' '.", e.name, cmdName, e.port)
  if err := cmd.Start(); err != nil {
    return err
  }
  e.cmd = cmd
  return nil
}

func (e endpoint) kill() error {
  log.Tracef("Killing process for '%s' endpoint.", e.name)
  if err := e.cmd.Process.Kill(); err != nil {
    return err
  }
  return nil
}

func (e endpoint) waitForEndpoint() error {
  time.Sleep(e.wait)

  for i := 1; i < 10; i++ {
    log.Debugf("Waiting for '%s' endpoint on port '%s'.", e.name, e.port)
    _, err := net.Dial("tcp", "localhost:"+e.port)
    if err == nil {
      log.Infof("Endpoint '%s' online.", e.name)
      return nil
    }

    log.Infof("Waiting for '%s' endpoint. Sleeping...", e.name)
    time.Sleep(e.wait)
  }
  log.Errorf("Output of endpoint command: \n\n Stdout: '%s' \n\n Stderr: '%s'\n", e.cmd.Stdout, e.cmd.Stderr)
  return fmt.Errorf("Timeout waiting for endpoint '%s' to start.", e.name)
}
