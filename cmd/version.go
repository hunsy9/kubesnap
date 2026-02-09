package main

import (
	"fmt"
	"io"

	"github.com/hunsy9/kubesnap/pkg/version"
)

type VersionCmd struct{}

func (VersionCmd) Run(stdout, _ io.Writer) error {
	return showVersion(stdout)
}

func showVersion(stdout io.Writer) error {
	versionMsg := "Current Version: v" + version.Version + "\n"
	fmt.Fprint(stdout, versionMsg)
	return nil
}
