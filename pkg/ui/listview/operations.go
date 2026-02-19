package listview

import (
	"fmt"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/client-go/tools/clientcmd"
)

type SwitchingOperation func(param string) tea.Cmd // switching function executed by UIModel

type DeletingOperation func(targets []string) tea.Cmd

// TODO: modify kubeconfig file's current-context area instead of using kubectl
func SwitchContextOperation(contextName string) tea.Cmd {
	return func() tea.Msg {
		exec_err := exec.Command("kubectl", "config", "use-context", contextName).Run()
		time.Sleep(time.Millisecond * 500)
		return SwitchResultMsg{Err: exec_err}
	}
}

// TODO: modify kubeconfig file's current-context's namespace area instead of using kubectl
func SwitchNamespaceOperation(namespaceName string) tea.Cmd {
	return func() tea.Msg {
		exec_err := exec.Command("kubectl", "config", "set-context", "--current", "--namespace="+namespaceName).Run()
		return SwitchResultMsg{Err: exec_err}
	}
}

func DeleteOperation(targets []string) tea.Cmd {
	return func() tea.Msg {
		config, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)

		if err != nil {
			return DeletionResultMsg{Err: err, Targets: nil}
		}

		for _, targetCtx := range targets {
			if _, exists := config.Contexts[targetCtx]; exists {
				delete(config.Contexts, targetCtx)
			}
		}

		err = clientcmd.WriteToFile(*config, clientcmd.RecommendedHomeFile)
		if err != nil {
			return DeletionResultMsg{Err: err, Targets: nil}
		}

		time.Sleep(time.Millisecond * 500)

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
