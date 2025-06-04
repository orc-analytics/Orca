package datalayers

import (
	"context"
	"fmt"
	"log/slog"

	psql "github.com/predixus/orca/core/internal/datalayers/postgresql"
	types "github.com/predixus/orca/core/internal/types"
	pb "github.com/predixus/orca/core/protobufs/go"
)

// represents the supported database platforms as the datalayer
type Platform string

const (
	PostgreSQL Platform = "postgresql"
)

// check if the platform is supported
func (p Platform) isValid() bool {
	switch p {
	case PostgreSQL:
		return true
	default:
		return false
	}
}

func NewDatalayerClient(
	ctx context.Context,
	platform Platform,
	connStr string,
) (types.Datalayer, error) {
	if !platform.isValid() {
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}

	switch platform {
	case PostgreSQL:
		return psql.NewClient(ctx, connStr)
	default:
		slog.Error("datalayer not supported", "platform", platform)
		return nil, fmt.Errorf("platform not implemented: %s", platform)
	}
}

// wrapper funcs that stitch together the logic
func RegisterProcessor(
	ctx context.Context,
	dlyr types.Datalayer,
	proc *pb.ProcessorRegistration,
) error {
	slog.Debug("creating processor", "protobuf", proc)
	tx, err := dlyr.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return err
	}

	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	// register/refresh the processor
	err = dlyr.RefreshProcessor(ctx, tx, proc)
	if err != nil {
		slog.Error("could not create processor", "error", err)
		return err
	}

	// add all algorithms first
	for _, algo := range proc.GetSupportedAlgorithms() {
		// create algos
		err = dlyr.AddAlgorithm(ctx, tx, algo, proc)
		if err != nil {
			slog.Error("error creating algorithm", "error", err)
			return err
		}
	}

	// then add the dependencies and associate the processor with all the algos
	for _, algo := range proc.GetSupportedAlgorithms() {

		dependencies := algo.GetDependencies()
		for _, algoDependentOn := range dependencies {
			err := dlyr.AddOverwriteAlgorithmDependency(
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
				return err
			}
		}
	}

	// then add datagetters
	for _, dg := range proc.GetDataGetters() {
		err := dlyr.AddOverwriteDataGetter(ctx, tx, dg, proc)
		if err != nil {
			slog.Error("could not create data getter", "data getter", dg, "error", err)
		}
		return err
	}

	return tx.Commit(ctx)
}

func RegisterWindowType(ctx context.Context, dlyr types.Datalayer, w *pb.WindowRegistration) error {
	slog.Debug("registering window type", "window registration", w)

	tx, err := dlyr.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return err
	}

	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()
	err = dlyr.CreateWindowType(ctx, tx, w.GetWindowType())
	if err != nil {
		return err
	}
	return nil
}
