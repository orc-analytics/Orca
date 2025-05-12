package main

const (
	pgContainerName = "orca-pg-instance"

	redisContainerName = "orca-redis-instance"
	networkName        = "orca-network"
)

var orcaContainers = []string{
	pgContainerName,
	redisContainerName,
}

// follows pattern of <container-name>-data
var orcaVolumes = []string{
	"orca-pg-instance-data",
	"orca-redis-instance-data",
}
