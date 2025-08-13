package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hunsy9/kubesnap/pkg/constant"
	"github.com/hunsy9/kubesnap/pkg/envutil"
	"github.com/hunsy9/kubesnap/pkg/model"
	"github.com/hunsy9/kubesnap/pkg/yamlutil"
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

	if len(unMarshalTarget.Contexts) == 0 {
		return errors.New("no contexts found")
	}

	items := []list.Item{}

	for _, ctx := range unMarshalTarget.Contexts {
		items = append(items, model.Item(ctx.Name))
	}

	md := model.NewUIModel(items, "Select a Context")
	p := tea.NewProgram(md, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return nil
}
