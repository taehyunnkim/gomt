package tui

import (
	"fmt"
	"strings"
	
	"github.com/taehyunnkim/gomt/internal/math"

	"code.cloudfoundry.org/bytefmt"
	"github.com/charmbracelet/lipgloss"
)

func getMaxWidth(prev int, s string) int {
	length := len(s)

	if prev < len(s) {
		return length
	}

	return prev
}

func createSystemRendering(m MtModel) string {
	var systemRendering string
	boxFrame := borderedBoxStyle.GetHorizontalBorderSize() + borderedBoxStyle.GetHorizontalMargins()
	systemRenderingWidth := m.width - appStyle.GetHorizontalFrameSize() - boxFrame

	var deviceRendering string
	var deviceMinWidth int = 0
	deviceRendering += fmt.Sprintf("%s\n", subHeaderStyle.Render("Device"))
	deviceName := fmt.Sprintf("%s %s\n", m.deviceInfo.Platform, m.deviceInfo.BoardName)
	deviceRendering += deviceName

	deviceMinWidth = getMaxWidth(deviceMinWidth, deviceName)

	routerOs := fmt.Sprintf("%s: %s\n", "RouterOS", m.deviceInfo.OsVersion)
	deviceRendering += routerOs
	deviceMinWidth = getMaxWidth(deviceMinWidth, routerOs)

	deviceUptime := fmt.Sprintf("%s: %s\n", "Uptime", m.resource.uptime)
	deviceRendering += deviceUptime
	deviceMinWidth = getMaxWidth(deviceMinWidth, deviceUptime)

	deviceBox := boxStyle.Copy()
	if m.width % 2 != 0 { deviceBox.MarginRight(1) }
	
	var healthRendering string
	healthRendering += fmt.Sprintf("%s\n", subHeaderStyle.Render("Health"))
	for _, health := range m.health.data {
		if health.name != "" {
			healthRendering += fmt.Sprintf("%s: %s\n", health.name, health.value)
		}
	}

	if deviceMinWidth >= systemRenderingWidth/2 {
		deviceRendering = deviceBox.Width(systemRenderingWidth).Render(deviceRendering)
		healthRendering = boxStyle.Copy().Width(systemRenderingWidth).Render(healthRendering)
		systemRendering = lipgloss.JoinVertical(lipgloss.Top, deviceRendering, healthRendering)
	} else {
		deviceRendering = deviceBox.Width(systemRenderingWidth/2).Render(deviceRendering)
		healthRendering = boxStyle.Copy().Width(systemRenderingWidth/2).Render(healthRendering)
		systemRendering = lipgloss.JoinHorizontal(lipgloss.Top, deviceRendering, healthRendering)
	}

	systemRendering = borderedBoxStyle.Copy().Width(systemRenderingWidth).Render(systemRendering) 
	
	return lipgloss.JoinVertical(lipgloss.Top, boxHeaderStyle.Render("System"), systemRendering)
}

func createResourceRendering(m MtModel) string {
	var resourceRendering string

	if m.cpu.err != nil {
		resourceRendering = borderedBoxStyle.Render(fmt.Sprintf(
				"A problem has occured :(\n" +
				"%s\n",
				m.cpu.err,	
		))
	} else {
		horizontalFrameSize := borderedBoxStyle.GetHorizontalFrameSize()
	
		var cpuCores string
		
		// Calculate Resource Window Width
		windowWidth := m.width - horizontalFrameSize - appStyle.GetHorizontalFrameSize()


		// Calculate number of cpu core bars in a row
		var numBarsInRow int

		if num := windowWidth / m.cpu.minBarWidth; num % 2 == 0 {
			numBarsInRow = math.Min(m.deviceInfo.CpuCoreCount, num)
		} else {
			numBarsInRow = math.Min(m.deviceInfo.CpuCoreCount, num+1)
		}

		// Calculate bar width
		// 2: core number
		// 2: bar border
		// 4: percentage
		var barWidth int = windowWidth / numBarsInRow - 2 - 2 - 4
		var remainingGap = windowWidth - (barWidth + 2 + 2 + 4) * numBarsInRow

		var newLines int = 0
		for i := 0; i < m.cpu.count; i++ {
			var gap string
			m.cpu.bar[i].Width = barWidth

			if (i+1) % numBarsInRow == 0 {
				gap  = "\n"
				newLines++
			} else {
				var start int = numBarsInRow / 2
				
				if (i+1) == start + (numBarsInRow * newLines) {
					gap = strings.Repeat(" ", remainingGap)
				} 
			}


			cpuCoreInfo := fmt.Sprintf(
				"%2d[%s]%3.0f%%%s", 
				i+1, 
				m.cpu.bar[i].ViewAs(m.cpu.data[i]), 
				m.cpu.data[i]*100,
				gap,
			)

			cpuCores += cpuCoreInfo
		}	

		cpuHeader := subHeaderStyle.Render("CPU") + "\n"
		cpuLoad := boxStyle.Copy().Width(windowWidth).Render(cpuHeader + cpuCores)

		var memory string
		memoryHeader := subHeaderStyle.Copy().MarginTop(1).Render("MEMORY") + "\n"
		m.resource.memoryBar.Width = windowWidth

		memory = fmt.Sprintf(
			"%s / %s\n%s", 
			bytefmt.ByteSize(m.resource.freeMem), 
			bytefmt.ByteSize(m.resource.totalMem),
			m.resource.memoryBar.ViewAs(float64(m.resource.freeMem) / float64(m.resource.totalMem)),
		)

		memory = boxStyle.Copy().Width(windowWidth).Render(memoryHeader + memory)		

		resourceRendering += borderedBoxStyle.Copy().Width(m.width - horizontalFrameSize).Render(cpuLoad + memory)
	}

	return lipgloss.JoinVertical(lipgloss.Top, boxHeaderStyle.Render("Resources"), resourceRendering)
}

func (m MtModel) View() string {
	var contentRendering string

	if m.state == fetching {
		contentRendering = "Fetching data...\n"
	} else {
		systemRendering := createSystemRendering(m)
		resourceRendering := createResourceRendering(m)

		contentRendering += lipgloss.JoinVertical(lipgloss.Top, systemRendering, resourceRendering)
	}

	contentRendering = contentStyle.Render(contentRendering)
	contentHeight := strings.Count(contentRendering, "\n")
	
	appStyle.Width(m.width)
	appStyle.Height(m.height)


	helpView := m.help.View(m.keys)
	blankHeight := m.height - contentHeight - strings.Count(helpView, "\n") - 3

	if blankHeight < 0 {
		blankHeight = 1
	}

	app := fmt.Sprintf(
		"%s%s%s",
		contentRendering,
		strings.Repeat("\n", blankHeight),
		helpView,
	)
	
	var debugString string
	
	if m.debug {
		debugString = fmt.Sprintf(
			"%d:%d",
			m.cpu.bar[0].Width + 8,
			m.width - 6,
		)
	}

	return appStyle.Render(app + debugString)
}