package main

import (
	"fmt"
	"io"
)

type InfoCmd struct{}

func (InfoCmd) Run(stdout, stderr io.Writer) error {
	info := "kubesnap info"
	fmt.Fprintln(stdout, info)
	return nil
}
