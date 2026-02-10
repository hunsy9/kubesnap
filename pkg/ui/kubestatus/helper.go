package kubestatus

import tea "github.com/charmbracelet/bubbletea"

// helper func: return tea.Msg when async tasks done
func (m *StatusModel) incrementProgress() tea.Cmd {
	m.completedTasks++
	if m.completedTasks >= m.totalAsyncTasks {
		return func() tea.Msg {
			return AsyncTaskAllDoneMsg{}
		}
	}
	return nil
}

// helper func: render spinner if is_Loading variable is true
func (m *StatusModel) renderLoadingOrContent(loading bool, content string) string {
	if loading {
		return m.spinner.View()
	}
	return content
}
