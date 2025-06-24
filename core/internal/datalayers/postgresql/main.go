package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/predixus/orca/core/internal/dag"
	pb "github.com/predixus/orca/core/protobufs/go"
)

func (d *Datalayer) RegisterProcessor(
	ctx context.Context,
	proc *pb.ProcessorRegistration,
) error {
	slog.Debug("creating processor", "protobuf", proc)
	tx, err := d.WithTx(ctx)

	defer tx.Rollback(ctx)

	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return err
	}

	// register the processor
	err = d.createProcessorAndPurgeAlgos(ctx, tx, proc)
	if err != nil {
		slog.Error("could not create processor", "error", err)
		return err
	}

	// add all algorithms first
	for _, algo := range proc.GetSupportedAlgorithms() {
		// add window types
		window_type := algo.GetWindowType()

		err := d.createWindowType(ctx, tx, window_type)
		if err != nil {
			slog.Error("could not create window type", "error", err)
			return err
		}

		// create algos
		err = d.addAlgorithm(ctx, tx, algo, proc)
		if err != nil {
			slog.Error("error creating algorithm", "error", err)
			return err
		}
	}

	// then add the dependencies and associate the processor with all the algos
	for _, algo := range proc.GetSupportedAlgorithms() {

		dependencies := algo.GetDependencies()
		for _, algoDependentOn := range dependencies {
			err := d.addOverwriteAlgorithmDependency(
				ctx,
				tx,
				algo,
				proc,
			)
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

	return tx.Commit(ctx)
}

func (d *Datalayer) EmitWindow(
	ctx context.Context,
	window *pb.Window,
) (pb.WindowEmitStatus, error) {
	slog.Debug("inserting window", "window", window)

	tx, err := d.WithTx(ctx)

	defer tx.Rollback(ctx)

	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return pb.WindowEmitStatus{}, err
	}

	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)

	// marshal metadata
	metadata := window.GetMetadata()
	metadataBytes, err := metadata.MarshalJSON()
	if err != nil {
		return pb.WindowEmitStatus{}, fmt.Errorf("could not marshal metadata: %v", err)
	}

	insertedWindow, err := qtx.RegisterWindow(ctx, RegisterWindowParams{
		WindowTypeName:    window.GetWindowTypeName(),
		WindowTypeVersion: window.GetWindowTypeVersion(),
		TimeFrom:          int64(window.GetTimeFrom()),
		TimeTo:            int64(window.GetTimeTo()),
		Origin:            window.GetOrigin(),
		Metadata:          metadataBytes,
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
	exec_paths, err := qtx.ReadAlgorithmExecutionPaths(
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
		}, tx.Commit(ctx)
	} else {
		return pb.WindowEmitStatus{
			Status: pb.WindowEmitStatus_NO_TRIGGERED_ALGORITHMS,
		}, nil
	}
}

func (d *Datalayer) ReadWindowTypes(
	ctx context.Context,
) (*pb.WindowTypes, error) {
	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}

	defer tx.Rollback(ctx)

	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)

	windowTypes, err := qtx.ReadWindowTypes(ctx)
	if err != nil {
		return &pb.WindowTypes{}, fmt.Errorf("could not read window types: %v", err)
	}

	windowTypesPb := pb.WindowTypes{
		Windows: make([]*pb.WindowType, len(windowTypes)),
	}

	for ii, window := range windowTypes {
		windowTypesPb.Windows[ii] = &pb.WindowType{
			Name:        window.Name,
			Version:     window.Version,
			Description: window.Description,
		}
	}
	return &windowTypesPb, tx.Commit(ctx)
}

// TODO - Add in dependencies
func (d *Datalayer) ReadAlgorithms(
	ctx context.Context,
) (*pb.Algorithms, error) {
	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)

	algorithms, err := qtx.ReadAlgorithms(ctx)
	if err != nil {
		return &pb.Algorithms{}, fmt.Errorf("could not read algorithms: %v", err)
	}

	algorithmsPb := pb.Algorithms{
		Algorithm: make([]*pb.Algorithm, len(algorithms)),
	}

	for ii, algorithm := range algorithms {
		algorithmsPb.Algorithm[ii] = &pb.Algorithm{
			Name:    algorithm.Name,
			Version: algorithm.Version,
			WindowType: &pb.WindowType{
				Name:    algorithm.WindowName,
				Version: algorithm.WindowVersion,
			},
		}
	}
	return &algorithmsPb, tx.Commit(ctx)
}

func (d *Datalayer) ReadProcessors(
	ctx context.Context,
) (*pb.Processors, error) {
	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)

	processors, err := qtx.ReadProcessors(ctx)
	if err != nil {
		return &pb.Processors{}, fmt.Errorf("could not read processors: %v", err)
	}

	processorsPb := pb.Processors{
		Processor: make([]*pb.Processors_Processor, len(processors)),
	}

	for ii, processor := range processors {
		processorsPb.Processor[ii] = &pb.Processors_Processor{
			Name:    processor.Name,
			Runtime: processor.Runtime,
		}
	}
	return &processorsPb, tx.Commit(ctx)
}

func (d *Datalayer) ReadResultsStats(
	ctx context.Context,
) (*pb.ResultsStats, error) {
	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)

	resultsStats, err := qtx.ReadResultsStats(ctx)
	if err != nil {
		return &pb.ResultsStats{}, fmt.Errorf("could not read results: %v", err)
	}

	resultsStatsPb := pb.ResultsStats{
		Count: resultsStats,
	}

	return &resultsStatsPb, tx.Commit(ctx)
}

func (d *Datalayer) ReadResultFieldsForAlgorithm(
	ctx context.Context,
	resultFieldsRead *pb.AlgorithmFieldsRead,
) (*pb.AlgorithmFields, error) {
	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)

	algorithmFields, err := qtx.ReadDistinctJsonResultFieldsForAlgorithm(
		ctx,
		ReadDistinctJsonResultFieldsForAlgorithmParams{
			TimeFrom:         resultFieldsRead.GetTimeFrom(),
			TimeTo:           resultFieldsRead.GetTimeTo(),
			AlgorithmName:    resultFieldsRead.GetAlgorithm().GetName(),
			AlgorithmVersion: resultFieldsRead.GetAlgorithm().GetVersion(),
		},
	)
	if err != nil {
		return &pb.AlgorithmFields{}, fmt.Errorf("could not read results: %v", err)
	}
	algorithmFieldsResult := pb.AlgorithmFields{
		Field: make([]string, len(algorithmFields)),
	}
	for ii, algoField := range algorithmFields {
		algorithmFieldsResult.Field[ii] = algoField
	}
	return &algorithmFieldsResult, tx.Commit(ctx)
}
