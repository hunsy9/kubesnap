package listview

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	c "github.com/hunsy9/kubesnap/pkg/constant"
	"k8s.io/client-go/tools/clientcmd"
)

func (m *SwitchingUIModel) RenameOperation(newName string) tea.Cmd {
	return func() tea.Msg {
		if m.choice == newName {
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

		ctx, exists := config.Contexts[m.choice]
		if !exists {
			return RenameResultMsg{Err: fmt.Errorf("context '%s' not found", m.choice)}
		}

		config.Contexts[newName] = ctx
		delete(config.Contexts, m.choice)

		if config.CurrentContext == m.choice {
			config.CurrentContext = newName
		}

		if err := clientcmd.ModifyConfig(loadingRules, *config, true); err != nil {
			return RenameResultMsg{Err: err}
		}

		return RenameResultMsg{OriginalName: m.choice, NewName: newName, Err: nil}
	}
}

type RenameResultMsg struct {
	OriginalName string
	NewName      string
	Err          error
}

type OperationResultMsg struct {
	Err error
}

type SwitchingOperation func(param string) tea.Cmd // switching function executed by UIModel

type SwitchingUIModel struct {
	spinner            spinner.Model // spinner model
	list               list.Model    // list model
	title              string
	tag                string             // operation tag
	choice             string             // target kubernetes cluster which user selected
	switching          bool               // variable which shows state of changing context/namespace
	renaming           bool               // variable which shows state of renaming context
	switchingOperation SwitchingOperation // switching function executed by UIModel
	output             string             // success output message for context switching
	textInput          textinput.Model
}

func NewSwitchingUIModel(items []list.Item, title string, tag string) *SwitchingUIModel {

	l := list.New(items, ItemDelegate{}, c.DefaultWidth, c.ListHeight)

	l.SetShowStatusBar(false)
	l.SetShowPagination(true)
	l.KeyMap.CloseFullHelp.Unbind()
	l.KeyMap.ShowFullHelp.Unbind()
	l.KeyMap.Filter.SetHelp("/", "search")

	l.Title = title

	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.Styles.FilterPrompt = helpStyle

	if tag == c.Context {
		l.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{
				key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rename")),
				key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
			}
		}
	}

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = spinnerStyle

	ti := textinput.New()
	ti.Placeholder = "New Context Name"
	ti.CharLimit = 156
	ti.Width = 150

	return &SwitchingUIModel{spinner: sp, list: l, tag: tag, textInput: ti, title: title}
}

func (m *SwitchingUIModel) Init() tea.Cmd {
	return nil
}

func (m *SwitchingUIModel) IsCustomKeyEnablementInValid() bool {
	return m.tag == c.Namespace || m.list.FilterState() != list.Unfiltered || m.switching
}

func (m *SwitchingUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// 1. Text Input(Rename) Mode
	if m.renaming {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				return m, m.RenameOperation(m.textInput.Value())
			case "esc":
				m.renaming = false
				m.textInput.Blur()
				m.list.Title = m.title
				return m, nil
			}
		case RenameResultMsg:
			m.renaming = false
			return m, func() tea.Msg { return msg }
		}

		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	// 2. Context/Namespace Switching Mode
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case tea.KeyMsg:

		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				m.choice = string(i.Name)
			}

			m.switching = true
			return m, tea.Batch(m.spinner.Tick, m.switchingOperation(m.choice))
		case "d":
			if m.IsCustomKeyEnablementInValid() {
				break
			}
			currentWidth := m.list.Width()
			items := m.list.Items()

			return NewDeletingModel(m, items, currentWidth, DeleteOperation), nil
		case "r":
			if m.IsCustomKeyEnablementInValid() {
				break
			}

			i, ok := m.list.SelectedItem().(Item)
			if ok {
				m.choice = string(i.Name)
				m.textInput.SetValue(m.choice)
				m.textInput.Focus()
				m.textInput.CursorEnd()
				m.renaming = true
				m.list.Title = "Rename Context"
			}

			return m, nil
		}
	case OperationResultMsg:
		m.switching = false
		if msg.Err != nil {
			fmt.Fprintf(os.Stderr, "Error switching %s: %v\n", m.tag, msg.Err)
			return m, tea.Quit
		}
		m.output = fmt.Sprintf("Switched to %s: %s\n", m.tag, m.choice)
		return m, tea.Quit

	case DeletionResultMsg:
		if msg.Err != nil {
			m.output = fmt.Sprintf("Error deleting %v\n", msg.Err)
		} else {
			count := len(msg.Targets)
			targets := make([]string, 0, count)
			for _, k := range msg.Targets {
				targets = append(targets, "• "+k)
			}
			targetstring := strings.Join(targets, "\n")
			m.output = fmt.Sprintf("Deleted following %v contexts.\n%s\n", count, targetstring)
		}
		return m, tea.Quit

	case RenameResultMsg:
		if msg.Err != nil {
			m.output = fmt.Sprintf("Error renaming: %v\n", msg.Err)
		} else {
			m.output = fmt.Sprintf("Renamed context successfully.\nbefore: %s\nafter: %s\n", msg.OriginalName, msg.NewName)
		}
		return m, tea.Quit

	case tea.QuitMsg:
		return m, tea.Quit
	}

	if m.switching {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *SwitchingUIModel) View() string {
	listView := boxStyle.Render(m.list.View())

	if m.switching {
		switchingMsg := fmt.Sprintf("%s Switching %s to %s...", m.spinner.View(), m.tag, m.choice)
		m.list.Title = switchingMsg
	}

	if m.renaming {
		inputView := fmt.Sprintf(
			"\n  Rename %s to:\n  %s\n\n  %s",
			targetContextStyle.Render(m.choice),
			m.textInput.View(),
			quitRenameModeFooterStyle.Render(c.QuitRenameMode),
		)
		return lipgloss.JoinVertical(lipgloss.Left, listView, inputView)
	}

	return listView
}

func (m *SwitchingUIModel) SetOperationFunc(operation SwitchingOperation) {
	m.switchingOperation = operation
}

func (m *SwitchingUIModel) GetOutput() string {
	return m.output
}
