package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	orca "github.com/predixus/orca/internal"
	dlyr "github.com/predixus/orca/internal/datalayers"
	pb "github.com/predixus/orca/protobufs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startGRPCServer(
	platform string,
	dbConnString string,
	port int,
	logLevel string,
) {
	orcaServer, err := orca.NewServer(context.Background(), dlyr.Platform(platform), dbConnString)
	if err != nil {
		slog.Error("issue launching Orca Server", "error", err)
		os.Exit(1)
	}
	go func(server *orca.OrcaCoreServer) {
		slog.Info("starting server", "port", port)
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			slog.Error("failed to listen", "message", err)
		}
		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)

		pb.RegisterOrcaCoreServer(
			grpcServer,
			server,
		)
		reflection.Register(grpcServer)
		err = grpcServer.Serve(lis)
		if err != nil {
			slog.Error("failed to serve", "error", err)
		}
	}(orcaServer)
}

func main() {
	flags := parseFlags()

	if err := validateFlags(flags); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	runCLI(flags)
}
