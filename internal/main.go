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

	err = o.client.CreateProcessor(context.Background(), proc)
	if err != nil {
		return nil, err
	}
	slog.Debug("registered processor", "processor", proc)
	return &pb.Status{
		Received: true,
		Message:  "Successfully registered processor",
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
