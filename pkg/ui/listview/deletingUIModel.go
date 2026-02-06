package listview

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	c "github.com/hunsy9/kubesnap/pkg/constant"
	"k8s.io/client-go/tools/clientcmd"
)

func DeleteOperation(targets []string) tea.Cmd {
	return func() tea.Msg {
		config, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)

		if err != nil {
			return DeletionResultMsg{Err: err, Count: 0}
		}

		for _, targetCtx := range targets {
			if _, exists := config.Contexts[targetCtx]; exists {
				delete(config.Contexts, targetCtx)
			}
		}

		err = clientcmd.WriteToFile(*config, clientcmd.RecommendedHomeFile)
		if err != nil {
			return DeletionResultMsg{Err: err, Count: 0}
		}

		return DeletionResultMsg{Err: nil, Count: len(targets)}
	}
}

type DeletingOperation func(targets []string) tea.Cmd

type DeletionResultMsg struct {
	Err   error
	Count int
}

type DeletingModel struct {
	parent    tea.Model
	list      list.Model
	spinner   spinner.Model
	selected  map[string]struct{}
	operation DeletingOperation
	deleting  bool
}

func NewDeletingModel(parent tea.Model, items []list.Item, width int, op DeletingOperation) *DeletingModel {
	selected := make(map[string]struct{})
	delegate := DeletingItemDelegate{Selected: selected}

	l := list.New(items, delegate, width, c.ListHeight)
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
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "toggle")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		}
	}

	sp := spinner.New()
	sp.Style = spinnerStyle

	return &DeletingModel{
		parent:    parent,
		list:      l,
		spinner:   sp,
		selected:  selected,
		operation: op,
	}
}

func (m *DeletingModel) Init() tea.Cmd {
	return nil
}

func (m *DeletingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		case "enter", " ":
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				val := string(i.Name)
				if _, exists := m.selected[val]; exists {
					delete(m.selected, val)
				} else {
					m.selected[val] = struct{}{}
				}
			}
			return m, nil

		case "d":
			if len(m.selected) > 0 {
				m.deleting = true
				targets := make([]string, 0, len(m.selected))
				for k := range m.selected {
					targets = append(targets, k)
				}
				return m, tea.Batch(m.spinner.Tick, m.operation(targets))
			}
		}

	case DeletionResultMsg:
		m.deleting = false
		return m.parent, func() tea.Msg { return msg }
	}

	var cmd tea.Cmd
	if m.deleting {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *DeletingModel) View() string {
	if m.deleting {
		return fmt.Sprintf("\n\n   %s Deleting %d items from %s...\n\n", m.spinner.View(), len(m.selected), "")
	}
	return boxStyle.Render(m.list.View())
}
