package model

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hunsy9/kubesnap/pkg/constant"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Background(lipgloss.Color("173"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("173"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type UIModel struct {
	spinner  spinner.Model
	list     list.Model
	choice   string
	quitting bool
}

func NewUIModel(items []list.Item, title string) *UIModel {
	l := list.New(items, ItemDelegate{}, constant.DefaultWidth, constant.ListHeight)

	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Title = title
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("206"))

	return &UIModel{spinner: sp, list: l}
}

func (m UIModel) Init() tea.Cmd {
	return nil
}
