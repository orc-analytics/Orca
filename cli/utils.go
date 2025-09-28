package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// checkCreateVolume checks if a volume exists for a container and if not creates it
func checkCreateVolume(containerName string) string {
	// Create a volume with a name specific to the orca storage container
	volumeName := containerName + "-data"

	// Check if the volume already exists
	volumeCheckCmd := exec.Command(
		"docker",
		"volume",
		"ls",
		"--filter",
		"name="+volumeName,
		"--format",
		"{{.Name}}",
	)
	volumeOutput, volumeErr := volumeCheckCmd.CombinedOutput()

	if volumeErr != nil || !strings.Contains(string(volumeOutput), volumeName) {
		fmt.Println(infoStyle.Render(fmt.Sprintf("Creating volume %s...", volumeName)))

		createVolumeCmd := exec.Command("docker", "volume", "create", volumeName)
		if err := createVolumeCmd.Run(); err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("Failed to create volume: %s", err)))
			os.Exit(1)
		}
		fmt.Println(successStyle.Render(fmt.Sprintf("Volume %s created successfully", volumeName)))
	} else {
		fmt.Println(infoStyle.Render(fmt.Sprintf("Using existing volume: %s", volumeName)))
	}

	return volumeName
}

func checkPostgresReady(ctx context.Context, containerName string) (bool, error) {
	// Command to run pg_isready inside the container
	healthCmd := exec.CommandContext(
		ctx,
		"docker",
		"exec",
		containerName,
		"pg_isready",
		"-U", "postgres", // Specify the default postgres user
	)

	// Run the command
	_, err := healthCmd.CombinedOutput()
	// pg_isready returns:
	// 0 - the server is accepting connections
	// 1 - the server is not accepting connections
	// 2 - the server is starting up
	// 3 - the server is not responding
	if err != nil {
		// Check if it's an exit error to interpret the return code
		if exitErr, ok := err.(*exec.ExitError); ok {
			// pg_isready uses exit codes to indicate status
			switch exitErr.ExitCode() {
			case 0:
				return true, nil // Server is ready
			case 1, 2, 3:
				return false, nil // Server not ready
			default:
				return false, fmt.Errorf("unexpected pg_isready error: %w", err)
			}
		}
		return false, fmt.Errorf("error running pg_isready: %w", err)
	}

	// If no error and command succeeded, the server is ready
	return true, nil
}

func waitForPgReady(
	ctx context.Context,
	containerName string,
	checkInterval time.Duration,
) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for container %s to be ready", containerName)
		default:
			// Use the postgres-specific ready check
			healthy, err := checkPostgresReady(ctx, containerName)
			if err != nil {
				// Log the error but continue trying
				fmt.Printf("Error checking container health: %v\n", err)
			} else if healthy {
				return nil // Container is ready
			}

			// Wait before next check
			select {
			case <-ctx.Done():
				return fmt.Errorf("timeout waiting for container %s to be ready", containerName)
			case <-time.After(checkInterval):
				// Continue to next iteration
			}
		}
	}
}

func checkStartContainer(containerName string) bool {
	// Check if container already exists
	checkCmd := exec.Command(
		"docker",
		"ps",
		"-a",
		"--filter",
		"name="+containerName,
		"--format",
		"{{.Names}}",
	)
	output, err := checkCmd.CombinedOutput()

	if err == nil && strings.Contains(string(output), containerName) {
		// Check if it's already running
		statusCmd := exec.Command(
			"docker",
			"ps",
			"--filter",
			"name="+containerName,
			"--format",
			"{{.Names}}",
		)
		statusOutput, statusErr := statusCmd.CombinedOutput()

		if statusErr == nil && strings.Contains(string(statusOutput), containerName) {
			fmt.Println(successStyle.Render(fmt.Sprintf("%s already running", containerName)))
			return true
		}

		// Start the container
		startCmd := exec.Command("docker", "start", containerName)
		streamCommandOutput(startCmd, "Starting container")

		fmt.Println(successStyle.Render("Container started successfully"))
		return true
	}

	return false
}

// helper function to stream command output
func streamCommandOutput(cmd *exec.Cmd, prefix string) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error creating stdout pipe: %s", err)))
		os.Exit(1)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error creating stderr pipe: %s", err)))
		os.Exit(1)
	}

	// start the command
	if err := cmd.Start(); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("%s failed: %s", prefix, err)))
		os.Exit(1)
	}

	// create a WaitGroup to wait for both goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// stream stdout
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(prefixStyle.Render(prefix) + " " + infoStyle.Render(scanner.Text()))
		}
	}()

	// stream stderr
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(prefixStyle.Render(prefix) + " " + warningStyle.Render(scanner.Text()))
		}
	}()

	// wait for both streams to finish
	wg.Wait()

	// wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("%s command failed: %s", prefix, err)))
		os.Exit(1)
	}
}

// createNetworkIfNotExists creates a bridge network if it doesn't already exist
func createNetworkIfNotExists() string {
	// Check if network exists
	checkCmd := exec.Command(
		"docker",
		"network",
		"ls",
		"--filter", "name="+networkName,
		"--format", "{{.Name}}",
	)
	output, err := checkCmd.CombinedOutput()

	if err != nil || !strings.Contains(string(output), networkName) {
		fmt.Println(infoStyle.Render(fmt.Sprintf("Creating network '%s'...", networkName)))

		// Create bridge network
		createCmd := exec.Command(
			"docker",
			"network",
			"create",
			"--driver", "bridge",
			networkName,
		)

		streamCommandOutput(createCmd, "Network creation:")
		fmt.Println(
			successStyle.Render(fmt.Sprintf("Network '%s' created successfully", networkName)),
		)
	} else {
		fmt.Println(infoStyle.Render(fmt.Sprintf("Using existing network: %s", networkName)))
	}

	return networkName
}

// showStatus prints the status of each container along with connection strings
func showStatus() {
	// PostgreSQL status
	pgStatus := getContainerStatus(pgContainerName)
	fmt.Println(subHeaderStyle.Render("PostgreSQL:"), statusColor(pgStatus).Render(pgStatus))

	if pgStatus == "running" {
		pgPort := getContainerPort(pgContainerName, 5432)
		conn := fmt.Sprintf("postgresql://orca:orca@localhost:%s/orca?sslmode=disable", pgPort)
		fmt.Println(infoStyle.Render("Connection string: " + conn))
	}

	fmt.Println()

	// Redis status
	redisStatus := getContainerStatus(redisContainerName)
	fmt.Println(subHeaderStyle.Render("Redis:"), statusColor(redisStatus).Render(redisStatus))

	if redisStatus == "running" {
		redisPort := getContainerPort(redisContainerName, 6379)
		conn := fmt.Sprintf("redis://localhost:%s", redisPort)
		fmt.Println(infoStyle.Render("Connection string: " + conn))
	}

	fmt.Println()

	// Orca status
	orcaStatus := getContainerStatus(orcaContainerName)
	fmt.Println(subHeaderStyle.Render("Orca:"), statusColor(orcaStatus).Render(orcaStatus))

	if orcaStatus == "running" {
		orcaPort := getContainerPort(orcaContainerName, 3335)
		conn := fmt.Sprintf("localhost:%s", orcaPort)
		fmt.Println(infoStyle.Render("Connection string: " + conn))
		fmt.Println()
		fmt.Println(
			prefixStyle.Render(
				"Set these environment variables in your Orca processors to connect to Orca:",
			),
		)
		fmt.Println(prefixStyle.Render("\tORCA_CORE=" + conn))
		fmt.Println(prefixStyle.Render("\tPROCESSOR_ADDRESS=host.docker.internal:<your-processor-port>"))
		fmt.Println()
		fmt.Println(prefixStyle.Render("Optional: Override the port Orca uses to contact your processor:"))
		fmt.Println(prefixStyle.Render("\tPROCESSOR_EXTERNAL_PORT=<custom-external-port>"))
	}
}

// getContainerStatus returns the status of a container (running, stopped, or not found)
func getContainerStatus(containerName string) string {
	cmd := exec.Command(
		"docker",
		"ps",
		"-a",
		"--filter",
		"name="+containerName,
		"--format",
		"{{.Status}}",
	)
	output, err := cmd.CombinedOutput()

	if err != nil || len(output) == 0 {
		return "not found"
	}

	status := strings.TrimSpace(string(output))
	if strings.HasPrefix(status, "Up") {
		return "running"
	} else if len(status) > 0 {
		return "stopped"
	}

	return "not found"
}

// getContainerPort retrieves the mapped port for a specific container and internal port
func getContainerPort(containerName string, internalPort int) string {
	cmd := exec.Command("docker", "port", containerName)
	output, err := cmd.Output()
	if err != nil {
		return strconv.Itoa(internalPort) // fallback to default if command fails
	}

	// convert output to string and split lines
	portInfo := string(output)
	lines := strings.Split(portInfo, "\n")

	// find the line with the internal port
	portStr := fmt.Sprintf("%d/tcp", internalPort)
	for _, line := range lines {
		if strings.Contains(line, portStr) {
			// extract the mapped port (after ->)
			parts := strings.Split(line, "->")
			if len(parts) > 1 {
				// trim whitespace and get the mapped port
				mappedPortParts := strings.Fields(parts[1])
				if len(mappedPortParts) > 0 {
					// remove any host information (like 0.0.0.0: or [::]:)
					mappedPort := strings.TrimPrefix(mappedPortParts[0], "0.0.0.0:")
					mappedPort = strings.TrimPrefix(mappedPort, "[::]:")
					return mappedPort
				}
			}
		}
	}

	// fallback to default internal port if no mapping found
	return strconv.Itoa(internalPort)
}

// getContainerIP returns the IP address of a container
func getContainerIP(containerName string) string {
	cmd := exec.Command(
		"docker",
		"inspect",
		"--format",
		"{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}",
		containerName,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

// stopContainers stops all running containers related to Orca
func stopContainers() {
	for _, containerName := range orcaContainers {
		status := getContainerStatus(containerName)

		switch status {
		case "running":
			fmt.Printf("%s Stopping %s... ", prefixStyle, containerName)

			cmd := exec.Command("docker", "stop", containerName)
			err := cmd.Run()

			if err != nil {
				fmt.Println(
					errorStyle.Render(fmt.Sprintf("ERROR: Failed to stop container: %v", err)),
				)
			} else {
				fmt.Println(successStyle.Render("STOPPED"))
			}

		case "stopped":
			fmt.Println(infoStyle.Render(fmt.Sprintf("%s is already stopped", containerName)))

		default:
			fmt.Println(warningStyle.Render(fmt.Sprintf("%s not found", containerName)))
		}
	}
}

// destroy tears down all Orca-related resources (containers, images, networks, and volumes)
// It requires user confirmation before executing destructive operations
func destroy() {
	fmt.Println(warningStyle.Render("\n!!! WARNING: DESTRUCTIVE OPERATION !!!"))
	fmt.Println(
		warningStyle.Render("This will remove all Orca containers, images, networks, and volumes."),
	)
	fmt.Println(errorStyle.Render("All data will be permanently lost."))
	fmt.Print(warningStyle.Render("\nAre you sure you want to continue? (y/N): "))

	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" {
		fmt.Println(infoStyle.Render("Operation cancelled."))
		return
	}

	// Stop all containers first
	stopContainers()

	// Remove containers
	for _, containerName := range orcaContainers {
		fmt.Printf("%s Removing container %s... ", prefixStyle, containerName)

		cmd := exec.Command("docker", "rm", "-f", containerName)
		err := cmd.Run()

		if err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("ERROR: %v", err)))
		} else {
			fmt.Println(successStyle.Render("REMOVED"))
		}
	}

	// Remove volumes
	for _, volumeName := range orcaVolumes {
		fmt.Printf("%s Removing volume %s... ", prefixStyle, volumeName)

		cmd := exec.Command("docker", "volume", "rm", volumeName)
		err := cmd.Run()

		if err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("ERROR: %v", err)))
		} else {
			fmt.Println(successStyle.Render("REMOVED"))
		}
	}

	// Remove the Orca network
	cmd := exec.Command("docker", "network", "rm", "orca-network")
	err := cmd.Run()

	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("ERROR: Failed to remove network: %v", err)))
	} else {
		fmt.Println(successStyle.Render("Network orca-network REMOVED"))
	}

	// Instead of automatically removing images, provide instructions to the user
	fmt.Println(
		infoStyle.Render("To clean up Docker images related to Orca, you can run these commands:"),
	)
	fmt.Println(infoStyle.Render("  docker rmi postgres               # Remove PostgreSQL image"))
	fmt.Println(infoStyle.Render("  docker rmi redis                  # Remove Redis image"))
	fmt.Println(infoStyle.Render("  docker rmi ghcr.io/orc-analytics/orca  # Remove Orca image"))
	fmt.Println()
	fmt.Println(infoStyle.Render("Or to remove all unused images:"))
	fmt.Println(infoStyle.Render("  docker image prune -a  # Remove all unused images"))
	fmt.Println()
	fmt.Println(
		infoStyle.Render(
			"Note: These commands will only work if the images are not used by other containers.",
		),
	)
	fmt.Println(successStyle.Render("\nOrca Environment Destroyed"))
}

// checkDockerInstalled verifies that Docker is installed and accessible
// If Docker is not installed, it exits with an error message
func checkDockerInstalled() {
	cmd := exec.Command("docker", "--version")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(errorStyle.Render("ERROR: Docker is not installed or not in PATH"))
		fmt.Println(infoStyle.Render("Please install Docker before continuing:"))
		fmt.Println(
			infoStyle.Render("  - For Windows/Mac: https://www.docker.com/products/docker-desktop"),
		)
		fmt.Println(infoStyle.Render("  - For Linux: https://docs.docker.com/engine/install/"))
		os.Exit(1)
	}

	// check if Docker daemon is running
	cmd = exec.Command("docker", "info")
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(errorStyle.Render("ERROR: Docker daemon is not running"))
		fmt.Println(infoStyle.Render("Please start the Docker service before continuing."))
		os.Exit(1)
	}
}

func getNetworkGatewayIP() string {
	cmd := exec.Command(
		"docker",
		"network",
		"inspect",
		networkName,
		"--format='{{range .IPAM.Config}}{{.Gateway}}{{end}}'",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(errorStyle.Render("ERROR: Issue grabbing the gateway IP of the orca network"))
		return ""
	}
	return string(output)
}
