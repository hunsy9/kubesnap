package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hunsy9/kubesnap/pkg/constant"
	"github.com/hunsy9/kubesnap/pkg/envutil"
	"github.com/hunsy9/kubesnap/pkg/model"
	switchcontextlistmodel "github.com/hunsy9/kubesnap/pkg/model/switchcontextlist"
	"github.com/hunsy9/kubesnap/pkg/yamlutil"
	"github.com/pkg/errors"
)

type SwitchContextCmd struct{}

func (_ SwitchContextCmd) Run(stdout, _ io.Writer) error {

	// get kubeconfig path

	kubeConfigPath := os.Getenv("HOME") + constant.DefaultKubeConfigLocation
	filepath := envutil.GetEnvOrDefault("KUBECONFIG", kubeConfigPath)

	// transform kubeconfig file to Kubeconfig model

	var unMarshalTarget model.KubeConfig

	yamlContext := yamlutil.NewParsingContext(filepath, &unMarshalTarget)
	err := yamlutil.ParseYaml(yamlContext)
	if err != nil {
		return errors.Wrap(err, "failed to parse yaml")
	}

	if len(unMarshalTarget.Contexts) == 0 {
		return errors.New("no contexts found")
	}

	// if there is contexts in kubeconfig file, push it to switchcontextlistmodel's Item list

	items := []list.Item{}
	currentContextMarker := lipgloss.NewStyle().Foreground(lipgloss.Color(constant.CurrentContextMarkerColor)).Bold(true).Render(constant.CurrentContextMarker)
	items = append(items, switchcontextlistmodel.Item(unMarshalTarget.CurrentContext+currentContextMarker)) // push current context first

	for _, ctx := range unMarshalTarget.Contexts {
		if unMarshalTarget.CurrentContext == ctx.Name {
			continue
		}
		items = append(items, switchcontextlistmodel.Item(ctx.Name))
	}

	// create new bubbletea program with bunch of contexts

	md := switchcontextlistmodel.NewUIModel(items, constant.DefaultSwitchingContextHeaderMessage)
	p := tea.NewProgram(md, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	//

	if uiModel, ok := finalModel.(*switchcontextlistmodel.UIModel); ok {
		if output := uiModel.GetOutput(); output != "" {
			fmt.Print(output)
		}
	}

	return nil
}
