package internal

import (
	"context"
	"log/slog"

	"github.com/bufbuild/protovalidate-go"
	dlyr "github.com/predixus/orca/internal/datalayers"
	pb "github.com/predixus/orca/protobufs/go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
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

// validate a protobuf via protovalidate
func validate[T proto.Message](msg T) error {
	v, err := protovalidate.New()
	if err != nil {
		return err
	}

	if err := v.Validate(msg); err != nil {
		return err
	}

	return nil
}

// Register a processor with orca-core. Called when a processor startsup.
func (o *OrcaCoreServer) RegisterProcessor(
	ctx context.Context,
	proc *pb.ProcessorRegistration,
) (*pb.Status, error) {
	err := validate[*pb.ProcessorRegistration](proc)
	if err != nil {
		return nil, err
	}
	slog.Info("registering processor")
	slog.Debug("registered processor", "processor", proc)

	err = o.client.CreateProcessor(context.Background(), proc)
	if err != nil {
		return nil, err
	}
	return &pb.Status{
		Received: true,
		Message:  "Successfully registered processor",
	}, nil
}

func (o *OrcaCoreServer) RegisterWindowType(
	ctx context.Context,
	windowType *pb.WindowType,
) (*pb.WindowTypeRegisterStatus, error) {
	err := validate[*pb.WindowType](windowType)
	if err != nil {
		return nil, err
	}

	slog.Info("registering window type",
		"name", windowType.Name)
	err = o.client.RegisterWindowType(ctx, windowType)
	if err != nil {
		slog.Error("failed to register window type", "error", err)
		status := pb.WindowTypeRegisterStatus{
			Status:  pb.WindowTypeRegisterStatus_WINDOW_NOT_REGISTERED,
			Message: err.Error(),
		}
		return &status, err
	}
	slog.Debug("registered window type", "windowType", windowType)

	return &pb.WindowTypeRegisterStatus{
		Status:  pb.WindowTypeRegisterStatus_WINDOW_REGISTERED,
		Message: "window type registered",
	}, nil
}

func (o *OrcaCoreServer) EmitWindow(
	ctx context.Context,
	window *pb.Window,
) (*pb.WindowEmitStatus, error) {
	err := validate[*pb.Window](window)
	if err != nil {
		return nil, err
	}
	slog.Info("emitting window")
	err = o.client.EmitWindow(ctx, window)
	if err != nil {
		return &pb.WindowEmitStatus{
			Status: pb.WindowEmitStatus_NO_TRIGGERED_ALGORITHMS,
		}, err
	}

	// TODO: actually trigger some algos
	return &pb.WindowEmitStatus{
		Status: pb.WindowEmitStatus_NO_TRIGGERED_ALGORITHMS,
	}, nil
}

// func (o *OrcaCoreServer) EmitWindow(
// 	ctx context.Context,
// 	window *pb.Window,
// ) (*pb.WindowEmitStatus, error) {
// 	slog.Info("received window",
// 		"name", window.Name,
// 		"from", window.From,
// 		"to", window.To)
// 	return &pb.WindowEmitStatus{
// 		Status: pb.WindowEmitStatus_NO_TRIGGERED_ALGORITHMS,
// 	}, nil
// }
//
// func (o *OrcaCoreServer) RegisterAlgorithm(
// 	ctx context.Context,
// 	algorithm *pb.Algorithm,
// ) (*pb.Status, error) {
// 	slog.Info("registering algorithm",
// 		"name", algorithm.Name,
// 		"version", algorithm.Version)
// 	return &pb.Status{
// 		Received: true,
// 	}, nil
// }
//
// func (o *OrcaCoreServer) SubmitResult(
// 	ctx context.Context,
// 	result *pb.Result,
// ) (*pb.Status, error) {
// 	slog.Info("received result",
// 		"algorithm", result.AlgorithmName,
// 		"version", result.Version,
// 		"status", result.Status)
// 	return &pb.Status{
// 		Received: true,
// 	}, nil
// }
//
// func (o *OrcaCoreServer) GetDagState(
// 	ctx context.Context,
// 	request *pb.DagStateRequest,
// ) (*pb.DagState, error) {
// 	slog.Info("getting DAG state",
// 		"window_id", request.WindowId)
// 	return &pb.DagState{}, nil
// }
