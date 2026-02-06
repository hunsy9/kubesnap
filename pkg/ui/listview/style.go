package listview

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	c "github.com/hunsy9/kubesnap/pkg/constant"
)

var (
	// Common listview style
	boxStyle          = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(c.DefaultThemeColor)).Padding(0, 2)
	titleStyle        = lipgloss.NewStyle().MarginTop(1).Bold(true).Italic(true).Foreground(lipgloss.Color(c.DefaultThemeColor))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Bold(true).Foreground(lipgloss.Color(c.DefaultThemeColor))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(2)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(2).PaddingBottom(1)
	spinnerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultThemeColor))

	// Renaming listview style
	renameTargetContextStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultThemeColor)).Bold(true)
	quitRenameModeFooterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultUrlColor))
)
