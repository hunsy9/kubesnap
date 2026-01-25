package switchingList

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/hunsy9/kubesnap/pkg/constant"
)

var (
	boxStyle          = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(constant.DefaultThemeColor)).Padding(0, 2).Margin(0, 0, 1, 0)
	titleStyle        = lipgloss.NewStyle().MarginTop(1).Bold(true).Italic(true).Foreground(lipgloss.Color(constant.DefaultThemeColor))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Bold(true).Foreground(lipgloss.Color(constant.DefaultThemeColor))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(2)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(2).PaddingBottom(1)
	spinnerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(constant.DefaultThemeColor))
)
