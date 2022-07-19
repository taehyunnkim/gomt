package tui

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-routeros/routeros"
)


func fetchData(c *routeros.Client, sub chan dataMessage) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Second)

			var reply *routeros.Reply
			var err error
			
			reply, err = c.RunArgs([]string{"/system/resource/print"})
			resourceData := routerFetchData{reply, err}

			reply, err = c.RunArgs([]string{"/system/resource/cpu/print"})
			cpuData := routerFetchData{reply, err}

			reply, err = c.RunArgs([]string{"/system/health/print"})
			healthData := routerFetchData{reply, err}

			sub <- dataMessage{
				resourceData,
				cpuData,
				healthData,
			}
		}
	}
}

func waitForMessage(sub chan dataMessage) tea.Cmd {
	return func() tea.Msg {
		return dataMessage(<-sub)
	}
}

func parseResourceData(data routerFetchData, m *resourceData) {
	if data.err != nil {
		m.err = data.err
	}

	if len(data.reply.Re) > 0 {
		message := data.reply.Re[0]

		m.uptime = message.Map["uptime"]

		freeMem, err := strconv.ParseUint(message.Map["free-memory"], 10, 64)
		if err == nil {
			m.freeMem = freeMem
		}
		
		totalMem, err := strconv.ParseUint(message.Map["total-memory"], 10, 64)
		if err == nil {
			m.totalMem = totalMem
		}
		
		freeHdd, err := strconv.ParseUint(message.Map["free-hdd-space"], 10, 64)
		if err == nil {
			m.freeHdd = freeHdd
		}
		
		totalHdd, err := strconv.ParseUint(message.Map["total-hdd-space"], 10, 64)
		if err == nil {
			m.totalHdd = totalHdd
		}
	}
}

func parseCpuData(data routerFetchData, m *cpuData) {
	if data.err != nil {
		m.err = data.err
	}

	for _, message := range data.reply.Re {
		var cpu int
		_, err := fmt.Sscanf(message.Map[".id"], "*%d", &cpu)

		if err != nil {
			log.Fatal(err)	
		}

		value, _ := strconv.ParseFloat(message.Map["load"], 64)
		m.data[cpu] = value / 100
	}
}

func parseHealthData(data routerFetchData, m *healthData) {
	if data.err != nil {
		m.err = data.err
	}

	for i, message := range data.reply.Re {
		value := fmt.Sprintf("%s%s", message.Map["value"], message.Map["type"])
		m.data[i] = health{
			message.Map["name"],
			value,
		}
	}
}

func (m MtModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		if message.Width > m.minWidth {
			m.width = message.Width
		}
	
		m.height = message.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(message, m.keys.Quit):
			return m, tea.Quit
		default:
			return m, nil
		}
	case dataMessage:
		parseResourceData(message.resourceData, m.resource)
		parseCpuData(message.cpuData, m.cpu)
		parseHealthData(message.healthData, m.health)
		m.state = ready

		return m, waitForMessage(m.sub)
	}

	return m, nil
}