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
  ks [OPTIONS] | [SUBCOMMAND]

OPTIONS:
  -h, --help    Show this message

COMMANDS:
  ks        Interactively list and select kubernetes context
`
	fmt.Fprint(stdout, help)
	return nil
}
