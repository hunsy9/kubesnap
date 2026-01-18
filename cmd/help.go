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
  ks		Show current cluster/namespace detail
  ks ctx        Interactively list and select kubernetes context
  ks ns         Interactively list and select kubernetes namespace
`
	fmt.Fprint(stdout, help)
	return nil
}
