package internal

import (
	"context"
	"log/slog"

	dlyr "github.com/predixus/orca/internal/datalayers"
	pb "github.com/predixus/orca/protobufs/go"
	"google.golang.org/grpc"
)

type (
	OrcaCoreServer struct {
		pb.UnimplementedOrcaCoreServer
		client dlyr.Datalayer
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

// NewServer produces a new ORCA gRPC server
func NewServer(
	ctx context.Context,
	platform dlyr.Platform,
	connStr string,
) (*OrcaCoreServer, error) {
	client, err := dlyr.NewDatalayerClient(ctx, platform, connStr)
	if err != nil {
		slog.Error(
			"Could not initialise new platform client whilst initialising server",
			"platform",
			platform,
			"error",
			err,
		)

		return nil, err
	}

	s := &OrcaCoreServer{
		client: client,
	}
	return s, nil
}

// Register a processor with orca-core. Called when a processor startsup.
func (o *OrcaCoreServer) RegisterProcessor(
	ctx context.Context,
	proc *pb.ProcessorRegistration,
) (*pb.Status, error) {
	var status pb.Status

	slog.Info("registering processor",
		"runtime", proc.Runtime)

	err := o.client.AddProcessor(context.Background(), proc)
	if err != nil {
		status = pb.Status{
			Received: false,
			Message:  "Could not register processor",
		}
		return &status, err
	}
	return &pb.Status{
		Received: true,
		Message:  "Successfull registered processor",
	}, nil
}

func (o *OrcaCoreServer) EmitWindow(
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

func (o *OrcaCoreServer) RegisterWindowType(
	ctx context.Context,
	windowType *pb.WindowType,
) (*pb.Status, error) {
	slog.Info("registering window type",
		"name", windowType.Name)
	return &pb.Status{
		Received: true,
	}, nil
}

func (o *OrcaCoreServer) RegisterAlgorithm(
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

func (o *OrcaCoreServer) SubmitResult(
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

func (o *OrcaCoreServer) GetDagState(
	ctx context.Context,
	request *pb.DagStateRequest,
) (*pb.DagState, error) {
	slog.Info("getting DAG state",
		"window_id", request.WindowId)
	return &pb.DagState{}, nil
}
