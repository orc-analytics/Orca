package postgresql

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	types "github.com/predixus/orca/core/internal/types"
	pb "github.com/predixus/orca/core/protobufs/go"
)

type Datalayer struct {
	queries *Queries
	conn    *pgx.Conn
	closeFn func(context.Context) error
}

type PgTx struct {
	tx pgx.Tx
}

func (t *PgTx) Rollback(ctx context.Context) {
	t.tx.Rollback(ctx)
}

func (t *PgTx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

// generate a new client for the postgres datalayer
func NewClient(ctx context.Context, connStr string) (*Datalayer, error) {
	if connStr == "" {
		return nil, errors.New("connection string empty")
	}

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		slog.Error("Issue connecting to postgres", "error", err)
		return nil, err
	}

	return &Datalayer{
		queries: New(conn),
		conn:    conn,
		closeFn: conn.Close,
	}, nil
}

func (d *Datalayer) WithTx(ctx context.Context) (types.Tx, error) {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		slog.Error("could not start transaction", "error", err)
		return nil, err
	}
	return &PgTx{tx: tx}, nil
}

func (d *Datalayer) createProcessorAndPurgeAlgos(
	ctx context.Context,
	tx types.Tx,
	proc *pb.ProcessorRegistration,
) error {
	pgTx := tx.(*PgTx)

	qtx := d.queries.WithTx(pgTx.tx)
	// register the processor
	err := qtx.CreateProcessorAndPurgeAlgos(ctx, CreateProcessorAndPurgeAlgosParams{
		Name:             proc.GetName(),
		Runtime:          proc.GetRuntime(),
		ConnectionString: proc.GetConnectionStr(),
	})
	if err != nil {
		slog.Error("could not create processor", "error", err)
		return err
	}
	return nil
}

func (d *Datalayer) createWindowType(
	ctx context.Context,
	tx types.Tx,
	windowType *pb.WindowType,
) error {
	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)
	err := qtx.CreateWindowType(ctx, CreateWindowTypeParams{
		Name:    windowType.Name,
		Version: windowType.Version,
	})
	if err != nil {
		slog.Error("could not create window type", "error", err)
		return err
	}
	return nil
}

func (d *Datalayer) addAlgorithm(
	ctx context.Context,
	tx types.Tx,
	algo *pb.Algorithm,
	proc *pb.ProcessorRegistration,
) error {
	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)
	// create algos
	params := CreateAlgorithmParams{
		Name:              algo.GetName(),
		Version:           algo.GetVersion(),
		ProcessorName:     proc.GetName(),
		ProcessorRuntime:  proc.GetRuntime(),
		WindowTypeName:    algo.GetWindowType().GetName(),
		WindowTypeVersion: algo.GetWindowType().GetVersion(),
	}
	err := qtx.CreateAlgorithm(ctx, params)
	if err != nil {
		slog.Error("error creating algorithm", "error", err)
		return err
	}
	return nil
}

func (d *Datalayer) addOverwriteAlgorithmDependency(
	ctx context.Context,
	tx types.Tx,
	algo *pb.Algorithm,
	proc *pb.ProcessorRegistration,
) error {
	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)
	// get algorithm id
	algoId, err := qtx.ReadAlgorithmId(ctx, ReadAlgorithmIdParams{
		AlgorithmName:    algo.GetName(),
		AlgorithmVersion: algo.GetVersion(),
		ProcessorName:    proc.GetName(),
		ProcessorRuntime: proc.GetRuntime(),
	})
	if err != nil {
		slog.Error("could not get algorithm ID", "algorithm", algo)
		return err
	}
	dependencies := algo.GetDependencies()
	for _, algoDependentOn := range dependencies {
		// get algorithm id
		algoDependentOnId, err := qtx.ReadAlgorithmId(ctx, ReadAlgorithmIdParams{
			AlgorithmName:    algoDependentOn.GetName(),
			AlgorithmVersion: algoDependentOn.GetVersion(),
			ProcessorName:    proc.GetName(),
			ProcessorRuntime: proc.GetRuntime(),
		})
		if err != nil {
			slog.Error("could not get algorithm ID of dependant", "algorithm", algoDependentOn)
			return err
		}

		// get the algo execution path
		execPaths, err := qtx.ReadAlgorithmExecutionPathsForAlgo(ctx, algoDependentOnId)
		if err != nil {
			slog.Error("could not obtain execution paths", "algorithm_id", algoDependentOnId)
			return err
		}
		for _, algoPath := range execPaths {
			algoIds := strings.Split(algoPath.AlgoIDPath, ".")
			if slices.Contains(algoIds, strconv.Itoa(int(algoId))) {
				slog.Error(
					"found circular dependency",
					"from_algo",
					algoDependentOn,
					"to_algo",
					algo,
				)
				return &types.CircularDependencyError{
					FromAlgoName:      algoDependentOn.GetName(),
					ToAlgoName:        algo.GetName(),
					FromAlgoVersion:   algoDependentOn.GetVersion(),
					ToAlgoVersion:     algo.GetVersion(),
					FromAlgoProcessor: proc.GetName(),
					ToAlgoProcessor:   algoDependentOn.GetProcessorName(),
				}
			} else {
				err = qtx.CreateAlgorithmDependency(ctx, CreateAlgorithmDependencyParams{
					FromAlgorithmName:    algoDependentOn.GetName(),
					FromAlgorithmVersion: algoDependentOn.GetVersion(),
					FromProcessorName:    algoDependentOn.GetProcessorName(),
					FromProcessorRuntime: algoDependentOn.GetProcessorRuntime(),
					ToAlgorithmName:      algo.GetName(),
					ToAlgorithmVersion:   algo.GetVersion(),
					ToProcessorName:      proc.GetName(),
					ToProcessorRuntime:   proc.GetRuntime(),
				})
				if err != nil {
					slog.Error(
						"cloud not create algotrithm dependency",
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
	}
	return nil
}
