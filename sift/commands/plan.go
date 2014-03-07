package command

import (
  "encoding/json"
  log "github.com/cihub/seelog"
  "github.com/codegangsta/cli"
  "github.com/brettweavnet/sift/sift/engine"
  "github.com/brettweavnet/sift/sift/plan"
  "io/ioutil"
  "os"
  "strconv"
)

func planCommand(app *cli.App) cli.Command {
  c := cli.Command{
    Name:        "plan",
    Usage:       "sift plan -f PLAN_FILE -l LOG_LEVEL",
    Description: "Manage evaluations from plan.",
    Action: func(c *cli.Context) {
      args, err := readArgs(planFlags(), c)
      setLogLevel(args["log"])
      defer log.Flush()
      if err != nil {
        log.Critical(err)
        os.Exit(1)
      }
      file := args["file"]
      if file == "" {
        log.Critical("File required.")
        os.Exit(1)
      }
      log.Debugf("Loading plan from file '%s'.", file)

      data, err := ioutil.ReadFile(file)
      if err != nil {
        log.Critical("File required.")
        os.Exit(1)
      }
      p := plan.NewPlan()
      if err := p.LoadJSON(data); err != nil {
        log.Critical(err)
        os.Exit(1)
      }
      evaluations, err := p.Evaluations()
      if err != nil {
        log.Critical(err)
        os.Exit(1)
      }
      evaluationsJSON, err := json.MarshalIndent(evaluations, "", "  ")
      if err != nil {
        log.Critical(err)
        os.Exit(1)
      }
      if args["evaluations"] == "true" {
        log.Infof("Generated Evaluations JSON: \n%s\n", evaluationsJSON)
        return
      }
      e := engine.NewEngine()
      e.LoadEvaluationsFromJSON(evaluationsJSON)
      if err != nil {
        log.Critical(err)
        os.Exit(1)
      }
      timeout, err := strconv.Atoi(args["timeout"])
      if err != nil {
        log.Critical(err)
        os.Exit(1)
      }
      rc := engine.NewRunConfig(timeout, 300)
      if err := e.Execute(rc); err != nil {
        log.Critical(err)
        os.Exit(1)
      }
      run := e.Run()
      runJSON, err := json.MarshalIndent(run, "", "  ")
      if err != nil {
        log.Critical(err)
        os.Exit(1)
      }
      log.Tracef("Run JSON results:\n%s", runJSON)
      summary := engine.NewSummary(run)
      summary.Output()
      log.Infof("Completed successfully, exiting.")
    },
    Flags: planFlags(),
  }
  return c
}

func planFlags() []cli.Flag {
  f := []cli.Flag{
    cli.BoolFlag{"evaluations, e", "Display evaluations JSON, do not execute."},
    cli.StringFlag{"file, f", "", "Plan file to load"},
    cli.StringFlag{"log, l", "info", "Set log level"},
    cli.IntFlag{"timeout, t", 300, "Seconds to wait for state to load before timeout"},
  }
  return f
}
