package kubestatus

import (
	"github.com/charmbracelet/lipgloss"
	c "github.com/hunsy9/kubesnap/pkg/constant"
)

var (
	spinnerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultThemeColor))
	boxStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(c.DefaultThemeColor)).Padding(0, 2).Margin(0, 0, 1, 0)
	headerStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultThemeColor)).Bold(true).Italic(true)
	urlStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultUrlColor))
	activeStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultActiveColor))
	warningStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultWarningColor))
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultErrorColor))
	updageMsgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultUpdateMsgColor))
	iconStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(c.DefaultThemeColor))
)
