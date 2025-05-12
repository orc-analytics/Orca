package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// checkCreateVolume checks if a volume exists for a container and if not creates it
func checkCreateVolume(containerName string) string {
	// create a volume with a name specific to the orca storage container
	volumeName := containerName + "-data"

	// check if the volume already exists
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
		// volume doesn't exist, create it
		fmt.Printf("Creating volume %s...\n", volumeName)
		createVolumeCmd := exec.Command("docker", "volume", "create", volumeName)
		if err := createVolumeCmd.Run(); err != nil {
			fmt.Printf("Failed to create volume: %s\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Using existing volume: %s\n", volumeName)
	}
	return volumeName
}

func checkStartContainer(containerName string) bool {
	// check if container already exists
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
		// container exists, restart it
		fmt.Printf("Container %s already exists, checking run status...\n", containerName)

		// first check if it's already running
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
			// container is running
			fmt.Println(containerName + " already running")
			return true
		}

		// start the container
		startCmd := exec.Command("docker", "start", containerName)

		// stream start logs
		streamCommandOutput(startCmd, "Starting container")

		fmt.Printf("Container started successfully\n")
		return true
	}
	return false
}

// helper function to stream command output
func streamCommandOutput(cmd *exec.Cmd, prefix string) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe: %s\n", err)
		os.Exit(1)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error creating stderr pipe: %s\n", err)
		os.Exit(1)
	}

	// start the command
	if err := cmd.Start(); err != nil {
		fmt.Printf("%s failed: %s\n", prefix, err)
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
			fmt.Printf("%s (stdout): %s\n", prefix, scanner.Text())
		}
	}()

	// stream stderr
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Printf("%s (stderr): %s\n", prefix, scanner.Text())
		}
	}()

	// wait for both streams to finish
	wg.Wait()

	// wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Printf("%s command failed: %s\n", prefix, err)
		os.Exit(1)
	}
}

// createNetworkIfNotExists creates a bridge network if it doesn't already exist
func createNetworkIfNotExists() string {
	// check if network exists
	checkCmd := exec.Command(
		"docker",
		"network",
		"ls",
		"--filter", "name="+networkName,
		"--format", "{{.Name}}",
	)
	output, err := checkCmd.CombinedOutput()

	if err != nil || !strings.Contains(string(output), networkName) {
		fmt.Printf("Creating network '%s'...\n", networkName)

		// Create bridge network
		createCmd := exec.Command(
			"docker",
			"network",
			"create",
			"--driver", "bridge",
			networkName,
		)

		streamCommandOutput(createCmd, "Network creation:")
		fmt.Printf("Network '%s' created successfully\n", networkName)
	} else {
		fmt.Printf("Using existing network: %s\n", networkName)
	}
	return networkName
}

// showStatus prints the status of each container along with connection strings
func showStatus() {
	fmt.Println("\n=== Container Status ===")

	// Check PostgreSQL status
	pgStatus := getContainerStatus(pgContainerName)
	fmt.Printf("PostgreSQL: %s\n", pgStatus)
	if pgStatus == "running" {
		// Get PostgreSQL connection info
		pgIP := getContainerIP("orca-pg-instance")
		if pgIP != "" {
			fmt.Printf("Connection string: postgresql://orca:orca@%s:5432/orca\n", pgIP)
		}
	}

	// Check Redis status
	redisStatus := getContainerStatus(redisContainerName)
	fmt.Printf("\nRedis: %s\n", redisStatus)
	if redisStatus == "running" {
		// Get Redis connection info
		redisIP := getContainerIP("orca-redis-instance")
		if redisIP != "" {
			fmt.Printf("Connection string: redis://%s:6379\n", redisIP)
		}
	}

	fmt.Println("=======================")
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
	fmt.Println("\n=== Stopping Orca Containers ===")

	for _, containerName := range orcaContainers {
		status := getContainerStatus(containerName)

		if status == "running" {
			fmt.Printf("Stopping %s... ", containerName)

			cmd := exec.Command("docker", "stop", containerName)
			err := cmd.Run()

			if err != nil {
				fmt.Printf("ERROR: Failed to stop container: %v\n", err)
			} else {
				fmt.Println("STOPPED")
			}
		} else if status == "stopped" {
			fmt.Printf("%s is already stopped\n", containerName)
		} else {
			fmt.Printf("%s not found\n", containerName)
		}
	}
}

// destroy tears down all Orca-related resources (containers, images, networks, and volumes)
// It requires user confirmation before executing destructive operations
func destroy() {
	fmt.Println("\n!!! WARNING: DESTRUCTIVE OPERATION !!!")
	fmt.Println("This will remove all Orca containers, images, networks, and volumes.")
	fmt.Println("All data will be permanently lost.")
	fmt.Print("\nAre you sure you want to continue? (y/N): ")

	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" {
		fmt.Println("Operation cancelled.")
		return
	}

	// Stop all containers first
	stopContainers()

	// Remove containers
	fmt.Println("\n=== Removing Orca Containers ===")
	for _, containerName := range orcaContainers {
		fmt.Printf("Removing container %s... ", containerName)

		cmd := exec.Command("docker", "rm", "-f", containerName)
		err := cmd.Run()

		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else {
			fmt.Println("REMOVED")
		}
	}

	// Remove volumes
	fmt.Println("\n=== Removing Orca Volumes ===")
	for _, volumeName := range orcaVolumes {
		fmt.Printf("Removing volume %s... ", volumeName)

		cmd := exec.Command("docker", "volume", "rm", volumeName)
		err := cmd.Run()

		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else {
			fmt.Println("REMOVED")
		}
	}

	// Remove the Orca network
	fmt.Println("\n=== Removing Orca Network ===")
	cmd := exec.Command("docker", "network", "rm", "orca-network")
	err := cmd.Run()

	if err != nil {
		fmt.Printf("ERROR: Failed to remove network: %v\n", err)
	} else {
		fmt.Println("Network orca-network REMOVED")
	}

	// Instead of automatically removing images, provide instructions to the user
	fmt.Println("\n=== Orca Image Cleanup Instructions ===")
	fmt.Println("To clean up Docker images related to Orca, you can run these commands:")
	fmt.Println("  docker rmi postgres    # Remove PostgreSQL image")
	fmt.Println("  docker rmi redis       # Remove Redis image")
	fmt.Println("")
	fmt.Println("Or to remove all unused images:")
	fmt.Println("  docker image prune -a  # Remove all unused images")
	fmt.Println("")
	fmt.Println(
		"Note: These commands will only work if the images are not used by other containers.",
	)
	fmt.Println("\n=== Orca Environment Destroyed ===")
}

// checkDockerInstalled verifies that Docker is installed and accessible
// If Docker is not installed, it exits with an error message
func checkDockerInstalled() {
	fmt.Println("Checking for Docker installation...")

	cmd := exec.Command("docker", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("ERROR: Docker is not installed or not in PATH")
		fmt.Println("Please install Docker before continuing:")
		fmt.Println("  - For Windows/Mac: https://www.docker.com/products/docker-desktop")
		fmt.Println("  - For Linux: https://docs.docker.com/engine/install/")
		os.Exit(1)
	}

	version := strings.TrimSpace(string(output))
	fmt.Printf("Docker found: %s\n", version)

	// Check if Docker daemon is running
	cmd = exec.Command("docker", "info")
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("ERROR: Docker daemon is not running")
		fmt.Println("Please start the Docker service before continuing.")
		os.Exit(1)
	}

	fmt.Println("Docker is installed and running correctly.")
}
