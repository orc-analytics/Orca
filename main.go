package main

import (
	"sync"

	api "github.com/predixus/analytics_framework/internal/api"
	grpc "github.com/predixus/analytics_framework/internal/grpc"
)

func main() {
	// start a waitgroup
	var wg sync.WaitGroup

	// Start the gRPC framework server
	wg.Add(1)
	go func() {
		grpc.StartGRPCServer(&wg)
	}()

	// Start the API - TODO: kick out into seperate module
	wg.Add(1)
	go func() {
		api.StartHTTPServer(&wg)
	}()

	wg.Wait()
}
