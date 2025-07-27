package main

import (
	"fmt"
	"io"
)

type ErrorCmd struct{ Err error }

func (cmd ErrorCmd) Run(_, _ io.Writer) error {
	return cmd.Err
}

func parseCmd(argv []string) Cmd {
	if len(argv) == 0 {
		return HelpCmd{}
	}
	return ErrorCmd{Err: fmt.Errorf("unknown command")}
}
