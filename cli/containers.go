package main

import (
	"os/exec"
)

// startPostgres starts the postgres instance that orca needs.
func startPostgres(networkName string) {
	exists := checkStartContainer(pgContainerName)

	if !exists {
		// create or start a volume
		volumeName := checkCreateVolume(pgContainerName)

		// run container with volume mounted
		args := []string{
			"run",
			"-d",
			"--name",
			pgContainerName,
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
	exists := checkStartContainer(redisContainerName)

	if !exists {
		// create or start a volume
		volumeName := checkCreateVolume(redisContainerName)

		// run container with volume mounted
		args := []string{
			"run",
			"--name", redisContainerName,
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
