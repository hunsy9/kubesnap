package kubestatus

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	c "github.com/hunsy9/kubesnap/pkg/constant"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ApiHealthMsg struct {
	Msg string
}

type AuthInfoMsg struct {
	Msg *authenticationv1.SelfSubjectReview
	Err error
}

type NodeInfoMsg struct {
	TotalNodes       int
	ReadyNodes       int
	NodeAPIStatusErr string
}

type PodInfoMsg struct {
	TotalPods       int
	RunningPods     int
	PodAPIStatusErr string
}

type EventInfoMsg struct {
	WarningEvents     int
	EventAPIStatusErr string
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

type StatusModel struct {
	spinner         spinner.Model
	clientSet       *kubernetes.Clientset
	clusterName     string
	namespace       string
	endpoint        string
	apiStatus       string
	userInfo        string
	groupsInfo      string
	nodeInfo        string
	podInfo         string
	eventInfo       string
	isHealthLoading bool
	isAuthLoading   bool
	isNodeLoading   bool
	isPodLoading    bool
	isEventLoading  bool
}

func NewStatusModel(clientSet *kubernetes.Clientset, clusterName string, namespace string, endpoint string) *StatusModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = spinnerStyle
	return &StatusModel{
		spinner:         sp,
		clientSet:       clientSet,
		clusterName:     clusterName,
		namespace:       namespace,
		endpoint:        endpoint,
		isHealthLoading: true,
		isAuthLoading:   true,
		isNodeLoading:   true,
		isPodLoading:    true,
		isEventLoading:  true,
	}
}

func (m *StatusModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.checkApiHealth,
		m.getAuthInfo,
		m.getNodeInfo,
		m.getPodInfo,
		m.getEventInfo,
	)
}

func (m *StatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var spinCmd tea.Cmd
		m.spinner, spinCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinCmd)

	case ApiHealthMsg:
		m.apiStatus = msg.Msg
		m.isHealthLoading = false

	case AuthInfoMsg:
		if msg.Err == nil {
			m.userInfo = msg.Msg.Status.UserInfo.Username
			m.groupsInfo = fmt.Sprintf("%s", msg.Msg.Status.UserInfo.Groups)
		} else {
			m.userInfo = iconStyle.Bold(true).Render("-")
			m.groupsInfo = iconStyle.Bold(true).Render("-")
		}
		m.isAuthLoading = false

	case NodeInfoMsg:
		if msg.NodeAPIStatusErr == "" {
			m.nodeInfo = fmt.Sprintf("%v Nodes ( %v Running )", msg.TotalNodes, msg.ReadyNodes)
		} else {
			m.nodeInfo = msg.NodeAPIStatusErr
		}
		m.isNodeLoading = false

	case PodInfoMsg:
		if msg.PodAPIStatusErr == "" {
			m.podInfo = fmt.Sprintf("%v Pods ( %v Running )", msg.TotalPods, msg.RunningPods)
		} else {
			m.podInfo = msg.PodAPIStatusErr
		}
		m.isPodLoading = false

	case EventInfoMsg:
		if msg.EventAPIStatusErr == "" {
			if msg.WarningEvents > 0 {
				m.eventInfo = warningStyle.Render(fmt.Sprintf("%v Warning ⚠︎", msg.WarningEvents))
			} else {
				m.eventInfo = activeStyle.Render("Healthy")
			}
		} else {
			m.eventInfo = msg.EventAPIStatusErr
		}
		m.isEventLoading = false
	}

	return m, tea.Batch(cmds...)
}

// helper func: render spinner if is_Loading variable is true
func (m *StatusModel) renderLoadingOrContent(loading bool, content string) string {
	if loading {
		return m.spinner.View()
	}
	return content
}

func (m *StatusModel) View() string {

	apiStatusView := m.renderLoadingOrContent(m.isHealthLoading, m.apiStatus)
	userView := m.renderLoadingOrContent(m.isAuthLoading, m.userInfo)
	groupView := m.renderLoadingOrContent(m.isAuthLoading, m.groupsInfo)
	nodeView := m.renderLoadingOrContent(m.isNodeLoading, m.nodeInfo)
	podView := m.renderLoadingOrContent(m.isPodLoading, m.podInfo)
	eventView := m.renderLoadingOrContent(m.isEventLoading, m.eventInfo)

	clusterConnection := lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render("Cluster Connection\n"),
		fmt.Sprintf(" %s Cluster: %s\n %s API Server: %s\n    └─ Status: %s\n",
			iconStyle.Render(c.Bullet),
			m.clusterName, iconStyle.Render(c.Bullet),
			urlStyle.Render(m.endpoint),
			apiStatusView),
	)

	clusterAuth := lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render("Auth & Scope\n"),
		fmt.Sprintf(" %s User: %s\n %s Groups: %s\n %s Namespace: %s\n",
			iconStyle.Render(c.Bullet),
			userView,
			iconStyle.Render(c.Bullet),
			groupView,
			iconStyle.Render(c.Bullet),
			m.namespace),
	)

	clusterOverview := lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render("Resource Overview\n"),
		fmt.Sprintf(" %s Nodes: %s\n %s Pods: %s\n %s Events: %s\n",
			iconStyle.Render(c.Bullet),
			nodeView,
			iconStyle.Render(c.Bullet),
			podView,
			iconStyle.Render(c.Bullet),
			eventView),
	)

	styledOverview := lipgloss.NewStyle().Width(35).Render(clusterOverview)
	styledAuth := lipgloss.NewStyle().MarginLeft(3).Render(clusterAuth)

	bottomLayout := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styledOverview,
		styledAuth,
	)

	layout := boxStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			clusterConnection,
			bottomLayout,
		),
	)
	return layout
}
