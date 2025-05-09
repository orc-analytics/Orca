package main

import (
	"os/exec"
)

// startPostgres starts the postgres instance that orca needs.
func startPostgres(networkName string) {
	containerName := "orca-pg-instance"

	exists := checkStartContainer(containerName)

	if !exists {
		// create or start a volume
		volumeName := checkCreateVolume(containerName)

		// run container with volume mounted
		args := []string{
			"run",
			"-d",
			"--name",
			containerName,
			"--network",
			networkName,
			"-e",
			"POSTGRES_USER=orca",
			"-e",
			"POSTGRES_PASSWORD=orca",
			"-e",
			"POSTGRES_DB=orca",
			"-v",
			volumeName + ":/var/lib/postgresql/data",
			"postgres",
		}

		runCmd := exec.Command("docker", args...)
		// stream container creation logs
		streamCommandOutput(runCmd, "PostgreSQL Store:")
	}
}

func startRedis(networkName string) {
	containerName := "orca-redis-instance"
	exists := checkStartContainer(containerName)

	if !exists {
		// create or start a volume
		volumeName := checkCreateVolume(containerName)

		// run container with volume mounted
		args := []string{
			"run",
			"--name", containerName,
			"--network", networkName,
			"-d",
			"-v", volumeName + ":/data",
			"redis",
			"redis-server", "--appendonly", "yes",
		}

		runCmd := exec.Command("docker", args...)
		// stream container creation logs
		streamCommandOutput(runCmd, "Redis Cache:")
	}
}
