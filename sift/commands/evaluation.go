package command

import (
  "encoding/json"
  "github.com/brettweavnet/sift/sift/engine"
  log "github.com/cihub/seelog"
  "github.com/codegangsta/cli"
  "os"
  "strconv"
)

func evaluationCommand(app *cli.App) cli.Command {
  c := cli.Command{
    Name:        "evaluation",
    Usage:       "sift evaluation -f FILE -l LOG_LEVEL",
    Description: "Manage sift evaluations",
    Action: func(c *cli.Context) {
      args, err := readArgs(evaluationFlags(), c)
      setLogLevel(args["log-level"])
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

      e := engine.NewEngine()
      e.LoadEvaluationsFromFile(file)
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
      results := e.Evaluations()
      data, err := json.MarshalIndent(results, "", "  ")
      if err != nil {
        log.Critical(err)
        os.Exit(1)
      }
      log.Infof("Evaluation Results:\n%s", data)
      log.Infof("Completed successfully, exiting.")
    },
    Flags: evaluationFlags(),
  }
  return c
}

func evaluationFlags() []cli.Flag {
  f := []cli.Flag{
    cli.StringFlag{"file, f", "", "File to load"},
    cli.StringFlag{"log, l", "info", "Set log level"},
    cli.IntFlag{"timeout, t", 300, "Seconds to wait for state to load before timeout"},
  }
  return f
}
