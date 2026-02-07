package listview

import (
	"fmt"
	"io"

	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type RenamingItemDelegate struct {
	Selected map[string]struct{}
}

func (d RenamingItemDelegate) Height() int                             { return 1 }
func (d RenamingItemDelegate) Spacing() int                            { return 0 }
func (d RenamingItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d RenamingItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.DisplayName)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("• " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
