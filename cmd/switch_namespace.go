package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"

	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hunsy9/kubesnap/pkg/constant"
	"github.com/hunsy9/kubesnap/pkg/envutil"
	slm "github.com/hunsy9/kubesnap/pkg/model/switchingList"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SwitchNamespaceCmd struct{}

func (_ SwitchNamespaceCmd) Run(stdout, _ io.Writer) error {

	// get kubeconfig path

	kubeConfigPath := os.Getenv("HOME") + constant.DefaultKubeConfigLocation
	filepath := envutil.GetEnvOrDefault("KUBECONFIG", kubeConfigPath)

	// get the clientset of current context

	clientset, err := getClientSet(filepath)
	if err != nil {
		return errors.Wrap(err, "failed to getting kubernetes clientset")
	}

	// get namespace list via clientset

	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to listing namespaces from the cluster")
	}

	// if there is current-context's namespace in kubeconfig file, push it to switchinglistmodel's Item list first

	items := []list.Item{}
	currentNamespace := getCurrentNamespace()
	currentNamespaceMarker := lipgloss.NewStyle().Foreground(lipgloss.Color(constant.CurrentNamespaceMarkerColor)).Bold(true).Render(constant.CurrentNamespaceMarker)

	items = append(items, slm.Item(currentNamespace+currentNamespaceMarker)) // push current namespace first

	for _, namespace := range namespaces.Items {
		if namespace.Name == currentNamespace {
			continue
		}
		items = append(items, slm.Item(namespace.Name))
	}

	// create new bubbletea program with bunch of namespaces

	md := slm.NewUIModel(items, constant.DefaultSwitchingNamespaceHeaderMessage, constant.Namespace)
	md.SetOperationFunc(SwitchNamespace)
	p := tea.NewProgram(md, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if uiModel, ok := finalModel.(*slm.UIModel); ok {
		if output := uiModel.GetOutput(); output != "" {
			fmt.Print(output)
		}
	}

	return nil
}

func getClientSet(kubeconfigPath string) (*kubernetes.Clientset, error) {
	currentConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(currentConfig)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func getCurrentNamespace() string {
	config, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err.Error())
	}

	currentContext := config.CurrentContext
	if context, exists := config.Contexts[currentContext]; exists {
		if context.Namespace != "" {
			return context.Namespace
		}
	}

	return "default"
}

// TODO: modify kubeconfig file's current-context's namespace area instead of using kubectl
func SwitchNamespace(namespaceName string) tea.Cmd {
	return func() tea.Msg {
		exec_err := exec.Command("kubectl", "config", "set-context", "--current", "--namespace="+namespaceName).Run()
		return slm.OperationResultMsg{Err: exec_err}
	}
}
