package tui

import tea "github.com/charmbracelet/bubbletea"

func (m MtModel) Init() tea.Cmd {
	return tea.Batch(
		fetchData(m.client, m.sub),
		waitForMessage(m.sub),
	)
}