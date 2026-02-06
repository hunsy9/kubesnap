package listview

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	c "github.com/hunsy9/kubesnap/pkg/constant"
)

type OperationResultMsg struct {
	Err error
}

type SwitchingOperation func(param string) tea.Cmd // switching function executed by UIModel

type SwitchingUIModel struct {
	spinner   spinner.Model      // spinner model
	list      list.Model         // list model
	tag       string             // operation tag
	choice    string             // target kubernetes cluster which user selected
	switching bool               // variable which shows state of changing context/namespace
	operation SwitchingOperation // switching function executed by UIModel
	output    string             // success output message for context switching
}

func NewSwitchingUIModel(items []list.Item, title string, tag string) *SwitchingUIModel {

	l := list.New(items, ItemDelegate{}, c.DefaultWidth, c.ListHeight)

	l.SetShowStatusBar(false)
	l.SetShowPagination(true)
	l.SetShowHelp(true)
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

	return &SwitchingUIModel{spinner: sp, list: l, tag: tag}
}

func (m *SwitchingUIModel) Init() tea.Cmd {
	return nil
}

func (m *SwitchingUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, tea.Batch(m.spinner.Tick, m.operation(m.choice))
		case "d":
			if m.tag == c.Namespace {
				break
			}
			currentWidth := m.list.Width()
			items := m.list.Items()

			return NewDeletingModel(m, items, currentWidth, DeleteOperation), nil
		case "r":
			if m.tag == c.Namespace {
				break
			}
			currentWidth := m.list.Width()
			items := m.list.Items()

			return NewRenamingModel(m, items, currentWidth, RenameOperation), nil
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
			m.output = fmt.Sprintf("Deleted %d contexts.\n", msg.Count)
		}
		return m, tea.Quit

	case RenameResultMsg:
		if msg.Err != nil {
			m.output = fmt.Sprintf("Error renaming: %v\n", msg.Err)
		} else {
			m.output = "Renamed context successfully.\n"
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

func (m *SwitchingUIModel) View() string {
	if m.switching {
		switchingMsg := fmt.Sprintf("%s Switching %s to %s...", m.spinner.View(), m.tag, m.choice)
		m.list.Title = switchingMsg
	}

	return boxStyle.Render(m.list.View())
}

func (m *SwitchingUIModel) SetOperationFunc(operation SwitchingOperation) {
	m.operation = operation
}

func (m *SwitchingUIModel) GetOutput() string {
	return m.output
}
