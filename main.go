package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	pb "github.com/predixus/orca/protobufs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type (
	orcaService struct {
		pb.UnimplementedOrcaServiceServer
	}
)

func (w *orcaService) RegisterWindow(ctx context.Context, window *pb.Window) (*pb.Status, error) {
	slog.Debug("Recieved window", "window", window)
	return &pb.Status{
		Recieved: true,
	}, nil
}

func (w *orcaService) RegisterWindowType(
	ctx context.Context,
	windowType *pb.WindowType,
) (*pb.Status, error) {
	slog.Debug("Recieved window type", "windowType", windowType)
	return &pb.Status{
		Recieved: true,
	}, nil
}

func (w *orcaService) RegisterAlgorithmType(
	ctx context.Context,
	algorithmType *pb.AlgorithmType,
) (*pb.Status, error) {
	slog.Debug("Recieved algorithm type", "AlgorithmType", algorithmType)
	return &pb.Status{
		Recieved: true,
	}, nil
}

func (w *orcaService) RegisterResult(ctx context.Context, result *pb.Result) (*pb.Status, error) {
	slog.Debug("Recieved result", "result", result)
	return &pb.Status{
		Recieved: true,
	}, nil
}

func newServer() *orcaService {
	s := &orcaService{}
	return s
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	port := 4040
	slog.Debug("Running the server", "port", port)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		slog.Error("failed to listen", "message", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterOrcaServiceServer(grpcServer, newServer())
	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		slog.Error("message", err)
	}
}
