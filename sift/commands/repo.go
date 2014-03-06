package command

import (
  "encoding/json"
  log "github.com/cihub/seelog"
  "github.com/codegangsta/cli"
  "github.com/siftproject/sift/sift/engine"
  "github.com/siftproject/sift/sift/plan"
  "os"
  "strconv"
  "time"
)

func repoCommand(app *cli.App) cli.Command {
  c := cli.Command{
    Name:        "repo",
    Usage:       "sift repo -d DIRECTORY -l LOG_LEVEL",
    Description: "Manage evaluations from repo.",
    Action: func(c *cli.Context) {
      args, err := readArgs(repoFlags(), c)
      setLogLevel(args["log"])
      defer log.Flush()
      if err != nil {
        log.Critical(err)
        os.Exit(1)
      }
      repoDirectory := args["directory"]
      if repoDirectory == "" {
        log.Critical("Repo directory required.")
        os.Exit(1)
      }
      startTime := time.Now()
      p := plan.NewPlan()
      if err := p.LoadRepo(repoDirectory); err != nil {
        log.Critical(err)
        os.Exit(1)
      }
      if args["plan"] == "true" {
        data, err := json.MarshalIndent(p, "", "  ")
        if err != nil {
          log.Critical(err)
          os.Exit(1)
        }
        log.Infof("Generated Plan JSON: \n\n%s\n", data)
        return
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
        log.Infof("Generated Evaluations JSON: \n\n%s\n", evaluationsJSON)
        return
      }
      e := engine.NewEngine()
      if err := e.LoadEvaluationsFromJSON(evaluationsJSON); err != nil {
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
      duration := time.Since(startTime)
      log.Infof("Completed successfully (%3.2f seconds), exiting.", duration.Seconds())
    },
    Flags: repoFlags(),
  }
  return c
}

func repoFlags() []cli.Flag {
  f := []cli.Flag{
    cli.StringFlag{"directory, d", "", "Repo directory to load"},
    cli.BoolFlag{"evaluations, e", "Display evaluations JSON, do not execute."},
    cli.StringFlag{"log, l", "info", "Set log level"},
    cli.BoolFlag{"plan, p", "Display plan JSON, do not execute."},
    cli.IntFlag{"timeout, t", 300, "Seconds to wait for state to load before timeout"},
  }
  return f
}
