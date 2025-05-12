package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Define subcommands
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	destroyCmd := flag.NewFlagSet("destroy", flag.ExitOnError)
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

	// Check if a subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println()
		showHelp()
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println()

	// Parse the appropriate subcommand
	switch os.Args[1] {

	case "start":
		checkDockerInstalled()

		fmt.Println()
		startCmd.Parse(os.Args[2:])

		fmt.Println(headerStyle.Render("Starting Orca stack..."))
		fmt.Println()

		networkName := createNetworkIfNotExists()
		fmt.Println()

		startPostgres(networkName)
		fmt.Println()

		startRedis(networkName)
		fmt.Println()

		fmt.Println(successStyle.Render("✅ Orca stack started successfully."))
		fmt.Println()

	case "stop":
		checkDockerInstalled()

		fmt.Println()
		stopCmd.Parse(os.Args[2:])

		fmt.Println()
		stopContainers()

		fmt.Println()
		fmt.Println(successStyle.Render("✅ All containers stopped."))
		fmt.Println()

	case "status":
		checkDockerInstalled()
		statusCmd.Parse(os.Args[2:])

		fmt.Println()
		showStatus()
		fmt.Println()

	case "destroy":
		checkDockerInstalled()
		fmt.Println()
		destroyCmd.Parse(os.Args[2:])
		fmt.Println()
		destroy()
		fmt.Println()

	case "help":
		helpCmd.Parse(os.Args[2:])
		fmt.Println()
		if helpCmd.NArg() > 0 {
			showCommandHelp(os.Args[2])
		} else {
			showHelp()
		}
		fmt.Println()

	default:
		fmt.Println()
		fmt.Println(errorStyle.Render(fmt.Sprintf("Unknown subcommand: %s", os.Args[1])))
		fmt.Println(infoStyle.Render("Run 'help' for usage information."))
		fmt.Println()
		os.Exit(1)
	}
}
