package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	orca "github.com/predixus/orca/internal"
	dlyr "github.com/predixus/orca/internal/datalayers"
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

func startGRPCServer(platform string, dbConnString string, port int) {
	// Create .orca/logs directory if it doesn't exist
	logDir := ".orca/logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		slog.Error("Could not create log directory", "error", err)
		os.Exit(1)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	logPath := fmt.Sprintf("%s/debug-%s.log", logDir, timestamp)

	f, err := tea.LogToFile(logPath, "debug")
	if err != nil {
		slog.Error("Could not open logging file", err)
		os.Exit(1)
	}
	// Configure slog to use the same file
	logger := slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	orcaServer, err := orca.NewServer(context.Background(), dlyr.Platform(platform), dbConnString)
	if err != nil {
		slog.Error("Issue launching Orca Server", "error", err)
		os.Exit(1)
	}
	go func(server *orca.OrcaCoreServer) {
		slog.Info("Launching server", "port", port)
		defer f.Close()
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
	Run()
}
