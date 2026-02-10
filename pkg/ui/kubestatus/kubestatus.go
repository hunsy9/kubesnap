package kubestatus

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	c "github.com/hunsy9/kubesnap/pkg/constant"
	"github.com/hunsy9/kubesnap/pkg/version"
	"k8s.io/client-go/kubernetes"
)

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
	totalAsyncTasks int
	completedTasks  int
	latestVersion   string
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
		totalAsyncTasks: 6,
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
		m.checkForUpdate,
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
		cmds = append(cmds, m.incrementProgress())

	case AuthInfoMsg:
		if msg.Err == nil {
			m.userInfo = msg.Msg.Status.UserInfo.Username
			m.groupsInfo = fmt.Sprintf("%s", msg.Msg.Status.UserInfo.Groups)
		} else {
			m.userInfo = iconStyle.Bold(true).Render("-")
			m.groupsInfo = iconStyle.Bold(true).Render("-")
		}
		m.isAuthLoading = false
		cmds = append(cmds, m.incrementProgress())

	case NodeInfoMsg:
		if msg.NodeAPIStatusErr == "" {
			m.nodeInfo = fmt.Sprintf("%v Nodes ( %v Running )", msg.TotalNodes, msg.ReadyNodes)
		} else {
			m.nodeInfo = msg.NodeAPIStatusErr
		}
		m.isNodeLoading = false
		cmds = append(cmds, m.incrementProgress())

	case PodInfoMsg:
		if msg.PodAPIStatusErr == "" {
			m.podInfo = fmt.Sprintf("%v Pods ( %v Running )", msg.TotalPods, msg.RunningPods)
		} else {
			m.podInfo = msg.PodAPIStatusErr
		}
		m.isPodLoading = false
		cmds = append(cmds, m.incrementProgress())

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
		cmds = append(cmds, m.incrementProgress())

	case UpdateAvailableMsg:
		m.latestVersion = string(msg)
		cmds = append(cmds, m.incrementProgress())

	case AsyncTaskAllDoneMsg:
		return m, tea.Quit

	}

	return m, tea.Batch(cmds...)
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

	if m.latestVersion != "" {
		updateMsg := updageMsgStyle.Render(fmt.Sprintf("⚡ New update available! version %s -> %s\n👉 Check update guide: %s", version.Version, m.latestVersion, updateUrlStyle.Render(c.MaintenanceGuideURL)))
		layout = lipgloss.JoinVertical(
			lipgloss.Left,
			updateMsg,
			layout,
		)
	}

	return layout
}
