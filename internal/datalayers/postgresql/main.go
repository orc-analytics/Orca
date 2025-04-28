package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/predixus/orca/internal/dag"
	pb "github.com/predixus/orca/protobufs/go"
	// "google.golang.org/grpc"
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
	executionPaths, err := dag.GetPathsForWindow(
		algoIdPaths,
		windowTypeIDPaths,
		procIdPaths,
		int(insertedWindow),
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

	slog.Info("calculated unique execution paths", "execution_paths", executionPaths)
	// get the processor details
	// processorIds := make([]int64, len(executionPaths), len(executionPaths))
	// for ii, execPath := range executionPaths {
	// 	processorIds[ii] = int64(execPath.ProcessorId)
	// }
	//
	// processors, err := d.queries.ReadProcessorsByIDs(ctx, processorIds)
	// // map exec paths to processors
	// execMap := make(map[dag.ExecutionPath]Processor, len(executionPaths))
	// for _, path := range executionPaths {
	// 	for _, processor := range processors {
	// 		if path.ProcessorId == int(processor.ID) {
	// 			execMap[path] = processor
	// 			break
	// 		}
	// 	}
	// }
	// slog.Info("built processor execution map", "exec map", execMap)
	// for _, execPath := range executionPaths {
	// 	conn, err := grpc.NewClient(
	// 		execMap[execPath].ConnectionString,
	// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 	)
	// 	if err != nil {
	// 		slog.Error("could not connect to processor", "error", err)
	// 		return err
	// 	}
	// 	defer conn.Close()
	//
	// 	client := pb.NewOrcaProcessorClient(conn)
	// 	stream, err := client.ExecuteDagPart(ctx, &pb.ExecutionRequest{
	// 		Window: window,
	// 		// TODO handle results from other processors
	// 	})
	// 	if err != nil {
	// 		// handle error
	// 	}
	//
	// 	for {
	// 		result, err := stream.Recv()
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		if err != nil {
	// 			// handle error
	// 		}
	// 		// process result
	// 	}
	// }
	return nil
}
