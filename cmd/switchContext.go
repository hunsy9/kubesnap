package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	c "github.com/hunsy9/kubesnap/pkg/constant"
	lv "github.com/hunsy9/kubesnap/pkg/ui/listview"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
)

type SwitchContextCmd struct{}

func (_ SwitchContextCmd) Run(stdout, _ io.Writer) error {

	// load kubeconfig using client-go
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	config, err := loadingRules.Load()
	if err != nil {
		return errors.Wrap(err, "failed to load kubeconfig")
	}

	if len(config.Contexts) == 0 {
		return errors.New("no contexts found")
	}

	// build list items from contexts
	items := []list.Item{}
	currentContextMarker := lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultActiveColor)).Bold(true).Render(c.CurrentContextMarker)

	// add current context first
	if config.CurrentContext != "" {
		currentContextItem := lv.Item{
			DisplayName: config.CurrentContext + currentContextMarker,
			Name:        config.CurrentContext,
		}
		items = append(items, currentContextItem)
	}

	// add other contexts
	for ctxName := range config.Contexts {
		if ctxName == config.CurrentContext {
			continue
		}
		contextItem := lv.Item{
			DisplayName: ctxName,
			Name:        ctxName,
		}
		items = append(items, contextItem)
	}

	// create new bubbletea program with bunch of contexts

	md := lv.NewSwitchingUIModel(items, c.DefaultSwitchingContextHeaderMessage, c.Context)
	md.SetOperationFunc(lv.SwitchContextOperation)
	p := tea.NewProgram(md, tea.WithAltScreen())

	updatedModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// print updatedModel's output
	// it represents switched target

	if switchingUIModel, ok := updatedModel.(*lv.SwitchingUIModel); ok {
		if output := switchingUIModel.GetOutput(); output != "" {
			fmt.Print(output)
		}
	}

	return nil
}
