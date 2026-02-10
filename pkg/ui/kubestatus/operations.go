package kubestatus

import (
	"context"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hunsy9/kubesnap/pkg/version"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *StatusModel) checkForUpdate() tea.Msg {
	latest, err := version.CheckVersionUpdate()
	if err != nil || latest == "" {
		return UpdateAvailableMsg("")
	}
	return UpdateAvailableMsg(latest)
}

func (m *StatusModel) getNodeInfo() tea.Msg {
	var msg NodeInfoMsg

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Lookup Nodes
	nodes, err := m.clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		if statusErr, ok := err.(interface{ Status() metav1.Status }); ok {
			msg.NodeAPIStatusErr = http.StatusText(int(statusErr.Status().Code))
		} else {
			msg.NodeAPIStatusErr = "Unknown"
		}
	}

	msg.TotalNodes = len(nodes.Items)
	for _, node := range nodes.Items {
		for _, cond := range node.Status.Conditions {
			if cond.Type == "Ready" && cond.Status == "True" {
				msg.ReadyNodes++
				break
			}
		}
	}

	return msg
}

func (m *StatusModel) getPodInfo() tea.Msg {
	var msg PodInfoMsg

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Lookup pods in all namespaces
	pods, err := m.clientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})

	if err == nil {
		msg.TotalPods = len(pods.Items)
		for _, pod := range pods.Items {
			if pod.Status.Phase == "Running" {
				msg.RunningPods++
			}
		}
	} else { // fallback: if all namespace pod list is unauthorized, find namespaces it can be listed
		nsList, nsErr := m.clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		var targetNamespaces []string
		if nsErr == nil {
			for _, ns := range nsList.Items {
				targetNamespaces = append(targetNamespaces, ns.Name)
			}
		} else {
			targetNamespaces = []string{m.namespace}
		}

		successCount := 0
		for _, ns := range targetNamespaces {
			nsPods, nsPodErr := m.clientSet.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
			if nsPodErr == nil {
				successCount++
				msg.TotalPods += len(nsPods.Items)
				for _, pod := range nsPods.Items {
					if pod.Status.Phase == "Running" {
						msg.RunningPods++
					}
				}
			}
		}

		if successCount == 0 {
			if statusErr, ok := err.(interface{ Status() metav1.Status }); ok {
				msg.PodAPIStatusErr = http.StatusText(int(statusErr.Status().Code))
			} else {
				msg.PodAPIStatusErr = "Unknown"
			}
		}
	}

	return msg
}

func (m *StatusModel) getEventInfo() tea.Msg {
	var msg EventInfoMsg

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Lookup events in all namespaces
	events, err := m.clientSet.CoreV1().Events("").List(ctx, metav1.ListOptions{})

	if err == nil {
		for _, ev := range events.Items {
			if ev.Type == "Warning" {
				msg.WarningEvents++
			}
		}
	} else { // fallback: if all namespace event list is unauthorized, find events it can be listed
		nsList, nsErr := m.clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		var targetNamespaces []string
		if nsErr == nil {
			for _, ns := range nsList.Items {
				targetNamespaces = append(targetNamespaces, ns.Name)
			}
		} else {
			targetNamespaces = []string{m.namespace}
		}

		successCount := 0
		for _, ns := range targetNamespaces {
			nsEvents, nsEventErr := m.clientSet.CoreV1().Events(ns).List(ctx, metav1.ListOptions{})
			if nsEventErr == nil {
				successCount++
				for _, e := range nsEvents.Items {
					if e.Type == "Warning" {
						msg.WarningEvents++
					}
				}
			}
		}

		if successCount == 0 {
			if statusErr, ok := err.(interface{ Status() metav1.Status }); ok {
				msg.EventAPIStatusErr = http.StatusText(int(statusErr.Status().Code))
			} else {
				msg.EventAPIStatusErr = "Unknown"
			}
		}
	}

	return msg
}

func (m *StatusModel) checkApiHealth() tea.Msg {

	var statusCode int
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.clientSet.Discovery().RESTClient().
		Get().
		AbsPath("/healthz").
		Do(ctx).
		StatusCode(&statusCode).
		Error()
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		var statusText string
		if ctx.Err() == context.DeadlineExceeded {
			statusText = "Network Timeout"
		} else if statusCode == 0 {
			statusText = "Network Unreachable"
		} else {
			statusText = http.StatusText(statusCode)
		}
		return ApiHealthMsg{
			Msg: errorStyle.Render(fmt.Sprintf("● InActive (%v)", statusText)),
		}
	}

	activeMsg := activeStyle.Render(fmt.Sprintf("● Active (%vms)", elapsed))
	return ApiHealthMsg{
		Msg: activeMsg,
	}
}

func (m *StatusModel) getAuthInfo() tea.Msg {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ssr := &authenticationv1.SelfSubjectReview{}
	result, err := m.clientSet.AuthenticationV1().SelfSubjectReviews().Create(ctx, ssr, metav1.CreateOptions{})

	if err != nil {
		return AuthInfoMsg{
			Err: err,
		}
	}
	return AuthInfoMsg{
		Msg: result,
	}
}
