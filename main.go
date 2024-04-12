package main

import (
	"flag"
	"sync"

	api "github.com/predixus/analytics_framework/internal/api"
	datalayer "github.com/predixus/analytics_framework/internal/datalayer"
	provision "github.com/predixus/analytics_framework/internal/datalayer/provision"
	grpc "github.com/predixus/analytics_framework/internal/grpc"
	li "github.com/predixus/analytics_framework/internal/logger"
)

func main() {
	// first check whether there are any command line arguments
	initPtr := flag.Bool("init-db", false, "Provision the local postgres db")
	flag.Parsed()

	if *initPtr {
		println("Initialising postgres DB")
		err := provision.Provision()
		if err != nil {
			li.Logger.Fatal(err)
		}
		return
	}

	// connect to the postgres store
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
	wg.Add(1)
	go func() {
		api.StartHTTPServer(&wg)
	}()

	wg.Wait()
}
