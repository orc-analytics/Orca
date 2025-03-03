package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	orca "github.com/predixus/orca/internal"
	pb "github.com/predixus/orca/protobufs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}

func startGRPCServer(dbConnString string, port int) error {
	slog.Debug("Running the server", "port", port)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		slog.Error("failed to listen", "message", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterOrcaCoreServer(grpcServer, orca.NewServer())
	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		slog.Error("failed to serve", "error", err)
	}
	return err
}

func main() {
	Run()
}
