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
  kubesnap [SUBCOMMAND] | [OPTIONS]

OPTIONS:
  -h, --help    Show this message
  -v, --version Show current version of kubesnap

COMMANDS:
  kubesnap		Show current cluster/namespace overview
  
  kubesnap ctx        Open context switcher
  └─'d' key     Switch to context delete mode   
  └─'r' key     Switch to context rename mode

  kubesnap ns         Open namespace switcher
  kubesnap ns ~       Switch to default namespace
`
	fmt.Fprint(stdout, help)
	return nil
}
