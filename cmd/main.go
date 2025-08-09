package main

import (
	"fmt"
	"io"
	"os"
)

type Cmd interface {
	Run(stdout, stderr io.Writer) error
}

func main() {
	cmd := parseCmd(os.Args[1:])
	if err := cmd.Run(os.Stdout, os.Stderr); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
