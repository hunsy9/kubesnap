package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	c "github.com/hunsy9/kubesnap/pkg/constant"
	"github.com/hunsy9/kubesnap/pkg/envutil"
	lv "github.com/hunsy9/kubesnap/pkg/ui/listview"
	"github.com/hunsy9/kubesnap/pkg/yamlutil"
	"github.com/pkg/errors"
)

type SwitchContextCmd struct{}

func (_ SwitchContextCmd) Run(stdout, _ io.Writer) error {

	// get kubeconfig path

	kubeConfigPath := os.Getenv("HOME") + c.DefaultKubeConfigLocation
	filepath := envutil.GetEnvOrDefault("KUBECONFIG", kubeConfigPath)

	// transform kubeconfig file to Kubeconfig model
	// TODO: use client-go functions instead of using yamlutil

	var unMarshalTarget yamlutil.KubeConfig

	yamlContext := yamlutil.NewParsingContext(filepath, &unMarshalTarget)
	err := yamlutil.ParseYaml(yamlContext)
	if err != nil {
		return errors.Wrap(err, "failed to parse yaml")
	}

	if len(unMarshalTarget.Contexts) == 0 {
		return errors.New("no contexts found")
	}

	// if there is contexts in kubeconfig file, push it to switchinglistmodel's Item list first

	items := []list.Item{}
	currentContextMarker := lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultActiveColor)).Bold(true).Render(c.CurrentContextMarker)
	items = append(items, lv.Item(unMarshalTarget.CurrentContext+currentContextMarker)) // push current context first

	for _, ctx := range unMarshalTarget.Contexts {
		if unMarshalTarget.CurrentContext == ctx.Name {
			continue
		}
		items = append(items, lv.Item(ctx.Name))
	}

	// create new bubbletea program with bunch of contexts

	md := lv.NewSwitchingUIModel(items, c.DefaultSwitchingContextHeaderMessage, c.Context)
	md.SetOperationFunc(SwitchContext)
	p := tea.NewProgram(md, tea.WithAltScreen())

	updatedModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// print updatedModel's output
	// it represents switched target

	if uiModel, ok := updatedModel.(*lv.UIModel); ok {
		if output := uiModel.GetOutput(); output != "" {
			fmt.Print(output)
		}
	}

	return nil
}

// TODO: modify kubeconfig file's current-context area instead of using kubectl
func SwitchContext(contextName string) tea.Cmd {
	return func() tea.Msg {
		exec_err := exec.Command("kubectl", "config", "use-context", contextName).Run()
		return lv.OperationResultMsg{Err: exec_err}
	}
}
