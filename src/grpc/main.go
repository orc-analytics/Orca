package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	infc "github.com/predixus/analytics_framework/protobufs/go"
)

// struct level api for the GRPC service
type epochServiceserver struct {
	infc.UnimplementedEpochServiceServer
}

// implement the GRPC server requests
func (s *epochServiceserver) RegisterEpoch(
	ctx context.Context,
	epochRequest *infc.EpochRequest,
) (*infc.EpochResponse, error) {
	return nil, nil
}

func StartGRPCServer() {
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
	service := &epochServiceserver{}

	infc.RegisterEpochServiceServer(grpcServer, service)

	grpcServer.Serve(lis)
}
