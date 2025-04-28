package postgresql

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/predixus/orca/internal/dag"
	pb "github.com/predixus/orca/protobufs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Datalayer struct {
	queries *Queries
	conn    *pgx.Conn
	closeFn func(context.Context) error
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

// CreateProcessor add a processor to the Orca server
func (d *Datalayer) CreateProcessor(ctx context.Context, proc *pb.ProcessorRegistration) error {
	slog.Debug("creating processor", "protobuf", proc)

	tx, err := d.conn.Begin(ctx)
	if err != nil {
		slog.Error("could not start transaction when", "error", err)
		return err
	}
	defer tx.Rollback(ctx)
	qtx := d.queries.WithTx(tx)

	// register the processor
	err = qtx.CreateProcessorAndPurgeAlgos(ctx, CreateProcessorAndPurgeAlgosParams{
		Name:             proc.GetName(),
		Runtime:          proc.GetRuntime(),
		ConnectionString: proc.GetConnectionStr(),
	})
	if err != nil {
		slog.Error("could not create processor", "error", err)
		return err
	}

	// add all algorithms first
	for _, algo := range proc.GetSupportedAlgorithms() {
		// add window types
		window_type := algo.GetWindowType()

		err := qtx.CreateWindowType(ctx, CreateWindowTypeParams{
			Name:    window_type.Name,
			Version: window_type.Version,
		})
		if err != nil {
			slog.Error("could not create window type", "error", err)
		}

		// create algos
		err = qtx.CreateAlgorithm(ctx,
			CreateAlgorithmParams{
				Name:              algo.GetName(),
				Version:           algo.GetVersion(),
				ProcessorName:     proc.GetName(),
				ProcessorRuntime:  proc.GetRuntime(),
				WindowTypeName:    algo.GetWindowType().GetName(),
				WindowTypeVersion: algo.GetWindowType().GetVersion(),
			})
		if err != nil {
			slog.Error("error creating algorithm", "error", err)
			return err
		}
	}

	// then add the dependencies and associate the processor with all the algos
	for _, algo := range proc.GetSupportedAlgorithms() {

		dependencies := algo.GetDependencies()
		for _, algoDependentOn := range dependencies {
			err := qtx.CreateAlgorithmDependency(ctx, CreateAlgorithmDependencyParams{
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

		err := qtx.AddProcessorAlgorithm(ctx, AddProcessorAlgorithmParams{
			ProcessorName:    proc.GetName(),
			ProcessorRuntime: proc.GetRuntime(),
			AlgorithmName:    algo.GetName(),
			AlgorithmVersion: algo.GetVersion(),
		})
		if err != nil {
			slog.Error(
				"could not associate algo with processor",
				"proc",
				proc,
				"algo",
				algo,
				"error",
				err,
			)
			return err
		}
	}

	return tx.Commit(ctx)
}

func (d *Datalayer) EmitWindow(ctx context.Context, window *pb.Window) error {
	slog.Info("recieved window")

	slog.Debug("inserting window", "window", window)
	insertedWindow, err := d.queries.RegisterWindow(ctx, RegisterWindowParams{
		WindowTypeName:    window.GetWindowTypeName(),
		WindowTypeVersion: window.GetWindowTypeVersion(),
		TimeFrom:          int64(window.GetFrom()),
		TimeTo:            int64(window.GetTo()),
		Origin:            window.GetOrigin(),
	})
	if err != nil {
		slog.Error("could not insert window", "error", err)
		if strings.Contains(err.Error(), "(SQLSTATE 23503)") {
			return fmt.Errorf(
				"window type does not exist - insert via window type registration: %v", err.Error(),
			)
		}
		return err
	}
	slog.Debug("window record inserted into the datalayer", "window", insertedWindow)
	exec_paths, err := d.queries.ReadAlgorithmExecutionPaths(ctx, strconv.Itoa(int(insertedWindow)))
	if err != nil {
		slog.Error(
			"could not read execution paths for window id",
			"window_id",
			insertedWindow,
			"error",
			err,
		)
		return err
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
		int64(insertedWindow),
	)
	if err != nil {
		slog.Error(
			"failed to construct execution paths for window",
			"window",
			insertedWindow,
			"error",
			err,
		)
		return err
	}

	slog.Info("calculated execution paths", "execution_paths", executionPlan)
	// get map of processors from processor ids
	processorMap := make(
		map[int64]Processor,
		len(executionPlan.AffectedProcessors),
	)
	processors, err := d.queries.ReadProcessorsByIDs(ctx, executionPlan.AffectedProcessors)
	for _, proc := range processors {
		processorMap[proc.ID] = proc
	}

	// get map of algorithms from algorithm ids
	algorithmMap := make(
		map[int64]Algorithm,
	)
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

			// Build list of affected Algorithms
			var affectedAlgorithms []*pb.Algorithm
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
			}

			execReq := &pb.ExecutionRequest{
				Window:           window,
				AlgorithmResults: nil, // TODO: dependency handling
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

			// Receive streamed execution results
			for {
				result, err := stream.Recv()
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
					"task_id", result.GetTaskId(),
					"status", result.GetStatus().String(),
				)

				// TODO: handle result persistence if needed
			}
		}
	}
	return nil
}
