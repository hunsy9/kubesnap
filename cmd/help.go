package main

import (
	"fmt"
	"io"
)

type HelpCmd struct{}

func (HelpCmd) Run(stdout, _ io.Writer) error {
	return showHelp(stdout)
}

func showHelp(stdout io.Writer) error {
	help := `USAGE:
  ks [SUBCOMMAND] | [OPTIONS]

OPTIONS:
  -h, --help    Show this message

COMMANDS:
  // General Command
  ks		Show current cluster/namespace detail
  
  // Context Command
  ks ctx        Open Context Switcher

  // Namespace Command
  ks ns         Open Namespace Switcher
  ks ns -       Switch to Default Namespace
`
	fmt.Fprint(stdout, help)
	return nil
}
