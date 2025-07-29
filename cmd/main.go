package main

import (
	"io"
	"os"

	"github.com/fatih/color"
)

type Cmd interface {
	Run(stdout, stderr io.Writer) error
}

func main() {
	cmd := parseCmd(os.Args[1:])
	if err := cmd.Run(color.Output, color.Error); err != nil {
		color.Red("Error: %v", err)
	}
}
