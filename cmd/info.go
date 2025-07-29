package main

import (
	"io"

	"github.com/fatih/color"
)

type InfoCmd struct{}

func (InfoCmd) Run(stdout, stderr io.Writer) error {
	info := `kubesnap info`

	color.New(color.FgBlue).Fprintf(stdout, "%s\n", info)
	return nil
}
