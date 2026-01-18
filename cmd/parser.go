package main

import (
	"io"

	"github.com/pkg/errors"
)

type ErrorCmd struct{ Err error }

func (cmd ErrorCmd) Run(_, _ io.Writer) error {
	return cmd.Err
}

func parseCmd(argv []string) Cmd {

	// feature: switching context and routing command
	if len(argv) == 0 {
		return InfoCmd{}
	}

	if argv[0] == "ctx" {
		return SwitchContextCmd{}
	}

	if argv[0] == "ns" {
		return SwitchNamespaceCmd{}
	}

	if argv[0] == "-h" || argv[0] == "--help" {
		return HelpCmd{}
	}

	return ErrorCmd{Err: errors.New("kubesnap: unknown command")}
}
