package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-routeros/routeros"
	"golang.org/x/term"
	"syscall"
	"log"
	"fmt"
)

type model struct {
	address string
	user string
	password string
	clientMessage string
	err error
}

type clientMessage string

type errMsg struct{ error }

func dial(m model) (*routeros.Client, error) {
	return routeros.Dial(m.address, m.user, m.password)	
}

func getClientMessage(m model) tea.Msg {
	client, err := dial(m)

	if err!= nil {
		return errMsg{err}
	}

	defer client.Close()

	r, err := client.RunArgs([]string{"/interface/print"})

	if err != nil {
		return errMsg{err}
	}

	return clientMessage(r.String())
}

func (m model) Init() tea.Cmd {
	startClient := func() tea.Msg {
		return getClientMessage(m)
	} 

	return startClient
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
	case clientMessage:
		m.clientMessage = string(message)
		return m, nil
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

	return fmt.Sprintf(
		"The program is running(\n" +
		"%s\n" +
		"Press q to exit...",
		m.clientMessage,	
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

	p := tea.NewProgram(model{
		address,
		user,
		string(password),
		"Fetching data...",
		nil,
	})

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}