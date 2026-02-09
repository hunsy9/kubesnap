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
  -v, --version Show current version of kubesnap

COMMANDS:
  ks		Show current cluster/namespace overview
  
  ks ctx        Open context switcher
  └─'d' key     Switch to context delete mode   
  └─'r' key     Switch to context rename mode

  ks ns         Open namespace switcher
  ks ns ~       Switch to default namespace
`
	fmt.Fprint(stdout, help)
	return nil
}
