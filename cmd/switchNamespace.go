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
	c "github.com/hunsy9/kubesnap/pkg/constant"
	"github.com/hunsy9/kubesnap/pkg/envutil"
	lv "github.com/hunsy9/kubesnap/pkg/ui/listview"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SwitchNamespaceCmd struct{}

type SwitchToDefaultNamespaceCmd struct{}

func (_ SwitchNamespaceCmd) Run(stdout, _ io.Writer) error {

	// get kubeconfig path

	kubeConfigPath := os.Getenv("HOME") + c.DefaultKubeConfigLocation
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
	currentNamespace := GetCurrentNamespace()
	currentNamespaceMarker := lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultActiveColor)).Bold(true).Render(c.CurrentNamespaceMarker)
	currentNamespaceItem := lv.Item{
		DisplayName: currentNamespace + currentNamespaceMarker,
		Name:        currentNamespace,
	}

	items = append(items, currentNamespaceItem) // push current namespace first

	for _, namespace := range namespaces.Items {
		if namespace.Name == currentNamespace {
			continue
		}
		namespaceItem := lv.Item{
			DisplayName: namespace.Name,
			Name:        namespace.Name,
		}
		items = append(items, namespaceItem)
	}

	// create new bubbletea program with bunch of namespaces

	md := lv.NewSwitchingUIModel(items, c.DefaultSwitchingNamespaceHeaderMessage, c.Namespace)
	md.SetOperationFunc(lv.SwitchNamespaceOperation)
	p := tea.NewProgram(md, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if switchingUIModel, ok := finalModel.(*lv.SwitchingUIModel); ok {
		if output := switchingUIModel.GetOutput(); output != "" {
			fmt.Print(output)
		}
	}

	return nil
}

func (_ SwitchToDefaultNamespaceCmd) Run(stdout, _ io.Writer) error {

	currentNamespace := GetCurrentNamespace()
	if currentNamespace == c.DefaultNamespace {
		return errors.New("You are already in the default namespace")
	}

	// TODO: modify kubeconfig file's current-context's namespace area instead of using kubectl
	exec_err := exec.Command("kubectl", "config", "set-context", "--current", "--namespace="+c.DefaultNamespace).Run()
	if exec_err != nil {
		return errors.Wrap(exec_err, "Failed to switch namespace to default")
	}

	fmt.Println("Namespace switched to default")
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

func GetCurrentNamespace() string {
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

	return c.DefaultNamespace
}
