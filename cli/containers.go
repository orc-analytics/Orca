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
			"ghcr.io/predixus/orca:latest",
			"-connStr", fmt.Sprintf("postgresql://orca:orca@%v:5432/orca?sslmode=disable", pgContainerName),
			"-migrate",
			"-platform", "postgresql",
			"-port", "3335",
		}

		runCmd := exec.Command("docker", args...)
		streamCommandOutput(runCmd, "Orca-Core:")
	}
}
