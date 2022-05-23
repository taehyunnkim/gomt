package cmd

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"strconv"

	"github.com/taehyunnkim/gomt/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/go-routeros/routeros"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command{
	Use:     "gomt [IP Address]",
	Short:   "Go MikroTik is a console monitor application for MikroTik devices",
	Version: "0.1.2",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Please enter the MikroTik device's IP")
		}

		address := args[0]

		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatal(err)
		}

		var user string
		fmt.Print("Enter the user: ")
		fmt.Scanf("%s", &user)

		fmt.Print("Enter the password: ")
		password, _ := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()

		client, err := routeros.Dial(address + ":" + port, user, string(password))

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

			m := tui.New(client, deviceInfo, cpuCoreCount)
			p := tea.NewProgram(m, tea.WithAltScreen())

			if err := p.Start(); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Error fetching data...")
		}
	},
}

// Execute runs the root command and starts the application.
func Execute() {	
	var port string

	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "8728", "Path to write to file on open.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}