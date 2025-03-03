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

func startGRPCServer(dbConnString string, port int) {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	// defer f.Close()

	// Configure slog to use the same file
	logger := slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	go func() {
		slog.Info("Launching server", "port", port)
		slog.Debug("Debugging")
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
	}()
}

func main() {
	Run()
}
