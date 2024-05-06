package setup

import (
	"sync"

	api "github.com/predixus/analytics_framework/internal/api"
	dlyr "github.com/predixus/analytics_framework/internal/datalayer"
	grpc "github.com/predixus/analytics_framework/internal/grpc"
)

func Setup(
	db_connector dlyr.DB,
	grpc_server grpc.GRPCServer,
	api_server api.HTTPServer,
) error {
	// connect to the db
	err := db_connector.Connect()
	if err != nil {
		return err
	}
	defer db_connector.Close()

	// start a waitgroup
	var wg sync.WaitGroup

	// start the gRPC framework server
	wg.Add(1)
	go func() {
		grpc_server.Start(&wg)
	}()

	// start the API server
	wg.Add(1)
	go func() {
		api_server.Start(&wg)
	}()

	wg.Wait()

	return nil
}
