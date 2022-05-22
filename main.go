package main

import (
	"fmt"
	"log"
	"syscall"
	"time"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/go-routeros/routeros"
	"golang.org/x/term"
)

var (
	version string = "v0.1.1"
)

type model struct {
	deviceInfo string
	client *routeros.Client
	sub chan dataMessage
	data *dataMessage
	cpu cpuData
}

type cpuData struct {
	count int
	bar map[int] progress.Model
	data map[int] float64
}

type dataMessage struct {
	resourceData data
	cpuData data 
}

type data struct {
	reply *routeros.Reply
	err error
}

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

func (m model) Init() tea.Cmd {
	return tea.Batch(
		fetchData(m.client, m.sub),
		waitForMessage(m.sub),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m model) View() string {
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

func main() {
	fmt.Println("gomt", version)

	var address string
	fmt.Print("Enter the full address: ")
	fmt.Scanf("%s", &address)

	var user string
	fmt.Print("Enter the user: ")
	fmt.Scanf("%s", &user)

	fmt.Print("Enter the password: ")
	password, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	client, err := routeros.Dial(address, user, string(password))

	if err!= nil {
		log.Fatal(err)
		return
	}

	defer client.Close()

	reply, err := client.RunArgs([]string{"/system/resource/print"})

	if err != nil {
		log.Fatal(err)
	}

	if len(reply.Re) > 0 {
		platform := reply.Re[0].Map["platform"]
		boardName := reply.Re[0].Map["board-name"]
		osVersion := reply.Re[0].Map["version"]
		cpuCoreCount, _ := strconv.Atoi(reply.Re[0].Map["cpu-count"])

		deviceInfo := fmt.Sprintf("%s %s | RouterOs %s | %s | %s\n", platform, boardName, osVersion, address, user)

		bars := make(map[int] progress.Model)

		for i := 0; i < cpuCoreCount; i++ {
			bars[i] = progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
		}

		p := tea.NewProgram(model{
			deviceInfo: deviceInfo,
			client: client,
			sub: make(chan dataMessage),
			data: nil,
			cpu: cpuData{
				count: cpuCoreCount,
				bar: bars,
				data: make(map[int] float64),
			},
		}, tea.WithAltScreen())

		if err := p.Start(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Error fetching data...")
	}
}