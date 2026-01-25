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
	spinner         spinner.Model // spinner model
	clientSet       *kubernetes.Clientset
	clusterName     string
	namespace       string
	endpoint        string
	auth            *authenticationv1.SelfSubjectReview
	apiStatus       string
	isHealthLoading bool
	isAuthLoading   bool
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
	}
}

func (m *StatusModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.checkApiHealth,
		m.getAuthInfo,
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
		m.auth = msg.Msg
		m.isAuthLoading = false
	}

	return m, tea.Batch(cmds...)
}

func (m *StatusModel) View() string {

	var apiStatusView string

	if m.isHealthLoading {
		apiStatusView = m.spinner.View()
	} else {
		apiStatusView = m.apiStatus
	}

	var userView, groupView string

	if m.isAuthLoading {
		userView = m.spinner.View()
		groupView = m.spinner.View()
	} else if m.auth != nil {
		userView = fmt.Sprintf("%v", m.auth.Status.UserInfo.Username)
		groupView = fmt.Sprintf("%v", m.auth.Status.UserInfo.Groups)
	} else {
		userView = iconStyle.Bold(true).Render("-")
		groupView = iconStyle.Bold(true).Render("-")
	}

	clusterConnection := fmt.Sprintf(" %s Cluster: %s\n %s API Server: %s\n    └─ Status: %s\n",
		iconStyle.Render(c.Bullet),
		m.clusterName, iconStyle.Render(c.Bullet),
		urlStyle.Render(m.endpoint),
		apiStatusView)

	clusterAuth := fmt.Sprintf(" %s User: %s\n %s Groups: %s\n %s Namespace: %s\n",
		iconStyle.Render(c.Bullet),
		userView,
		iconStyle.Render(c.Bullet),
		groupView,
		iconStyle.Render(c.Bullet),
		m.namespace)

	layout := boxStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			headerStyle.Render("CLUSTER CONNECTION\n"),
			clusterConnection,
			headerStyle.Render("AUTH & SCOPE\n"),
			clusterAuth),
	)
	return layout
}
