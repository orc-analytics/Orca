package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// starts the stack
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	// stops it
	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)
	// prints the status
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	// tears down all the images & data
	destroyCmd := flag.NewFlagSet("destroy", flag.ExitOnError)
	// shows help
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

	// check if a subcommand is provided
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	// parse the appropriate subcommand
	switch os.Args[1] {
	case "start":
		checkDockerInstalled()
		startCmd.Parse(os.Args[2:])
		// start the stack
		networkName := createNetworkIfNotExists()
		startPostgres(networkName)
		startRedis(networkName)
	case "stop":
		checkDockerInstalled()
		stopCmd.Parse(os.Args[2:])
		stopContainers()
	case "status":
		checkDockerInstalled()
		statusCmd.Parse(os.Args[2:])
		showStatus()
	case "destroy":
		checkDockerInstalled()
		destroyCmd.Parse(os.Args[2:])
		destroy()
	case "help":
		helpCmd.Parse(os.Args[2:])
		if helpCmd.NArg() > 0 {
			showCommandHelp(os.Args[2])
		} else {
			showHelp()
		}
	default:
		fmt.Printf("Unknown subcommand: %s\n", os.Args[1])
		fmt.Println("Run 'help' for usage information")
		os.Exit(1)
	}
}
