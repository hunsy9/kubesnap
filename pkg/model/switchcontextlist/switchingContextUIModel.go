package model

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hunsy9/kubesnap/pkg/constant"
)

type kubectlResult struct {
	err error
}

// TODO: modify kubeconfig file's current-context area instead of using kubectl
func (m *UIModel) executeKubectl() tea.Cmd {
	return func() tea.Msg {
		exec_err := exec.Command("kubectl", "config", "use-context", m.choice).Run()
		return kubectlResult{err: exec_err}
	}
}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color(constant.DefaultThemeColor))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color(constant.DefaultThemeColor))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(3)
)

type UIModel struct {
	spinner   spinner.Model // spinner model
	list      list.Model    // list model
	choice    string        // target kubernetes cluster which user selected
	switching bool          // variable which shows state of changing context
}

func NewUIModel(items []list.Item, title string) *UIModel {

	l := list.New(items, ItemDelegate{}, constant.DefaultWidth, constant.ListHeight)

	l.SetShowStatusBar(false)
	l.SetShowPagination(true)

	l.Title = title

	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.Styles.FilterPrompt = helpStyle

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(constant.DefaultThemeColor))

	return &UIModel{spinner: sp, list: l}
}

func (m UIModel) Init() tea.Cmd {
	return nil
}

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				m.choice = string(i)
			}

			m.switching = true
			return m, tea.Batch(m.spinner.Tick, m.executeKubectl())
		}
	case kubectlResult:
		m.switching = false
		if msg.err != nil {
			fmt.Fprintf(os.Stderr, "Error switching context: %v\n", msg.err)
			return m, tea.Quit
		}
		return m, tea.Quit
	case tea.QuitMsg:
		return m, tea.Quit
	}

	var cmd tea.Cmd
	if m.switching {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m UIModel) View() string {
	if m.switching {
		switchingMsg := fmt.Sprintf("%s Switching Context to %s...", m.spinner.View(), m.choice)
		m.list.Title = switchingMsg
	}
	return m.list.View()
}
