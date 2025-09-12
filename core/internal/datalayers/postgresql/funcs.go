package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	types "github.com/orc-analytics/orca/core/internal/types"
	pb "github.com/orc-analytics/orca/core/protobufs/go"
)

type Datalayer struct {
	queries *Queries
	conn    *pgxpool.Pool
	closeFn func()
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

	connPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		slog.Error("Issue connecting to postgres", "error", err)
		return nil, err
	}

	return &Datalayer{
		queries: New(connPool),
		conn:    connPool,
		closeFn: connPool.Close,
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
		Name:        windowType.GetName(),
		Version:     windowType.GetVersion(),
		Description: windowType.GetDescription(),
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
	var resultType ResultType
	if algo.GetResultType() == pb.ResultType_ARRAY {
		resultType = ResultTypeArray
	} else if algo.GetResultType() == pb.ResultType_VALUE {
		resultType = ResultTypeValue
	} else if algo.GetResultType() == pb.ResultType_STRUCT {
		resultType = ResultTypeStruct
	} else if algo.GetResultType() == pb.ResultType_NONE {
		resultType = ResultTypeNone
	} else {
		return fmt.Errorf("result type %v not supported", algo.GetResultType())
	}

	params := CreateAlgorithmParams{
		Name:              algo.GetName(),
		Version:           algo.GetVersion(),
		ProcessorName:     proc.GetName(),
		ProcessorRuntime:  proc.GetRuntime(),
		WindowTypeName:    algo.GetWindowType().GetName(),
		WindowTypeVersion: algo.GetWindowType().GetVersion(),
		ResultType:        resultType,
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
