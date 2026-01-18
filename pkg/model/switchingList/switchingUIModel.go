package model

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hunsy9/kubesnap/pkg/constant"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginTop(1).MarginLeft(2).Bold(true).Italic(true).Foreground(lipgloss.Color(constant.DefaultThemeColor))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color(constant.DefaultThemeColor))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(3)
	spinnerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(constant.DefaultThemeColor))
)

type OperationResultMsg struct {
	Err error
}

type SwitchingOperation func(param string) tea.Cmd // switching function executed by UIModel

type UIModel struct {
	spinner   spinner.Model      // spinner model
	list      list.Model         // list model
	tag       string             // operation tag
	choice    string             // target kubernetes cluster which user selected
	switching bool               // variable which shows state of changing context/namespace
	operation SwitchingOperation // switching function executed by UIModel
	output    string             // success output message for context switching
}

func NewUIModel(items []list.Item, title string, tag string) *UIModel {

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
	sp.Style = spinnerStyle

	return &UIModel{spinner: sp, list: l, tag: tag}
}

func (m *UIModel) Init() tea.Cmd {
	return nil
}

func (m *UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, tea.Batch(m.spinner.Tick, m.operation(m.choice))
		}
	case OperationResultMsg:
		m.switching = false
		if msg.Err != nil {
			fmt.Fprintf(os.Stderr, "Error switching %s: %v\n", m.tag, msg.Err)
			return m, tea.Quit
		}
		m.output = fmt.Sprintf("Switched to %s: %s\n", m.tag, m.choice)
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

func (m *UIModel) View() string {
	if m.switching {
		switchingMsg := fmt.Sprintf("%s Switching %s to %s...", m.spinner.View(), m.tag, m.choice)
		m.list.Title = switchingMsg
	}
	return m.list.View()
}

func (m *UIModel) SetOperationFunc(operation SwitchingOperation) {
	m.operation = operation
}

func (m *UIModel) GetOutput() string {
	return m.output
}
