package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	padding int = 1
)

var appStyle = lipgloss.NewStyle().
	PaddingTop(padding).
	PaddingBottom(padding).
	PaddingLeft(padding).
	PaddingRight(padding)

var boxHeaderStyle = lipgloss.NewStyle().
	Bold(true).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderBottom(false).
	BorderTop(true).
	BorderLeft(true).
 	BorderRight(true).
	BorderForeground(lipgloss.Color("#FAFAFA")).
	Foreground(lipgloss.Color("#FAFAFA")).
	PaddingLeft(padding).
	PaddingRight(padding)

var borderedBoxStyle = lipgloss.NewStyle().
	Bold(true).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#FAFAFA")).
	Foreground(lipgloss.Color("#FAFAFA")).
	PaddingTop(padding).
	PaddingBottom(padding).
	PaddingLeft(padding).
	PaddingRight(padding)

func createSystemRendering(m MtModel) string {
	var systemRendering string
	
	systemRendering += m.deviceInfo + "\n"

	for _, message := range m.data.resourceData.reply.Re {
		systemRendering += fmt.Sprintf("uptime: %s\n", message.Map["uptime"])
	}

	x, _ := borderedBoxStyle.GetFrameSize()
	systemRendering = borderedBoxStyle.Copy().Width(m.width-x).Render(systemRendering)
	return lipgloss.JoinVertical(lipgloss.Top, boxHeaderStyle.Render("System"), systemRendering)
}

func createResourceRendering(m MtModel) string {
	var resourceRendering string

	if m.data.cpuData.err != nil {
		resourceRendering = borderedBoxStyle.Render(fmt.Sprintf(
				"A problem has occured :(\n" +
				"%s\n",
				m.data.cpuData.err,	
		))
	} else {
		var cpuCoreInfo string

		x, _ := borderedBoxStyle.GetFrameSize()

		for i := 0; i < m.cpu.count; i++ {
			m.cpu.bar[i].Width = m.width / 2 - x

			cpuCoreInfo += fmt.Sprintf("core %d: ", i+1) + 
				m.cpu.bar[i].ViewAs(m.cpu.data[i])

			if i < m.cpu.count-1 {
				cpuCoreInfo += "\n\n"
			}
		}

		resourceRendering = borderedBoxStyle.Copy().Width(m.width / 2 - (x/2)).Render(cpuCoreInfo)
	}

	return lipgloss.JoinVertical(lipgloss.Top, boxHeaderStyle.Render("Resources"), resourceRendering)
}

func (m MtModel) View() string {
	var finalRendering string

	if m.data == nil {
		finalRendering = "Fetching data...\n"
	} else {
		systemRendering := createSystemRendering(m)
		resourceRendering := createResourceRendering(m)

		finalRendering += lipgloss.JoinVertical(lipgloss.Top, systemRendering, resourceRendering)
	}

	appStyle.Height(m.height)
	appStyle.Width(m.width)

	return appStyle.Render(fmt.Sprintf(
		"%s\n\n" +
		"Press q to exit...",
		finalRendering,
	))
}