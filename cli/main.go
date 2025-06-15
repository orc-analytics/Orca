package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"
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

	// Parse the appropriate subcommand
	switch os.Args[1] {

	case "start":
		checkDockerInstalled()

		startCmd.Parse(os.Args[2:])

		fmt.Println()
		networkName := createNetworkIfNotExists()
		fmt.Println()

		startPostgres(networkName)
		fmt.Println()

		startRedis(networkName)
		fmt.Println()

		// check for postgres instance running first
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		err := waitForPgReady(ctx, pgContainerName, time.Millisecond*500)
		if err != nil {
			fmt.Println(
				renderError(
					fmt.Sprintf("Issue waiting for Postgres store to start: %v", err.Error()),
				),
			)
			os.Exit(1)
		}
		startOrca(networkName)
		fmt.Println()

		fmt.Println(renderSuccess("✅ Orca stack started successfully."))
		fmt.Println()

	case "stop":
		checkDockerInstalled()

		stopCmd.Parse(os.Args[2:])

		fmt.Println()
		stopContainers()

		fmt.Println()
		fmt.Println(renderSuccess("✅ All containers stopped."))
		fmt.Println()

	case "status":
		checkDockerInstalled()
		statusCmd.Parse(os.Args[2:])

		fmt.Println()
		showStatus()
		fmt.Println()

	case "destroy":
		checkDockerInstalled()
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
		fmt.Println(renderError(fmt.Sprintf("Unknown subcommand: %s", os.Args[1])))
		fmt.Println(renderInfo("Run 'help' for usage information."))
		fmt.Println()
		os.Exit(1)
	}
}
