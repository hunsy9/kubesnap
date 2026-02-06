package listview

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	c "github.com/hunsy9/kubesnap/pkg/constant"
	"k8s.io/client-go/tools/clientcmd"
)

func RenameOperation(oldName, newName string) tea.Cmd {
	return func() tea.Msg {
		if oldName == newName {
			return RenameResultMsg{Err: nil}
		}

		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		config, err := loadingRules.Load()
		if err != nil {
			return RenameResultMsg{Err: err}
		}

		if _, exists := config.Contexts[newName]; exists {
			return RenameResultMsg{Err: fmt.Errorf("context named '%s' already exists", newName)}
		}

		ctx, exists := config.Contexts[oldName]
		if !exists {
			return RenameResultMsg{Err: fmt.Errorf("context '%s' not found", oldName)}
		}

		config.Contexts[newName] = ctx
		delete(config.Contexts, oldName)

		if config.CurrentContext == oldName {
			config.CurrentContext = newName
		}

		if err := clientcmd.ModifyConfig(loadingRules, *config, true); err != nil {
			return RenameResultMsg{Err: err}
		}

		return RenameResultMsg{Err: nil}
	}
}

type RenamingOperation func(oldName, newName string) tea.Cmd

type RenameResultMsg struct {
	Err error
}

type RenamingModel struct {
	parent    tea.Model
	list      list.Model
	spinner   spinner.Model
	textInput textinput.Model
	choice    string
	operation RenamingOperation
	renaming  bool
}

func NewRenamingModel(parent tea.Model, items []list.Item, width int, op RenamingOperation) *RenamingModel {

	l := list.New(items, RenamingItemDelegate{}, width, c.ListHeight)
	l.Title = c.DefaultRenamingContextHeaderMessage

	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(true)
	l.KeyMap.CloseFullHelp.Unbind()
	l.KeyMap.ShowFullHelp.Unbind()

	sp := spinner.New()
	sp.Style = spinnerStyle

	ti := textinput.New()
	ti.Placeholder = "New Context Name"
	ti.CharLimit = 156
	ti.Width = 150

	return &RenamingModel{
		parent:    parent,
		list:      l,
		spinner:   sp,
		textInput: ti,
		operation: op,
	}
}

func (m *RenamingModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *RenamingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// 1. Text Input Mode
	if m.renaming {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				return m, m.operation(m.choice, m.textInput.Value())
			case "esc":
				m.renaming = false
				m.textInput.Blur()
				return m, nil
			}
		case RenameResultMsg:
			m.renaming = false
			return m.parent, func() tea.Msg { return msg }
		}

		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	// 2. List Selection Mode
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m.parent, nil

		case "enter":
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				m.choice = string(i.Name)
				m.textInput.SetValue(m.choice)
				m.textInput.Focus()
				m.textInput.CursorEnd()
				m.renaming = true
			}
			return m, nil
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *RenamingModel) View() string {
	listView := boxStyle.Render(m.list.View())

	if m.renaming {
		inputView := fmt.Sprintf(
			"\n  Rename %s to:\n  %s\n\n  %s",
			lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultThemeColor)).Bold(true).Render(m.choice),
			m.textInput.View(),
			lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultUrlColor)).Render(c.QuitRenameMode),
		)
		return lipgloss.JoinVertical(lipgloss.Left, listView, inputView)
	}

	return listView
}
