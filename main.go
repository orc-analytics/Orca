package main

import (
	api "github.com/predixus/analytics_framework/src/api"
	grpc "github.com/predixus/analytics_framework/src/grpc"
)

func main() {
	// Start the gRPC framework server
	grpc.StartGRPCServer()

	// Start the API - TODO: kick out into seperate module
	api.StartHTTPServer()
}
