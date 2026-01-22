package kubestatus

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hunsy9/kubesnap/pkg/constant"
	authenticationv1 "k8s.io/api/authentication/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(constant.DefaultThemeColor))
)

type StatusModel struct {
	spinner         spinner.Model // spinner model
	clientSet       *kubernetes.Clientset
	clusterName     string
	namespace       string
	endpoint        string
	auth            *authenticationv1.SelfSubjectReview
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
	return nil
}
