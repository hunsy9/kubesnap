package listview

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	c "github.com/hunsy9/kubesnap/pkg/constant"
	"k8s.io/client-go/tools/clientcmd"
)

func DeleteOperation(targets []string) tea.Cmd {
	return func() tea.Msg {
		config, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)

		if err != nil {
			return DeletionResultMsg{Err: err, Targets: nil}
		}

		for _, targetCtx := range targets {
			if _, exists := config.Contexts[targetCtx]; exists {
				delete(config.Contexts, targetCtx)
			}
		}

		err = clientcmd.WriteToFile(*config, clientcmd.RecommendedHomeFile)
		if err != nil {
			return DeletionResultMsg{Err: err, Targets: nil}
		}

		time.Sleep(time.Millisecond * 500)

		return DeletionResultMsg{Err: nil, Targets: targets}
	}
}

type DeletingOperation func(targets []string) tea.Cmd

type DeletionResultMsg struct {
	Err     error
	Targets []string
}

type DeletingModel struct {
	parent     tea.Model
	list       list.Model
	spinner    spinner.Model
	selected   map[string]struct{}
	operation  DeletingOperation
	textInput  textinput.Model
	deleting   bool
	confirming bool
	errorMsg   string
}

func NewDeletingModel(parent tea.Model, items []list.Item, width int, op DeletingOperation) *DeletingModel {
	selected := make(map[string]struct{})
	delegate := DeletingItemDelegate{Selected: selected}

	l := list.New(items, delegate, width, c.ListHeight+2)
	l.Title = c.DefaultDeletingContextHeaderMessage

	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(true)
	l.KeyMap.CloseFullHelp.Unbind()
	l.KeyMap.ShowFullHelp.Unbind()

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		}
	}

	ti := textinput.New()
	ti.Placeholder = "yes"
	ti.CharLimit = 156
	ti.Width = 50

	sp := spinner.New()
	sp.Style = spinnerStyle

	return &DeletingModel{
		parent:    parent,
		list:      l,
		spinner:   sp,
		selected:  selected,
		operation: op,
		textInput: ti,
	}
}

func (m *DeletingModel) Init() tea.Cmd {
	return nil
}

func (m *DeletingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	// 1. Deleting confirmation mode
	if m.confirming {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.confirming = false
				m.textInput.Blur()
				return m, nil
			case "enter":
				if m.textInput.Value() == "yes" {
					targets := make([]string, 0, len(m.selected))
					for k := range m.selected {
						targets = append(targets, k)
					}
					m.deleting = true
					m.confirming = false
					return m, tea.Batch(m.spinner.Tick, m.operation(targets))
				} else {
					m.errorMsg = "Incorrect Input. Type 'yes' to confirm"
					m.textInput.SetValue("")
				}
				return m, nil
			}
			m.errorMsg = ""
		}

		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	// 2. Deleting list mode
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case tea.KeyMsg:
		if m.deleting {
			return m, nil
		}

		switch msg.String() {
		case "q", "esc":
			return m.parent, nil

		case " ":
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				if i.Name != i.DisplayName { // it is current context
					m.errorMsg = "  You can't delete current context."
					return m, nil
				}
				val := string(i.Name)
				if _, exists := m.selected[val]; exists {
					delete(m.selected, val)
				} else {
					m.selected[val] = struct{}{}
				}
			}
			return m, nil

		case "enter":
			if len(m.selected) > 0 {
				m.confirming = true
				m.textInput.Focus()
			}
		}

		m.errorMsg = ""

	case DeletionResultMsg:
		m.deleting = false
		return m.parent, func() tea.Msg { return msg }
	}

	if m.deleting {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *DeletingModel) View() string {
	if m.deleting {
		return fmt.Sprintf("\n\n   %s Deleting %d items from Kubeconfig\n\n", m.spinner.View(), len(m.selected))
	}

	if m.confirming {
		targets := make([]string, 0, len(m.selected))
		for k := range m.selected {
			targets = append(targets, "• "+k)
		}
		targetstring := strings.Join(targets, "\n ")

		inputView := fmt.Sprintf(
			"\n You are about to delete %v contexts: \n\n %s \n\n This action cannot be undone. \n Type \"%s\" to confirm\n\n  %s\n\n %s\n",
			len(targets),
			targetContextStyle.Render(targetstring),
			activeStyle.Render("yes"),
			m.textInput.View(),
			quitRenameModeFooterStyle.Render(c.QuitRenameMode),
		)

		if m.errorMsg != "" {
			inputView += "\n " + errorStyle.Render(m.errorMsg) + "\n"
		}

		return boxStyle.Render(inputView)
	}

	layout := m.list.View()

	if m.errorMsg != "" {
		layout += "\n " + errorStyle.Render(m.errorMsg) + "\n"
	}

	return boxStyle.Render(layout)
}
