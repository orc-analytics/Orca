package internal

import (
	"context"
	"log/slog"

	"github.com/bufbuild/protovalidate-go"
	dlyr "github.com/orc-analytics/orca/core/internal/datalayers"
	types "github.com/orc-analytics/orca/core/internal/types"
	pb "github.com/orc-analytics/orca/core/protobufs/go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type (
	OrcaCoreServer struct {
		pb.UnimplementedOrcaCoreServer
		client types.Datalayer
	}
)

var (
	MAX_PROCESSORS = 20
	processors     = make(
		[]grpc.ServerStreamingServer[pb.ProcessingTask],
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

// --------------------------- gRPC Services ---------------------------
// -------------------------- Core Operations --------------------------
// Register a processor with orca-core. Called when a processor startsup.
func (o *OrcaCoreServer) RegisterProcessor(
	ctx context.Context,
	proc *pb.ProcessorRegistration,
) (*pb.Status, error) {
	err := validate(proc)
	if err != nil {
		return nil, err
	}
	slog.Info("registering processor")
	err = o.client.RegisterProcessor(ctx, proc)
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
	slog.Debug("Recieved Window", "window", window)
	err := validate(window)
	if err != nil {
		return nil, err
	}
	slog.Info("emitting window", "window", window)
	windowEmitStatus, err := o.client.EmitWindow(ctx, window)
	return &windowEmitStatus, err
}

// -------------------------- Data Operations --------------------------
func (o *OrcaCoreServer) ReadWindowTypes(
	ctx context.Context,
	windowTypeReadStub *pb.WindowTypeRead,
) (*pb.WindowTypes, error) {
	return o.client.ReadWindowTypes(ctx)
}

func (o *OrcaCoreServer) ReadAlgorithms(
	ctx context.Context,
	algorithmsReadStub *pb.AlgorithmsRead,
) (*pb.Algorithms, error) {
	return o.client.ReadAlgorithms(ctx)
}

func (o *OrcaCoreServer) ReadProcessors(
	ctx context.Context,
	processorsReadStub *pb.ProcessorsRead,
) (*pb.Processors, error) {
	return o.client.ReadProcessors(ctx)
}

func (o *OrcaCoreServer) ReadResultsStats(
	ctx context.Context,
	ResultsStatsReadStub *pb.ResultsStatsRead,
) (*pb.ResultsStats, error) {
	return o.client.ReadResultsStats(ctx)
}

func (o *OrcaCoreServer) ReadResultFieldsForAlgorithm(
	ctx context.Context,
	algorithmFieldsRead *pb.AlgorithmFieldsRead,
) (*pb.AlgorithmFields, error) {
	return o.client.ReadResultFieldsForAlgorithm(ctx, algorithmFieldsRead)
}

func (o *OrcaCoreServer) ReadResultsForAlgorithm(
	ctx context.Context,
	resultsForAlgorithmRead *pb.ResultsForAlgorithmRead,
) (*pb.ResultsForAlgorithm, error) {
	return o.client.ReadResultsForAlgorithm(ctx, resultsForAlgorithmRead)
}

func (o *OrcaCoreServer) ReadWindows(
	ctx context.Context,
	windowReads *pb.WindowsRead,
) (*pb.Windows, error) {
	return o.client.ReadWindows(ctx, windowReads)
}

func (o *OrcaCoreServer) ReadDistinctMetadataForWindowType(
	ctx context.Context,
	windowMetadataRead *pb.DistinctMetadataForWindowTypeRead,
) (*pb.DistinctMetadataForWindowType, error) {
	return o.client.ReadDistinctMetadataForWindowType(ctx, windowMetadataRead)
}

func (o *OrcaCoreServer) ReadWindowsForMetadata(
	ctx context.Context,
	windowsForMetadataRead *pb.WindowsForMetadataRead,
) (*pb.WindowsForMetadata, error) {
	return o.client.ReadWindowsForMetadata(ctx, windowsForMetadataRead)
}

func (o *OrcaCoreServer) ReadResultsForAlgorithmAndMetadata(
	ctx context.Context,
	resultsForAlgorithmAndMetadata *pb.ResultsForAlgorithmAndMetadataRead,
) (*pb.ResultsForAlgorithmAndMetadata, error) {
	return o.client.ReadResultsForAlgorithmAndMetadata(ctx, resultsForAlgorithmAndMetadata)
}

// ---------------------- Labelling Operations ----------------------
func (o *OrcaCoreServer) Annotate(
	ctx context.Context,
	annotateWrite *pb.AnnotateWrite,
) (*pb.AnnotateResponse, error) {
	return o.client.Annotate(ctx, annotateWrite)
}
