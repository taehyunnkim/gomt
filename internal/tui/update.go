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

func (m MtModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		m.height = message.Height
	case tea.KeyMsg:
		switch message.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}
	case dataMessage:
		m.data = &message
		parseCpuData(m.data.cpuData, &m.cpu.data)
		return m, waitForMessage(m.sub)
	}

	return m, nil
}