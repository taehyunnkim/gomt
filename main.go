package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"strconv"

	"github.com/taehyunnkim/gomt/internal/tui"
	"github.com/taehyunnkim/gomt/internal/net"
	"github.com/taehyunnkim/gomt/internal/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var name string = `
  ___     __  __ _ _          _____ _ _   
 / __|___|  \/  (_) |___ _ __|_   _(_) |__
| (_ / _ \ |\/| | | / / '_/ _ \| | | | / /
 \___\___/_|  |_|_|_\_\_| \___/|_| |_|_\_\`

var version string = "0.1.3"

var rootCmd = &cobra.Command{
	Use:     "gomt [IP ADDRESS] [FLAGS]",
	Short:   "Go MikroTik is a console monitor application for MikroTik devices",
	Version: version,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
		fmt.Printf("%s v%s\n\n", name, version)
		
		var address string

		if len(args) < 1 {
			fmt.Print("Enter the IP Address: ")
			fmt.Scanf("%s ", &address)
		} else {
			address = args[0]
		}

		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatal(err)
		}

		rwc, err := net.Dial(address + ":" + port, config.Timeout)

		if err != nil {
			log.Fatal(err)	
		}

		var user string
		fmt.Print("Enter the username: ")
		fmt.Scanf("%s ", &user)

		fmt.Print("Enter the password: ")
		password, _ := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()

		client, err := net.NewClientAndLogin(rwc, user, string(password))

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
func main() {	
	var port string

	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "8728", "Path to write to file on open.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}