package setup

import (
	"sync"

	api "github.com/predixus/analytics_framework/internal/api"
	datalayer "github.com/predixus/analytics_framework/internal/datalayer"
	grpc "github.com/predixus/analytics_framework/internal/grpc"
)

func Setup() {
	// connect to the store
	datalayer.ConnectDB()
	defer datalayer.StorageDB.Close()

	// start a waitgroup
	var wg sync.WaitGroup

	// Start the gRPC framework server
	wg.Add(1)
	go func() {
		grpc.StartGRPCServer(&wg)
	}()

	// Start the API - TODO: kick out into seperate module
	// we want to separate out the services for task creation and
	// data getting.
	wg.Add(1)
	go func() {
		api.StartHTTPServer(&wg)
	}()

	wg.Wait()
}
