package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"	
	"code.cloudfoundry.org/bytefmt"
)

func createSystemRendering(m MtModel) string {
	var systemRendering string

	deviceHeader := subHeaderStyle.Render("Device")
	systemRendering += fmt.Sprintf("%s\n%s %s\n", deviceHeader, m.deviceInfo.Platform, m.deviceInfo.BoardName)

	systemRendering += fmt.Sprintf("%s: %s\n", "RouterOS", m.deviceInfo.OsVersion)

	systemRendering += fmt.Sprintf("%s: %s\n", "Uptime", m.resource.uptime)

	x, _ := borderedBoxStyle.GetFrameSize()
	systemRendering = borderedBoxStyle.Copy().Width(m.width / 2 - x).Render(systemRendering)
	return lipgloss.JoinVertical(lipgloss.Top, boxHeaderStyle.Render("System"), systemRendering)
}

func createResourceRendering(m MtModel) string {
	var resourceRendering string

	x, _ := borderedBoxStyle.GetFrameSize()

	width := m.width / 2 - x

	if m.cpu.err != nil {
		resourceRendering = borderedBoxStyle.Render(fmt.Sprintf(
				"A problem has occured :(\n" +
				"%s\n",
				m.cpu.err,	
		))
	} else {
		var cpuCoreInfo string

		for i := 0; i < m.cpu.count; i++ {
			m.cpu.bar[i].Width = width

			cpuCoreInfo += fmt.Sprintf("core %d: ", i+1) + 
				m.cpu.bar[i].ViewAs(m.cpu.data[i])

			if i < m.cpu.count-1 {
				cpuCoreInfo += "\n"
			}
		}
		
		cpuHeader := subHeaderStyle.Render("CPU") + "\n"
		cpuLoad := boxStyle.Copy().Width(width).Render(cpuHeader + cpuCoreInfo)

		var memory string

		memoryHeader := subHeaderStyle.Copy().MarginTop(1).Render("MEMORY") + "\n"
		m.resource.memoryBar.Width = width

		memory = fmt.Sprintf(
			"%s / %s\n%s", 
			bytefmt.ByteSize(m.resource.freeMem), 
			bytefmt.ByteSize(m.resource.totalMem),
			m.resource.memoryBar.ViewAs(float64(m.resource.freeMem) / float64(m.resource.totalMem)),
		)

		memory = boxStyle.Copy().Width(width).Render(memoryHeader + memory)		

		resourceRendering += borderedBoxStyle.Copy().Width(m.width / 2 - x).Render(cpuLoad + "\n" + memory)
	}

	return lipgloss.JoinVertical(lipgloss.Top, boxHeaderStyle.Render("Resources"), resourceRendering)
}

func (m MtModel) View() string {
	var finalRendering string

	if m.state == fetching {
		finalRendering = "Fetching data...\n"
	} else {
		systemRendering := createSystemRendering(m)
		resourceRendering := createResourceRendering(m)

		finalRendering += lipgloss.JoinHorizontal(lipgloss.Top, systemRendering, resourceRendering)
	}

	appStyle.Height(m.height)
	appStyle.Width(m.width)

	return appStyle.Render(fmt.Sprintf(
		"%s\n\n" +
		"Press q to exit...",
		finalRendering,
	))
}