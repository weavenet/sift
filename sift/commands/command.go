package command

import (
  "fmt"
  log "github.com/cihub/seelog"
  "github.com/codegangsta/cli"
  "strings"
)

func NewCommand(app *cli.App, cmdName string) (cmd cli.Command) {
  switch cmdName {
  case "repo":
    {
      cmd = repoCommand(app)
    }
  case "evaluation":
    {
      cmd = evaluationCommand(app)
    }
  case "plan":
    {
      cmd = planCommand(app)
    }
  }
  return cmd
}

func readArgs(flags []cli.Flag, context *cli.Context) (map[string]string, error) {
  args := map[string]string{}
  for _, arg := range flags {
    name := parseStringFlagName(arg)
    value := context.String(name)
    if value == "" {
      return args, fmt.Errorf("Required argument '%s' not provided.", name)
    }
    args[name] = value
  }
  return args, nil
}

func parseStringFlagName(flag cli.Flag) string {
  s := strings.Split(flag.String(), ",")[0]
  return strings.Replace(s, "--", "", 1)
}

func setLogLevel(level string) {
  logConfig := fmt.Sprintf("<seelog minlevel='%s'>", level)
  logger, _ := log.LoggerFromConfigAsBytes([]byte(logConfig))
  log.ReplaceLogger(logger)
}
