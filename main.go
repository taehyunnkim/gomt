package main

import (
	"fmt"
	"log"
	"syscall"
	"time"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-routeros/routeros"
	"golang.org/x/term"
)

type model struct {
	client *routeros.Client
	sub chan dataMessage
	data *dataMessage
	cpuData map[int] string
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
		var cpuCount int

		if m.data.resourceData.err != nil {
			result = fmt.Sprintf(
				"A problem has occured :(\n" +
				"%s\n",
				m.data.resourceData.err,	
			)
		} else {
			for _, message := range m.data.resourceData.reply.Re {
				for _, pair := range message.List {
					result += fmt.Sprintf("%s: %s\n", pair.Key, pair.Value)

					if pair.Key == "cpu-count" {
						cpuCount, _ = strconv.Atoi(pair.Value)
					}
				}
			}
		}

		result += "\n===== CPU =====\n"

		if m.data.cpuData.err != nil {
			result = fmt.Sprintf(
				"A problem has occured :(\n" +
				"%s\n",
				m.data.cpuData.err,	
			)
		} else {
			parseCpuData(m.data.cpuData, &m.cpuData)
			for i := 0; i < cpuCount; i++ {
				result += fmt.Sprintf("cpu-%d: %s%%\n", i+1, m.cpuData[i])	
			}
		}
	}

	return fmt.Sprintf(
		"The program is running\n\n" +
		"%s\n" +
		"Press q to exit...",
		result,	
	)
}

func parseCpuData(data data, m *map[int] string) {
	for _, message := range data.reply.Re {
		var cpu int

		for _, pair := range message.List {
			if pair.Key == ".id" {
				_, err := fmt.Sscanf(pair.Value, "*%d", &cpu)

				if err != nil {
					log.Println(pair.Value)
					log.Fatal(err)	
				}
			} else if pair.Key == "load" {
				(*m)[cpu] = pair.Value
			}
		}
	}
}

func main() {
	fmt.Println("ncmt v0.1.0")

	var address string
	fmt.Print("Enter the full address: ")
	fmt.Scanf("%s", &address)

	var user string
	fmt.Print("Enter the user: ")
	fmt.Scanf("%s", &user)

	fmt.Print("Enter the password: ")
	password, _ := term.ReadPassword(int(syscall.Stdin))

	client, err := routeros.Dial(address, user, string(password))

	if err!= nil {
		log.Fatal(err)
		return
	}

	defer client.Close()

	p := tea.NewProgram(model{
		client: client,
		sub: make(chan dataMessage),
		data: nil,
		cpuData: make(map[int]string),
	})

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}