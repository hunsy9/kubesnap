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
	fmt.Printf("help command")
	return nil
}
