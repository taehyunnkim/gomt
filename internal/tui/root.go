package tui

import (
	"fmt"
	"log"
	"time"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-routeros/routeros"
)


func fetchData(c *routeros.Client, sub chan dataMessage) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Second)
			
			reply, err := c.RunArgs([]string{"/system/resource/print"})
			resourceData := data{reply, err}

			reply2, err2 := c.RunArgs([]string{"/system/resource/cpu/print"})
			cpuData := data{reply2, err2}


			sub <- dataMessage{
				resourceData,
				cpuData,
			}
		}
	}
}

func waitForMessage(sub chan dataMessage) tea.Cmd {
	return func() tea.Msg {
		return dataMessage(<-sub)
	}
}

func (m MtModel) Init() tea.Cmd {
	return tea.Batch(
		fetchData(m.client, m.sub),
		waitForMessage(m.sub),
	)
}

func (m MtModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}
	case dataMessage:
		m.data = &message
		return m, waitForMessage(m.sub)
	}

	return m, nil
}

func (m MtModel) View() string {
	var result string

	if m.data == nil {
		result = "Fetching data...\n"
	} else {
		result += m.deviceInfo + "\n"

		for _, message := range m.data.resourceData.reply.Re {
			result += fmt.Sprintf("uptime: %s\n", message.Map["uptime"])
		}	

		result += "\n===== CPU =====\n"

		if m.data.cpuData.err != nil {
			result = fmt.Sprintf(
				"A problem has occured :(\n" +
				"%s\n",
				m.data.cpuData.err,	
			)
		} else {
			parseCpuData(m.data.cpuData, &m.cpu.data)
			for i := 0; i < m.cpu.count; i++ {
				result += fmt.Sprintf("cpu-%d: ", i+1) + 
					m.cpu.bar[i].ViewAs(m.cpu.data[i]) + "\n"
			}
		}
	}

	return fmt.Sprintf(
		"%s\n" +
		"Press q to exit...",
		result,	
	)
}

func parseCpuData(data data, m *map[int] float64) {
	for _, message := range data.reply.Re {
		var cpu int
		_, err := fmt.Sscanf(message.Map[".id"], "*%d", &cpu)

		if err != nil {
			log.Fatal(err)	
		}

		value, _ := strconv.ParseFloat(message.Map["load"], 64)
		(*m)[cpu] = value / 100
	}
}