package postgresql

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/predixus/orca/core/internal/dag"
	types "github.com/predixus/orca/core/internal/types"
	pb "github.com/predixus/orca/core/protobufs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
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

func (d *Datalayer) RefreshProcessor(
	ctx context.Context,
	tx types.Tx,
	proc *pb.ProcessorRegistration,
) error {
	pgTx := tx.(*PgTx)

	qtx := d.queries.WithTx(pgTx.tx)

	// delete any existing processors by the same name
	err := qtx.DeleteProcessor(ctx, DeleteProcessorParams{
		Name:    proc.GetName(),
		Runtime: proc.GetRuntime(),
	})
	if err != nil {
		slog.Error("issue deleteing processor", "processor", proc, "error", err)
		return err
	}

	// register the processor
	err = qtx.CreateProcessor(ctx, CreateProcessorParams{
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

func (d *Datalayer) CreateWindowType(
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

func (d *Datalayer) AddAlgorithm(
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

func (d *Datalayer) AddOverwriteAlgorithmDependency(
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

func (d *Datalayer) AddOverwriteDataGetter(
	ctx context.Context,
	tx types.Tx,
	dg *pb.DataGetter,
	proc *pb.ProcessorRegistration,
) error {
	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)
	_, err := qtx.CreateDataGetter(ctx, CreateDataGetterParams{
		Name:             dg.GetName(),
		WindowName:       dg.WindowType.GetName(),
		WindowVersion:    dg.WindowType.GetVersion(),
		TtlSeconds:       int64(dg.GetTtlSeconds()),
		MaxSizeBytes:     int64(dg.GetMaxSizeBytes()),
		ProcessorName:    proc.GetName(),
		ProcessorRuntime: proc.GetRuntime(),
	})
	if err != nil {
		slog.Error("issue adding data getter", "error", err)
		return err
	}
	return nil
}

func (d *Datalayer) EmitWindow(
	ctx context.Context,
	window *pb.Window,
) (pb.WindowEmitStatus, error) {
	slog.Info("recieved window")

	slog.Debug("inserting window", "window", window)
	insertedWindow, err := d.queries.RegisterWindow(ctx, RegisterWindowParams{
		WindowTypeName:    window.GetWindowTypeName(),
		WindowTypeVersion: window.GetWindowTypeVersion(),
		TimeFrom:          int64(window.GetTimeFrom()),
		TimeTo:            int64(window.GetTimeTo()),
		Origin:            window.GetOrigin(),
	})
	if err != nil {
		slog.Error("could not insert window", "error", err)
		if strings.Contains(err.Error(), "(SQLSTATE 23503)") {
			return pb.WindowEmitStatus{
					Status: pb.WindowEmitStatus_TRIGGERING_FAILED,
				}, fmt.Errorf(
					"window type does not exist - insert via window type registration: %v",
					err.Error(),
				)
		}
	}
	slog.Debug("window record inserted into the datalayer", "window", insertedWindow)
	exec_paths, err := d.queries.ReadAlgorithmExecutionPaths(
		ctx,
		strconv.Itoa(int(insertedWindow.WindowTypeID)),
	)
	if err != nil {
		slog.Error(
			"could not read execution paths for window id",
			"window_id",
			insertedWindow,
			"error",
			err,
		)
		return pb.WindowEmitStatus{Status: pb.WindowEmitStatus_TRIGGERING_FAILED}, err
	}

	// create the algo path args
	var algoIdPaths []string
	var windowTypeIDPaths []string
	var procIdPaths []string
	for _, path := range exec_paths {
		algoIdPaths = append(algoIdPaths, path.AlgoIDPath)
		windowTypeIDPaths = append(windowTypeIDPaths, path.WindowTypeIDPath)
		procIdPaths = append(procIdPaths, path.ProcIDPath)
	}

	// fire off processings
	executionPlan, err := dag.BuildPlan(
		algoIdPaths,
		windowTypeIDPaths,
		procIdPaths,
		int64(insertedWindow.WindowTypeID),
	)
	if err != nil {
		slog.Error(
			"failed to construct execution paths for window",
			"window",
			insertedWindow,
			"error",
			err,
		)
		return pb.WindowEmitStatus{Status: pb.WindowEmitStatus_TRIGGERING_FAILED}, err
	}

	if len(executionPlan.Stages) > 0 {
		go processTasks(d, executionPlan, window, insertedWindow)

		return pb.WindowEmitStatus{
			Status: pb.WindowEmitStatus_PROCESSING_TRIGGERED,
		}, nil
	} else {
		return pb.WindowEmitStatus{
			Status: pb.WindowEmitStatus_NO_TRIGGERED_ALGORITHMS,
		}, nil
	}
}

func processTasks(
	d *Datalayer,
	executionPlan dag.Plan,
	window *pb.Window,
	insertedWindow RegisterWindowRow,
) error {
	ctx := context.Background()
	slog.Info("calculated execution paths", "execution_paths", executionPlan)
	// get map of processors from processor ids
	processorMap := make(
		map[int64]Processor,
		len(executionPlan.AffectedProcessors),
	)
	processors, err := d.queries.ReadProcessorsByIDs(ctx, executionPlan.AffectedProcessors)
	if err != nil {
		slog.Error("Processors could not be read", "error", err)
		return err
	}

	for _, proc := range processors {
		processorMap[proc.ID] = proc
	}

	// get map of algorithms from algorithm ids
	algorithmMap := make(
		map[int64]Algorithm,
	)

	// map of algorithm Ids to results
	resultMap := make(
		map[int64]*pb.ExecutionResult,
	)

	// map of execution IDs and the algorithms requested
	algorithms, err := d.queries.ReadAlgorithmsForWindow(ctx, ReadAlgorithmsForWindowParams{
		WindowTypeName:    window.WindowTypeName,
		WindowTypeVersion: window.WindowTypeVersion,
	})

	for _, algo := range algorithms {
		algorithmMap[algo.ID] = algo
	}

	// for each stage, farm off processsings
	for _, stage := range executionPlan.Stages {
		for _, task := range stage.Tasks {
			proc, ok := processorMap[task.ProcId]
			if !ok {
				slog.Error("Processor not found for task", "proc_id", task.ProcId)
				return fmt.Errorf("processor ID %d not found", task.ProcId)
			}

			conn, err := grpc.NewClient(
				proc.ConnectionString,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				slog.Error("could not connect to processor", "proc_id", task.ProcId, "error", err)
				return err
			}
			// IMPORTANT: close conn when done (not deferred inside a loop)
			defer func(conn *grpc.ClientConn) {
				if err := conn.Close(); err != nil {
					slog.Warn("error closing gRPC connection", "error", err)
				}
			}(conn)

			client := pb.NewOrcaProcessorClient(conn)
			healthCheckResponse, err := client.HealthCheck(ctx, &pb.HealthCheckRequest{
				Timestamp: time.Now().Unix(),
			})
			if err != nil {
				slog.Error(
					"issue contacting processor",
					"response",
					healthCheckResponse,
					"processor",
					proc,
				)
				return err
			}
			if healthCheckResponse.Status != pb.HealthCheckResponse_STATUS_SERVING {
				slog.Error(
					"cannot execute stage, processor not serving",
					"status",
					healthCheckResponse.Status,
					"message",
					healthCheckResponse.Message,
				)
				return err
			}

			// build list of affected Algorithms
			var affectedAlgorithms []*pb.Algorithm

			// and their dependency's result
			algoDepsResults := []*pb.AlgorithmResult{}

			// generate an execution id
			execUuid := uuid.New()
			execId := strings.ReplaceAll(execUuid.String(), "-", "")

			for _, node := range task.Nodes {
				algo, ok := algorithmMap[node.AlgoId()]

				if !ok {
					slog.Error("algorithm not found", "algo_id", node.AlgoId())
					return fmt.Errorf("algorithm ID %d not found", node.AlgoId())
				}

				affectedAlgorithms = append(affectedAlgorithms, &pb.Algorithm{
					Name:    algo.Name,
					Version: algo.Version,
				})

				// determine which results need to be included
				for _, algoId := range node.AlgoDepIds() {
					algoDepsResults = append(algoDepsResults, resultMap[algoId].AlgorithmResult)
				}
			}

			execReq := &pb.ExecutionRequest{
				ExecId:           execId,
				Window:           window,
				AlgorithmResults: algoDepsResults,
				Algorithms:       affectedAlgorithms,
			}

			stream, err := client.ExecuteDagPart(ctx, execReq)
			if err != nil {
				slog.Error(
					"failed to start DAG part execution",
					"proc_id",
					task.ProcId,
					"error",
					err,
				)
				return err
			}

			// recieve streamed execution results
			for {
				result, err := stream.Recv()
				// error handling
				if err != nil {
					if errors.Is(err, context.Canceled) ||
						errors.Is(err, context.DeadlineExceeded) {
						slog.Warn(
							"context done while receiving execution result",
							"proc_id",
							task.ProcId,
						)
						break
					}
					if err == io.EOF {
						slog.Info("finished receiving execution results", "proc_id", task.ProcId)
						break
					}
					slog.Error(
						"error receiving execution result",
						"proc_id",
						task.ProcId,
						"error",
						err,
					)
					return err
				}

				slog.Info("received execution result",
					"exec_id", result.GetExecId(),
				)

				var algoResultId int
				for _, algo := range algorithms {
					if (algo.Name == result.AlgorithmResult.GetAlgorithm().Name) &&
						(algo.Version == result.AlgorithmResult.GetAlgorithm().Version) {
						algoResultId = int(algo.ID)
						break
					}
				}

				// add the result in to the result map
				resultMap[int64(algoResultId)] = result

				structResult, err := convertStructToJsonBytes(
					result.AlgorithmResult.Result.GetStructValue(),
				)
				if err != nil {
					slog.Error(
						"Issue converted algorithm struct result to bytes",
						"error",
						err,
						"struct",
						result.AlgorithmResult.Result.GetStructValue(),
					)
					return err
				}

				resultId, err := d.queries.CreateResult(ctx, CreateResultParams{
					WindowsID:    pgtype.Int8{Valid: true, Int64: insertedWindow.ID},
					WindowTypeID: pgtype.Int8{Valid: true, Int64: insertedWindow.WindowTypeID},
					AlgorithmID:  pgtype.Int8{Valid: true, Int64: int64(algoResultId)},
					ResultValue: pgtype.Float8{
						Valid:   true,
						Float64: float64(result.AlgorithmResult.Result.GetSingleValue()),
					},
					ResultArray: convertFloat32ToFloat64(
						result.AlgorithmResult.Result.GetFloatValues().GetValues(),
					),
					ResultJson: structResult,
				})
				if err != nil {
					slog.Error("Error inserting result", "error", err)
					return err
				}
				slog.Info("Inserted result", "resultId", resultId)
			}
		}
	}
	return nil
}

func convertFloat32ToFloat64(float32Slice []float32) []float64 {
	float64Slice := make([]float64, len(float32Slice), len(float32Slice))
	for i, value := range float32Slice {
		float64Slice[i] = float64(value)
	}
	return float64Slice
}

func convertStructToJsonBytes(s *structpb.Struct) ([]byte, error) {
	return protojson.Marshal(s)
}
