package main

import (
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	c "github.com/hunsy9/kubesnap/pkg/constant"
	"github.com/hunsy9/kubesnap/pkg/envutil"
	ks "github.com/hunsy9/kubesnap/pkg/model/kubestatus"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
)

type InfoCmd struct{}

func (InfoCmd) Run(stdout, stderr io.Writer) error {

	// get client set via kubeconfig file

	kubeConfigPath := os.Getenv("HOME") + c.DefaultKubeConfigLocation
	filepath := envutil.GetEnvOrDefault("KUBECONFIG", kubeConfigPath)

	clientSet, err := getClientSet(filepath)
	if err != nil {
		return errors.WithMessage(err, "unexpected error occured")
	}

	// get current cluster's name

	config, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)
	if err != nil {
		return errors.New("failed to load config from kubeconfig")
	}

	clusterName := config.CurrentContext
	namespace := GetCurrentNamespace()

	// get current cluster's api server endpoint

	endpoint, err := getCurrentApiServerEndpoint()
	if err != nil {
		return errors.Cause(err)
	}

	md := ks.NewStatusModel(clientSet, clusterName, namespace, endpoint)
	p := tea.NewProgram(md)

	_, err = p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return nil
}

func getCurrentApiServerEndpoint() (string, error) {
	config, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)
	if err != nil {
		return "", errors.New("failed to load config from kubeconfig")
	}

	currentContext := config.CurrentContext
	ctx, exists := config.Contexts[currentContext]

	if !exists {
		return "", errors.Errorf("context '%s' not found", currentContext)
	}

	cluster, exists := config.Clusters[ctx.Cluster]

	if !exists {
		return "", errors.Errorf("cluster '%s' not found", ctx.Cluster)
	}

	if cluster.Server == "" {
		return "", errors.Errorf("api server endpoint is empty for cluster '%s'", ctx.Cluster)
	}

	return cluster.Server, nil
}
