package main

import (
  "github.com/brettweavnet/sift/sift/commands"
  "github.com/codegangsta/cli"
  "os"
)

var commands = []string{"evaluation", "plan", "repo"}

func main() {
  app := cli.NewApp()
  app.Name = "sift"
  app.Usage = "A form of sift..."
  app.Version = "0.0.1"
  for _, c := range commands {
    app.Commands = append(app.Commands, command.NewCommand(app, c))
  }
  app.Run(os.Args)
}
