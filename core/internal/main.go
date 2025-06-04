package internal

import (
	"context"
	"log/slog"

	"github.com/bufbuild/protovalidate-go"
	dlyr "github.com/predixus/orca/core/internal/datalayers"
	types "github.com/predixus/orca/core/internal/types"
	pb "github.com/predixus/orca/core/protobufs/go"
	"google.golang.org/protobuf/proto"
)

type (
	OrcaCoreServer struct {
		pb.UnimplementedOrcaCoreServer
		client types.Datalayer
	}
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

// Orca gRPC server function implementations
func (o *OrcaCoreServer) RegisterWindowType(
	ctx context.Context,
	w *pb.WindowRegistration,
) (*pb.WindowTypeRegResponse, error) {
	err := validate(w)
	if err != nil {
		return nil, err
	}
	slog.Info("registering window type")

	tx, err := o.client.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}

	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	err = o.client.CreateWindowType(ctx, tx, w.GetWindowType())
	if err != nil {
		return nil, err
	}
	slog.Info("Registered window")
	return &pb.WindowTypeRegResponse{}, tx.Commit(ctx)
}

// Register a processor with orca-core. Called when a processor startsup.
func (o *OrcaCoreServer) RegisterProcessor(
	ctx context.Context,
	proc *pb.ProcessorRegistration,
) (*pb.ProcRegResponse, error) {
	err := validate(proc)
	if err != nil {
		return nil, err
	}
	slog.Info("registering processor")

	tx, err := o.client.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}

	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	// register/refresh the processor
	err = o.client.RefreshProcessor(ctx, tx, proc)
	if err != nil {
		slog.Error("could not create processor", "error", err)
		return nil, err
	}

	// add all algorithms first
	for _, algo := range proc.GetSupportedAlgorithms() {
		// create algos
		err = o.client.AddAlgorithm(ctx, tx, algo, proc)
		if err != nil {
			slog.Error("error creating algorithm", "error", err)
			return nil, err
		}
	}

	// then add the dependencies and associate the processor with all the algos
	for _, algo := range proc.GetSupportedAlgorithms() {

		dependencies := algo.GetDependencies()
		for _, algoDependentOn := range dependencies {
			err := o.client.AddOverwriteAlgorithmDependency(
				ctx,
				tx,
				algo,
				proc,
			)
			if err != nil {
				slog.Error(
					"could not create algotrithm dependency",
					"algorithm",
					algo,
					"depends_on",
					algoDependentOn,
					"error",
					err,
				)
				return nil, err
			}
		}
	}

	// then add datagetters
	for _, dg := range proc.GetDataGetters() {
		err := o.client.AddOverwriteDataGetter(ctx, tx, dg, proc)
		if err != nil {
			slog.Error("could not create data getter", "data getter", dg, "error", err)
			return nil, err
		}
	}

	slog.Info("registered processor")
	return &pb.ProcRegResponse{}, tx.Commit(ctx)
}

func (o *OrcaCoreServer) EmitWindow(
	ctx context.Context,
	window *pb.Window,
) (*pb.WindowEmitResponse, error) {
	err := validate(window)
	if err != nil {
		return nil, err
	}
	slog.Info("emitting window", "window", window)
	windowEmitStatus, err := o.client.EmitWindow(ctx, window)
	return &windowEmitStatus, err
}
