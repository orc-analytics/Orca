package internal

import (
	"context"
	"log/slog"

	pb "github.com/predixus/orca/protobufs/go"
	"google.golang.org/grpc"
)

type (
	orcaCoreServer struct {
		pb.UnimplementedOrcaCoreServer
	}
)

var (
	MAX_PROCESSORS = 20
	processors     = make(
		[]grpc.ServerStreamingServer[pb.ProcessingTask],
		MAX_PROCESSORS,
		MAX_PROCESSORS,
	)
)

// Register a processor with orca-core. Called on processor startup.
func (orcaCoreServer) RegisterProcessor(
	reg *pb.ProcessorRegistration,
	stream grpc.ServerStreamingServer[pb.ProcessingTask],
) error {
	slog.Info("registering processor",
		"runtime", reg.Runtime)

	// do stuff

	return nil
}

func (orcaCoreServer) EmitWindow(
	ctx context.Context,
	window *pb.Window,
) (*pb.WindowEmitStatus, error) {
	slog.Info("received window",
		"name", window.Name,
		"from", window.From,
		"to", window.To)
	return &pb.WindowEmitStatus{
		Status: pb.WindowEmitStatus_NO_TRIGGERED_ALGORITHMS,
	}, nil
}

func (orcaCoreServer) RegisterWindowType(
	ctx context.Context,
	windowType *pb.WindowType,
) (*pb.Status, error) {
	slog.Info("registering window type",
		"name", windowType.Name)
	return &pb.Status{
		Received: true,
	}, nil
}

func (orcaCoreServer) RegisterAlgorithm(
	ctx context.Context,
	algorithm *pb.Algorithm,
) (*pb.Status, error) {
	slog.Info("registering algorithm",
		"name", algorithm.Name,
		"version", algorithm.Version)
	return &pb.Status{
		Received: true,
	}, nil
}

func (orcaCoreServer) SubmitResult(
	ctx context.Context,
	result *pb.Result,
) (*pb.Status, error) {
	slog.Info("received result",
		"algorithm", result.AlgorithmName,
		"version", result.Version,
		"status", result.Status)
	return &pb.Status{
		Received: true,
	}, nil
}

func (orcaCoreServer) GetDagState(
	ctx context.Context,
	request *pb.DagStateRequest,
) (*pb.DagState, error) {
	slog.Info("getting DAG state",
		"window_id", request.WindowId)
	return &pb.DagState{}, nil
}

func NewServer() *orcaCoreServer {
	s := &orcaCoreServer{}
	return s
}
