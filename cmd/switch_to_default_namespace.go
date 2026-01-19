package main

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/hunsy9/kubesnap/pkg/constant"
	"github.com/pkg/errors"
)

type SwitchToDefaultNamespaceCmd struct{}

func (_ SwitchToDefaultNamespaceCmd) Run(stdout, _ io.Writer) error {

	currentNamespace := GetCurrentNamespace()
	if currentNamespace == constant.DefaultNamespace {
		return errors.New("You are already in the default namespace")
	}

	exec_err := exec.Command("kubectl", "config", "set-context", "--current", "--namespace="+constant.DefaultNamespace).Run()
	if exec_err != nil {
		return errors.Wrap(exec_err, "Failed to switch namespace to default")
	}

	fmt.Println("Namespace switched to default")
	return nil
}
