package listview

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type DeletingItemDelegate struct {
	Selected map[string]struct{}
}

func (d DeletingItemDelegate) Height() int                             { return 1 }
func (d DeletingItemDelegate) Spacing() int                            { return 0 }
func (d DeletingItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d DeletingItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}
	str := string(i)

	_, selected := d.Selected[str]
	checkbox := "[ ]"
	if selected {
		checkbox = "[X]"
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(s...)
		}
	}

	fmt.Fprint(w, fn(fmt.Sprintf("%s %s", checkbox, str)))
}
