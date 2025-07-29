package main

import (
	"io"
	"os"
	"os/exec"

	"github.com/hunsy9/kubesnap/internal/constant"
	"github.com/hunsy9/kubesnap/internal/envutil"
	"github.com/hunsy9/kubesnap/internal/yamlutil"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
)

type Clusters struct {
	Contexts []ContextInfo `yaml:"contexts"`
}

type ContextInfo struct {
	Context Context `yaml:"context"`
	Name    string  `yaml:"name"`
}

type Context struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type SwitchContextCmd struct{}

func (_ SwitchContextCmd) Run(stdout, _ io.Writer) error {

	kubeConfigPath := os.Getenv("HOME") + constant.DefaultKubeConfigLocation
	filepath := envutil.GetEnvOrDefault("KUBECONFIG", kubeConfigPath)

	var unMarshalTarget Clusters

	yamlContext := yamlutil.NewParsingContext(filepath, &unMarshalTarget)
	err := yamlutil.ParseYaml(yamlContext)
	if err != nil {
		return errors.Wrap(err, "failed to parse yaml")
	}

	clusterNameList := make([]string, len(unMarshalTarget.Contexts))
	for i, ctx := range unMarshalTarget.Contexts {
		clusterNameList[i] = ctx.Name
	}

	prompt := promptui.Select{
		Label: "Select Context",
		Items: clusterNameList,
		Size:  50,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return errors.Wrap(err, "Prompt failed")
	}

	exec.Command("kubectl", "config", "use-context", result).Run()

	return nil
}
