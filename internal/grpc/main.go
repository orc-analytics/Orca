package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/predixus/pdb_framework/internal/grpc/epoch"
	infc "github.com/predixus/pdb_framework/protobufs/go"
	"google.golang.org/grpc"
)

type GRPCServer interface {
	Start(wg *sync.WaitGroup)
}

type GrpcServer struct {
	*grpc.Server
}

// implement the GRPC server requests
func (s *GrpcServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "8080"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", os.Getenv("GRPC_PORT")))
	if err != nil {
		log.Fatalf("Cannot create listener: %s", err)
	}

	// register your gRPC service here
	service := &grpc_epoch.EpochServiceServer{}
	infc.RegisterEpochServiceServer(s.Server, service)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Println("GRPC: Shutting Down")
}
