package main

const (
	pgContainerName    = "orca-pg-instance"
	redisContainerName = "orca-redis-instance"
	orcaContainerName  = "orca-instance"
	networkName        = "orca-network"
)

var orcaContainers = []string{
	pgContainerName,
	redisContainerName,
	orcaContainerName,
}

// follows pattern of <container-name>-data
var orcaVolumes = []string{
	"orca-pg-instance-data",
	"orca-redis-instance-data",
}
