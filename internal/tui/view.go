package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	padding int = 1
)

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

var box = lipgloss.NewStyle().
	Bold(true).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#FAFAFA")).
	Foreground(lipgloss.Color("#FAFAFA")).
	PaddingTop(padding).
	PaddingBottom(padding).
	PaddingLeft(padding).
	PaddingRight(padding)

func createResourceRendering(m MtModel) string {
	var resourceRendering string

	if m.data.cpuData.err != nil {
		resourceRendering = box.Render(fmt.Sprintf(
				"A problem has occured :(\n" +
				"%s\n",
				m.data.cpuData.err,	
		))
	} else {
		var cpuCoreInfo string
		for i := 0; i < m.cpu.count; i++ {
			cpuCoreInfo += fmt.Sprintf("core %d: ", i+1) + 
				m.cpu.bar[i].ViewAs(m.cpu.data[i])

			if i < m.cpu.count-1 {
				cpuCoreInfo += "\n\n"
			}
		}

		resourceRendering = box.Render(cpuCoreInfo)
	}

	return lipgloss.JoinVertical(lipgloss.Top, boxHeaderStyle.Render("Resources"), resourceRendering)
}

func (m MtModel) View() string {
	var rendering string

	if m.data == nil {
		rendering = "Fetching data...\n"
	} else {
		rendering += m.deviceInfo + "\n"

		for _, message := range m.data.resourceData.reply.Re {
			rendering += fmt.Sprintf("uptime: %s\n", message.Map["uptime"])
		}

		rendering += createResourceRendering(m)
	}

	return fmt.Sprintf(
		"%s\n" +
		"Press q to exit...",
		rendering,	
	)
}