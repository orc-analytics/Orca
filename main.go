package main

import (
	api "github.com/predixus/analytics_framework/internal/api"
	cli "github.com/predixus/analytics_framework/internal/cli"
	dlyr "github.com/predixus/analytics_framework/internal/datalayer"
	grpc "github.com/predixus/analytics_framework/internal/grpc"
	li "github.com/predixus/analytics_framework/internal/logger"
	provision "github.com/predixus/analytics_framework/internal/provision_store"
	setup "github.com/predixus/analytics_framework/internal/setup"
)

func main() {
	CliArgs := cli.ParseInputs()

	if CliArgs.InitialiseDB {
		println("Initialising local postgres DB")
		err := provision.Provision()
		if err != nil {
			li.Logger.Fatal(err)
		}
		if !CliArgs.Continue {
			return
		}
	}

	db := &dlyr.Db{}
	grpc_server := &grpc.GrpcServer{}
	api_server := &api.HttpServer{}
	setup.Setup(db, grpc_server, api_server)
}
