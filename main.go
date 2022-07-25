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
	"github.com/daviddengcn/go-colortext"
)

var name string = `
  ___     __  __ _ _          _____ _ _   
 / __|___|  \/  (_) |___ _ __|_   _(_) |__
| (_ / _ \ |\/| | | / / '_/ _ \| | | | / /
 \___\___/_|  |_|_|_\_\_| \___/|_| |_|_\_\`

var version string = "0.1.5"

var rootCmd = &cobra.Command{
	Use:     "gomt [IP ADDRESS] [FLAGS]",
	Short:   "Go MikroTik is a console monitor application for MikroTik devices",
	Version: version,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
		fmt.Printf("%s v%s\n\n", name, version)

		useEnv, err := cmd.Flags().GetBool("env")

		if err != nil {
			log.Fatal(err)			
		}
	
		debug, err := cmd.Flags().GetBool("debug")

		if err != nil {
			log.Fatal(err)			
		}

		var address string
		if useEnv {
			address = os.Getenv("GOMT_IP")
		}
		if address == "" {
			if len(args) > 0 {
				address = args[0]
			} else {
				fmt.Print("Enter the IP Address: ")
				fmt.Scanf("%s ", &address)
			}
		}

		port, err := cmd.Flags().GetString("port")

		if useEnv {
			p := os.Getenv("GOMT_PORT")

			if p != "" {
				port = p
			}
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("Connecting to the device...")
		rwc, err := net.Dial(address + ":" + port, config.Timeout)

		if err != nil {
			ct.Foreground(ct.Red, true)
			log.Fatal(err)	
			ct.ResetColor()
		}

		routerosClient, err := net.NewRouterOsClient(rwc)

		if err != nil {
			ct.Foreground(ct.Red, true)
			log.Fatal(err)
			ct.ResetColor()
		}

		ct.Foreground(ct.Green, true)
		fmt.Println("Connected!")
		ct.ResetColor()
		
		defer routerosClient.Close()

		var user string
		if useEnv {
			user = os.Getenv("GOMT_USER")
		}
		if user == "" {
			fmt.Print("Enter the username: ")
			fmt.Scanf("%s ", &user)
		}

		var password string
		if useEnv {
			password = os.Getenv("GOMT_PASSWORD")
		}
		if password == "" {
			fmt.Print("Enter the password: ")
			pwd, _ := term.ReadPassword(int(syscall.Stdin))
			password = string(pwd)
		}

		err = net.Login(routerosClient, user, password)
		if err!= nil {
			ct.Foreground(ct.Red, true)
			log.Fatalf("\n%s", err)
			ct.ResetColor()
			return
		}
		
		ct.Foreground(ct.Green, true)
		fmt.Printf("\nSuccessfully logged in as %s!\n", user)
		ct.ResetColor()

		reply, err := routerosClient.RunArgs([]string{"/system/resource/print"})

		if err != nil {
			log.Fatal(err)
		}

		if len(reply.Re) > 0 {
			platform := reply.Re[0].Map["platform"]
			boardName := reply.Re[0].Map["board-name"]
			osVersion := reply.Re[0].Map["version"]
			cpuCoreCount, _ := strconv.Atoi(reply.Re[0].Map["cpu-count"])

			deviceInfo := tui.DeviceInfo {
				Platform: platform,
				BoardName: boardName,
				OsVersion: osVersion, 
				CpuCoreCount: cpuCoreCount,
			}

			m := tui.New(routerosClient, deviceInfo, debug, config.MinWindowWidth)
			p := tea.NewProgram(m, tea.WithAltScreen())

			if err := p.Start(); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Error fetching data...")
		}
	},
}

func main() {	
	rootCmd.PersistentFlags().StringP(
		"port", 
		"p", 
		strconv.Itoa(config.DefaultApiPort), 
		"Path to write to file on open.",
	)

	rootCmd.PersistentFlags().BoolP(
		"env",
		"e",
		false,
		"Use environment variables.\n[GOMT_IP, GOMT_PORT, GOMT_USER, GOMT_PASSWORD]",
	)

	rootCmd.PersistentFlags().BoolP(
		"debug",
		"d",
		false,
		"Debug mode",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
