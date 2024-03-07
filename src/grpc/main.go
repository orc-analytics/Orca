package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"

	infc "github.com/predixus/analytics_framework/protobufs/go"
	"github.com/predixus/analytics_framework/src/grpc/epoch"
)

// implement the GRPC server requests
func StartGRPCServer(wg *sync.WaitGroup) {
	defer wg.Done()

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "8080"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", os.Getenv("GRPC_PORT")))
	if err != nil {
		log.Fatalf("Cannnot create listener: %s", err)
	}

	// server opts
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	service := &grpc_epoch.EpochServiceServer{}

	// register all the services here
	defer wg.Done()
	infc.RegisterEpochServiceServer(grpcServer, service)

	grpcServer.Serve(lis)

	log.Println("GRPC: Shutting Down")
}
