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
