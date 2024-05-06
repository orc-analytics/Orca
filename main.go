package main

import (
	api "github.com/predixus/analytics_framework/internal/api"
	dlyr "github.com/predixus/analytics_framework/internal/datalayer"
	grpc "github.com/predixus/analytics_framework/internal/grpc"
	setup "github.com/predixus/analytics_framework/internal/setup"
)

func main() {
	db := &dlyr.DbConnector{}
	grpc_server := &grpc.GrpcServer{}
	api_server := &api.HttpServer{}
	setup.Setup(db, grpc_server, api_server)
}
