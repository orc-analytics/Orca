package main

import (
	"fmt"
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
			"-p", "0:5432",
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
			"-p", "0:6379",
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

func startOrca(networkName string) {
	exists := checkStartContainer(orcaContainerName)
	if !exists {
		args := []string{
			"run",
			"-d",
			"--name",
			orcaContainerName,
			"--network",
			networkName,
			"-p", "0:3335",
			"-e", fmt.Sprintf("ORCA_CONNECTION_STRING=postgresql://orca:orca@%s:5432/orca?sslmode=disable", pgContainerName),
			"-e", "ORCA_PORT=3335",
			"-e", "ORCA_LOG_LEVEL=DEBUG",
			"ghcr.io/orc-analytics/orca:latest",
			"-migrate",
		}
		runCmd := exec.Command("docker", args...)
		streamCommandOutput(runCmd, "Orca-Core:")
	}
}
