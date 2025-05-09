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
	fmt.Println("  help    - Show this help message or help for a specific command")
	fmt.Println("\nGlobal Options:")
	fmt.Println("  -image string    Docker image to use (default: nginx:latest)")
	fmt.Println("  -name string     Name for the container")
	fmt.Println(
		"  -port string     Port mapping (e.g., '8080:80' or multiple with '8080:80,9000:9000')",
	)
	fmt.Println("  -env string      Environment variables (e.g., 'KEY1=VAL1,KEY2=VAL2')")
	fmt.Println("  -volume string   Volume mapping (e.g., '/host/path:/container/path')")
	fmt.Println("\nExamples:")
	fmt.Println("  start -image redis -name cache-server -port 6379:6379")
	fmt.Println("  stop -name cache-server")
	fmt.Println("  status")
	fmt.Println("  status -name cache-server")
	fmt.Println("  help start")
}

func showCommandHelp(command string) {
	switch command {
	case "start":
		fmt.Println("'start' command - Start a Docker container")
		fmt.Println("\nUsage:")
		fmt.Println("  start [options]")
		fmt.Println("\nOptions:")
		fmt.Println("  -image string    Docker image to use (default: nginx:latest)")
		fmt.Println("  -name string     Name for the container")
		fmt.Println(
			"  -port string     Port mapping (e.g., '8080:80' or multiple with '8080:80,9000:9000')",
		)
		fmt.Println("  -env string      Environment variables (e.g., 'KEY1=VAL1,KEY2=VAL2')")
		fmt.Println("  -volume string   Volume mapping (e.g., '/host/path:/container/path')")
		fmt.Println("\nExample:")
		fmt.Println("  start -image redis -name cache-server -port 6379:6379")
	case "stop":
		fmt.Println("'stop' command - Stop a running Docker container")
		fmt.Println("\nUsage:")
		fmt.Println("  stop [options]")
		fmt.Println("\nOptions:")
		fmt.Println("  -name string     Name of the container to stop (required)")
		fmt.Println("\nExample:")
		fmt.Println("  stop -name cache-server")
	case "status":
		fmt.Println("'status' command - Show status of Docker containers")
		fmt.Println("\nUsage:")
		fmt.Println("  status [options]")
		fmt.Println("\nOptions:")
		fmt.Println("  -name string     Name of the container to check (optional)")
		fmt.Println("\nExamples:")
		fmt.Println("  status           # Shows all containers")
		fmt.Println("  status -name cache-server")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
	}
}
