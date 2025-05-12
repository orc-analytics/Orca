package main

import "fmt"

func showHelp() {
	fmt.Println("Orca CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  command [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  start   - Start the Orca stack")
	fmt.Println("  stop    - Stop the Orca stack")
	fmt.Println("  status  - Show status of the Orca components")
	fmt.Println("  destroy - Delete all Orca resources")
	fmt.Println("  help    - Show this help message or help for a specific command")

	fmt.Println("\nExamples:")
	fmt.Println("  start")
	fmt.Println("  stop")
	fmt.Println("  status")
	fmt.Println("  destroy")
	fmt.Println("  help start")
}

func showCommandHelp(command string) {
	switch command {
	case "start":
		fmt.Println("'start' command - Start the Orca stack")
		fmt.Println("\nUsage:")
		fmt.Println("  start")
		fmt.Println("\nExample:")
		fmt.Println("  start")
	case "stop":
		fmt.Println("'stop' command - Stop the Orca stack")
		fmt.Println("\nUsage:")
		fmt.Println("  stop")
		fmt.Println("\nExample:")
		fmt.Println("  stop")
	case "status":
		fmt.Println("'status' command - Show status of the Orca Services")
		fmt.Println("\nUsage:")
		fmt.Println("  status")
		fmt.Println("\nExamples:")
		fmt.Println("  status")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
	}
}
