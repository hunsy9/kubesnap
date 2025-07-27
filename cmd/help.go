package main

import (
	"fmt"
	"io"
)

type HelpCmd struct{}

func (_ HelpCmd) Run(stdout, _ io.Writer) error {
	return showHelp(stdout)
}

func showHelp(out io.Writer) error {
	help := `USAGE:
  ks [OPTIONS] | [SUBCOMMAND]

OPTIONS:
  -h, --help    Show this message

COMMANDS:
  ks ctx        Interactively list and select kubernetes context
`
	fmt.Fprint(out, help)
	return nil
}
