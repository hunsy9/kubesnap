package listview

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/client-go/tools/clientcmd"
)

type SwitchingOperation func(param string) tea.Cmd // switching function executed by UIModel

type DeletingOperation func(targets []string) tea.Cmd

func SwitchContextOperation(contextName string) tea.Cmd {
	return func() tea.Msg {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		config, err := loadingRules.Load()

		if err != nil {
			return SwitchResultMsg{Err: err}
		}

		if _, exists := config.Contexts[contextName]; !exists {
			return SwitchResultMsg{Err: fmt.Errorf("context '%s' not found", contextName)}
		}

		config.CurrentContext = contextName

		err = clientcmd.ModifyConfig(loadingRules, *config, true)
		if err != nil {
			return SwitchResultMsg{Err: err}
		}

		time.Sleep(time.Millisecond * 400)
		return SwitchResultMsg{Err: nil}
	}
}

func SwitchNamespaceOperation(namespaceName string) tea.Cmd {
	return func() tea.Msg {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		config, err := loadingRules.Load()

		if err != nil {
			return SwitchResultMsg{Err: err}
		}

		ctxName := config.CurrentContext
		if ctxName == "" {
			return SwitchResultMsg{Err: fmt.Errorf("no current context set")}
		}

		ctx, exists := config.Contexts[ctxName]
		if !exists {
			return SwitchResultMsg{Err: fmt.Errorf("context '%s' not found", ctxName)}
		}

		ctx.Namespace = namespaceName

		err = clientcmd.ModifyConfig(loadingRules, *config, true)
		if err != nil {
			return SwitchResultMsg{Err: err}
		}

		time.Sleep(time.Millisecond * 300)
		return SwitchResultMsg{Err: nil}
	}
}

func DeleteOperation(targets []string) tea.Cmd {
	return func() tea.Msg {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		config, err := loadingRules.Load()

		if err != nil {
			return DeletionResultMsg{Err: err, Targets: nil}
		}

		for _, targetCtx := range targets {
			delete(config.Contexts, targetCtx)
		}

		err = clientcmd.ModifyConfig(loadingRules, *config, true)
		if err != nil {
			return DeletionResultMsg{Err: err, Targets: nil}
		}

		time.Sleep(time.Millisecond * 400)

		return DeletionResultMsg{Err: nil, Targets: targets}
	}
}

func (m *SwitchingUIModel) RenameOperation(newName string) tea.Cmd {
	return func() tea.Msg {
		if m.choice == newName {
			return RenameResultMsg{Err: nil}
		}

		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		config, err := loadingRules.Load()
		if err != nil {
			return RenameResultMsg{Err: err}
		}

		if _, exists := config.Contexts[newName]; exists {
			return RenameResultMsg{Err: fmt.Errorf("context named '%s' already exists", newName)}
		}

		ctx, exists := config.Contexts[m.choice]
		if !exists {
			return RenameResultMsg{Err: fmt.Errorf("context '%s' not found", m.choice)}
		}

		config.Contexts[newName] = ctx
		delete(config.Contexts, m.choice)

		if config.CurrentContext == m.choice {
			config.CurrentContext = newName
		}

		if err := clientcmd.ModifyConfig(loadingRules, *config, true); err != nil {
			return RenameResultMsg{Err: err}
		}

		return RenameResultMsg{OriginalName: m.choice, NewName: newName, Err: nil}
	}
}
