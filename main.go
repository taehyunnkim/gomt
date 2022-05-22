package main

import (
	"fmt"
	"log"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
	"golang.org/x/term"
)

type model struct {
	client *routeros.Client
	sub chan fetchMessage
	clientMessage *routeros.Reply
	err error
}

type errMsg struct{ error }

type fetchMessage struct {
	reply *routeros.Reply
	err error
}

func fetchResourceInfo(c *routeros.Client, sub chan fetchMessage) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Second)
			reply, err := c.RunArgs([]string{"/system/resource/print"})
			sub <- fetchMessage{reply, err}
		}
	}
}

func waitForMessage(sub chan fetchMessage) tea.Cmd {
	return func() tea.Msg {
		return fetchMessage(<-sub)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		fetchResourceInfo(m.client, m.sub),
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
	case fetchMessage:
		m.clientMessage = message.reply
		m.err = message.err
		return m, waitForMessage(m.sub)
	case errMsg:
		m.err = message
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf(
			"A problem has occured :(\n" +
			"%s\n",
			m.err,	
		)
	}

	var result string

	for _, message := range m.clientMessage.Re {
		for _, pair := range message.List {
			result += fmt.Sprintf("%s: %s\n", pair.Key, pair.Value)
		}
	}

	return fmt.Sprintf(
		"The program is running\n" +
		"%s\n" +
		"Press q to exit...",
		result,	
	)

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
		client,
		make(chan fetchMessage),
		&routeros.Reply{
			Re: []*proto.Sentence{proto.NewSentence()},
		},
		nil,
	})

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}