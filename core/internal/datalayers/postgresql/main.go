package postgresql

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/orc-analytics/orca/core/internal/dag"
	pb "github.com/orc-analytics/orca/core/protobufs/go"
)

// RegisterProcessor with Orca Core
func (d *Datalayer) RegisterProcessor(
	ctx context.Context,
	proc *pb.ProcessorRegistration,
) error {
	slog.Debug("registering processor", "processor", proc)

	tx, err := d.WithTx(ctx)

	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

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
		windowType := algo.GetWindowType()

		// if there are metadata fields, add them
		var metadataFieldIds []int64
		if len(windowType.MetadataFields) > 0 {
			for _, metadataField := range windowType.MetadataFields {
				metadataFieldId, err := d.createMetadataField(ctx, tx, metadataField)
				if err != nil {
					slog.Error("could not create metadata field", "error", err)
					return err
				}
				metadataFieldIds = append(metadataFieldIds, metadataFieldId)
			}
		}

		windowTypeId, err := d.createWindowType(ctx, tx, windowType)
		if err != nil {
			slog.Error("could not create window type", "error", err)
			return err
		}

		// // remove any existing metadata field references for this window
		// err = d.flushMetadataFields(ctx, tx, windowTypeId)
		// if err != nil {
		// 	return err
		// }

		// if there were metadata fields, create the bridge
		if len(metadataFieldIds) > 0 {
			for _, metadataFieldId := range metadataFieldIds {
				err := d.createMetadataFieldBridge(ctx, tx, windowTypeId, metadataFieldId)
				if err != nil {
					slog.Error("could not create metadata field bridge", "error", err)
					return err
				}
			}
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
					algoDependentOn, "error",
					err,
				)
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

// EmitWindow with Orca core
func (d *Datalayer) EmitWindow(
	ctx context.Context,
	window *pb.Window,
) (pb.WindowEmitStatus, error) {
	slog.Debug("recieved emitted window", "window", window)

	tx, err := d.WithTx(ctx)

	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

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

	// check whether metadata is needed
	metadataFields, err := qtx.ReadMetadataFieldsByWindowType(ctx, ReadMetadataFieldsByWindowTypeParams{
		WindowTypeName:    window.GetWindowTypeName(),
		WindowTypeVersion: window.GetWindowTypeVersion(),
	})
	if err != nil {
		return pb.WindowEmitStatus{}, fmt.Errorf("could not read metadata for window: %v", err)
	}

	// confident that any required metadata is being supplied to the processor
	if len(metadataFields) > 0 {
		var metadataMap map[string]any
		if err := json.Unmarshal(metadataBytes, &metadataMap); err != nil {
			return pb.WindowEmitStatus{}, fmt.Errorf("could not unmarshal metadata for validation: %v", err)
		}

		for _, mDataField := range metadataFields {
			fieldName := mDataField.MetadataFieldName
			if _, exists := metadataMap[fieldName]; !exists {
				return pb.WindowEmitStatus{}, fmt.Errorf("required metadata field '%s' is missing", fieldName)
			}
		}
	}

	insertedWindow, err := qtx.RegisterWindow(ctx, RegisterWindowParams{
		WindowTypeName:    window.GetWindowTypeName(),
		WindowTypeVersion: window.GetWindowTypeVersion(),
		TimeFrom: pgtype.Timestamp{
			Time:  window.GetTimeFrom().AsTime().UTC(),
			Valid: true,
		},
		TimeTo: pgtype.Timestamp{
			Time:  window.GetTimeTo().AsTime().UTC(),
			Valid: true,
		},
		Origin:   window.GetOrigin(),
		Metadata: metadataBytes,
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
	execPaths, err := qtx.ReadAlgorithmExecutionPaths(
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

	slog.Info("EXEC PATHS: ", "paths", execPaths)

	// create the algo path args
	var algoIDPaths []string
	var windowTypeIDPaths []string
	var procIDPaths []string
	for _, path := range execPaths {
		algoIDPaths = append(algoIDPaths, path.AlgoIDPath)
		windowTypeIDPaths = append(windowTypeIDPaths, path.WindowTypeIDPath)
		procIDPaths = append(procIDPaths, path.ProcIDPath)
	}

	// fire off processings
	executionPlan, err := dag.BuildPlan(
		algoIDPaths,
		windowTypeIDPaths,
		procIDPaths,
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
	}
	return pb.WindowEmitStatus{
		Status: pb.WindowEmitStatus_NO_TRIGGERED_ALGORITHMS,
	}, nil
}

// ReadWindowTypes reads the types of windows registered with Orca core
func (d *Datalayer) ReadWindowTypes(
	ctx context.Context,
) (*pb.WindowTypes, error) {
	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}

	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

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

// ReadAlgorithms read the agorithms registered with Orca core
func (d *Datalayer) ReadAlgorithms(
	ctx context.Context,
) (*pb.Algorithms, error) {
	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

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
		var resultType pb.ResultType
		switch algorithm.ResultType {
		case ResultTypeNone:
			resultType = pb.ResultType_NONE
		case ResultTypeArray:
			resultType = pb.ResultType_ARRAY
		case ResultTypeStruct:
			resultType = pb.ResultType_STRUCT
		case ResultTypeValue:
			resultType = pb.ResultType_VALUE
		default:
			return &pb.Algorithms{}, fmt.Errorf(
				"read a result type that is not supported: %v",
				algorithm.ResultType,
			)
		}

		algorithmsPb.Algorithm[ii] = &pb.Algorithm{
			Name:    algorithm.Name,
			Version: algorithm.Version,
			WindowType: &pb.WindowType{
				Name:    algorithm.WindowName,
				Version: algorithm.WindowVersion,
			},
			ResultType: resultType,
		}
	}
	return &algorithmsPb, tx.Commit(ctx)
}

// ReadProcessors reads processors registered with Orca core
func (d *Datalayer) ReadProcessors(
	ctx context.Context,
) (*pb.Processors, error) {
	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

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
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

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
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)

	algorithmFields, err := qtx.ReadDistinctJsonResultFieldsForAlgorithm(
		ctx,
		ReadDistinctJsonResultFieldsForAlgorithmParams{
			TimeFrom: pgtype.Timestamp{
				Time:  resultFieldsRead.GetTimeFrom().AsTime().UTC(),
				Valid: true,
			},
			TimeTo: pgtype.Timestamp{
				Time:  resultFieldsRead.GetTimeTo().AsTime().UTC(),
				Valid: true,
			},
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
	copy(algorithmFieldsResult.Field, algorithmFields)
	return &algorithmFieldsResult, tx.Commit(ctx)
}

func (d *Datalayer) ReadResultsForAlgorithm(
	ctx context.Context,
	resultsForAlgorithmRead *pb.ResultsForAlgorithmRead,
) (*pb.ResultsForAlgorithm, error) {
	if resultsForAlgorithmRead.GetAlgorithm().GetResultType() == pb.ResultType_NONE {
		return &pb.ResultsForAlgorithm{}, fmt.Errorf(
			"cannot get results for algorithm that does not produce results",
		)
	}

	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)

	results, err := qtx.ReadResultsForAlgorithm(ctx, ReadResultsForAlgorithmParams{
		TimeFrom: pgtype.Timestamp{
			Time:  resultsForAlgorithmRead.GetTimeFrom().AsTime().UTC(),
			Valid: true,
		},
		TimeTo: pgtype.Timestamp{
			Time:  resultsForAlgorithmRead.GetTimeTo().AsTime().UTC(),
			Valid: true,
		},
		AlgorithmName:    resultsForAlgorithmRead.GetAlgorithm().GetName(),
		AlgorithmVersion: resultsForAlgorithmRead.GetAlgorithm().GetVersion(),
	})
	if err != nil {
		return &pb.ResultsForAlgorithm{}, fmt.Errorf(
			"could not read results for algorithm: %v",
			err,
		)
	}
	resultsPb := pb.ResultsForAlgorithm{
		Results: make([]*pb.ResultsForAlgorithm_ResultsRow, len(results)),
	}
	var _midpointPb *timestamppb.Timestamp
	if resultsForAlgorithmRead.GetAlgorithm().GetResultType() == pb.ResultType_VALUE {
		for ii, res := range results {
			_midpointPb = timestamppb.New(
				res.TimeFrom.Time.Add(res.TimeTo.Time.Sub(res.TimeFrom.Time) / 2),
			)

			resultsPb.Results[ii] = &pb.ResultsForAlgorithm_ResultsRow{
				Time: _midpointPb,
				ResultData: &pb.ResultsForAlgorithm_ResultsRow_SingleValue{
					SingleValue: float32(res.ResultValue.Float64),
				},
			}
		}
	} else if resultsForAlgorithmRead.GetAlgorithm().GetResultType() == pb.ResultType_ARRAY {
		for ii, res := range results {
			_midpointPb = timestamppb.New(
				res.TimeFrom.Time.Add(res.TimeTo.Time.Sub(res.TimeFrom.Time) / 2),
			)
			resultsPb.Results[ii] = &pb.ResultsForAlgorithm_ResultsRow{
				Time: _midpointPb,
				ResultData: &pb.ResultsForAlgorithm_ResultsRow_ArrayValues{
					ArrayValues: &pb.FloatArray{
						Values: convertFloat64ToFloat32(res.ResultArray),
					},
				},
			}
		}
	} else if resultsForAlgorithmRead.GetAlgorithm().GetResultType() == pb.ResultType_STRUCT {
		for ii, res := range results {
			_midpointPb = timestamppb.New(
				res.TimeFrom.Time.Add(res.TimeTo.Time.Sub(res.TimeFrom.Time) / 2),
			)
			newStruct, err := unmarshalToStruct(res.ResultJson)
			if err != nil {
				return &pb.ResultsForAlgorithm{}, fmt.Errorf("unable to parse struct data for algorithm %v: %v", resultsForAlgorithmRead.Algorithm, err)
			}

			resultsPb.Results[ii] = &pb.ResultsForAlgorithm_ResultsRow{
				Time: _midpointPb,
				ResultData: &pb.ResultsForAlgorithm_ResultsRow_StructValue{
					StructValue: newStruct,
				},
			}
		}
	} else {
		return &pb.ResultsForAlgorithm{}, fmt.Errorf("unhandled result type: %v", resultsForAlgorithmRead.GetAlgorithm().GetResultType())
	}
	return &resultsPb, tx.Commit(ctx)
}

func (d *Datalayer) ReadWindows(
	ctx context.Context,
	windowsRead *pb.WindowsRead,
) (*pb.Windows, error) {
	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)

	readWindowsRows, err := qtx.ReadWindows(ctx, ReadWindowsParams{
		TimeFrom: pgtype.Timestamp{
			Time:  windowsRead.GetTimeFrom().AsTime().UTC(),
			Valid: true,
		},
		TimeTo: pgtype.Timestamp{
			Time:  windowsRead.GetTimeTo().AsTime().UTC(),
			Valid: true,
		},
		WindowTypeName:    windowsRead.GetWindow().GetName(),
		WindowTypeVersion: windowsRead.GetWindow().GetVersion(),
	})
	if err != nil {
		return &pb.Windows{}, fmt.Errorf("could not read windows: %v", err)
	}
	windowsPb := pb.Windows{
		Window: make([]*pb.Window, len(readWindowsRows)),
	}

	var _timeFrom *timestamppb.Timestamp
	var _timeTo *timestamppb.Timestamp
	for ii, windowRow := range readWindowsRows {
		metadata, err := unmarshalToStruct(windowRow.Metadata)
		if err != nil {
			return &pb.Windows{}, fmt.Errorf("could not unpack specific window metadata: %v", err)
		}
		_timeFrom = timestamppb.New(windowRow.TimeFrom.Time)
		_timeTo = timestamppb.New(windowRow.TimeTo.Time)

		windowsPb.Window[ii] = &pb.Window{
			TimeFrom:          _timeFrom,
			TimeTo:            _timeTo,
			Origin:            windowRow.Origin,
			Metadata:          metadata,
			WindowTypeName:    windowRow.Name,
			WindowTypeVersion: windowRow.Version,
		}
	}
	return &windowsPb, tx.Commit(ctx)
}

func (d *Datalayer) ReadDistinctMetadataForWindowType(
	ctx context.Context,
	windowMetadataRead *pb.DistinctMetadataForWindowTypeRead,
) (*pb.DistinctMetadataForWindowType, error) {
	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)

	windowMetadata, err := qtx.ReadDistinctWindowMetadata(ctx, ReadDistinctWindowMetadataParams{
		TimeFrom: pgtype.Timestamp{
			Time:  windowMetadataRead.GetTimeFrom().AsTime(),
			Valid: true,
		},
		TimeTo: pgtype.Timestamp{
			Time:  windowMetadataRead.GetTimeTo().AsTime(),
			Valid: true,
		},
		WindowTypeName:    windowMetadataRead.GetWindowType().GetName(),
		WindowTypeVersion: windowMetadataRead.GetWindowType().GetVersion(),
	})
	if err != nil {
		return nil, fmt.Errorf("could not read window metadata: %w", err)
	}

	// Convert raw JSON blobs into []any of map[string]any
	metadataList := make([]any, len(windowMetadata))
	for i, raw := range windowMetadata {
		var m map[string]any
		if err := json.Unmarshal(raw, &m); err != nil {
			return nil, fmt.Errorf("could not unmarshal metadata at index %d: %w", i, err)
		}
		metadataList[i] = m
	}

	// Build a protobuf ListValue directly from []any
	listValue, err := structpb.NewList(metadataList)
	if err != nil {
		return nil, fmt.Errorf("could not convert metadata list to protobuf ListValue: %w", err)
	}

	return &pb.DistinctMetadataForWindowType{
		Metadata: listValue,
	}, tx.Commit(ctx)
}

func (d *Datalayer) ReadWindowsForMetadata(
	ctx context.Context,
	windowsForMetadataRead *pb.WindowsForMetadataRead,
) (*pb.WindowsForMetadata, error) {
	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}

	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)

	metadataFields := windowsForMetadataRead.GetMetadata()
	metadata := make(map[string]any, len(metadataFields))

	for _, m := range metadataFields {
		switch m.GetValue().GetKind().(type) {
		case *structpb.Value_BoolValue:
			metadata[m.Field] = m.Value.GetBoolValue()
		case *structpb.Value_NumberValue:
			metadata[m.Field] = m.Value.GetNumberValue()
		case *structpb.Value_StringValue:
			metadata[m.Field] = m.Value.GetStringValue()
		case *structpb.Value_NullValue:
			metadata[m.Field] = nil
		case *structpb.Value_ListValue:
			metadata[m.Field] = m.Value.GetListValue()
		case *structpb.Value_StructValue:
			metadata[m.Field] = m.Value.GetStructValue()
		}
	}
	metadataJson, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("could not parse metadata as json: %v", err)
	}

	rows, err := qtx.ReadWindowsForMetadata(ctx, ReadWindowsForMetadataParams{
		WindowTypeName:    windowsForMetadataRead.GetWindow().GetName(),
		WindowTypeVersion: windowsForMetadataRead.GetWindow().GetVersion(),
		TimeFrom: pgtype.Timestamp{
			Time:  windowsForMetadataRead.GetTimeFrom().AsTime(),
			Valid: true,
		},
		TimeTo: pgtype.Timestamp{
			Time:  windowsForMetadataRead.GetTimeTo().AsTime(),
			Valid: true,
		},
		MetadataFilter: metadataJson,
	})
	if err != nil {
		return nil, fmt.Errorf("could not read rows: %v", err)
	}
	result := pb.WindowsForMetadata{
		Window: make([]*pb.Window, len(rows)),
	}
	for ii, row := range rows {
		metadataStructPb, err := unmarshalToStruct(row.Metadata)
		if err != nil {
			return nil, fmt.Errorf(
				"could not unmarshal postgres metadata jsonb to structpb: %v",
				err,
			)
		}
		result.Window[ii] = &pb.Window{
			TimeFrom:          timestamppb.New(row.TimeFrom.Time),
			TimeTo:            timestamppb.New(row.TimeTo.Time),
			Origin:            row.Origin,
			WindowTypeName:    row.Name,
			WindowTypeVersion: row.Version,
			Metadata:          metadataStructPb,
		}
	}

	return &result, tx.Commit(ctx)
}

func (d *Datalayer) ReadResultsForAlgorithmAndMetadata(
	ctx context.Context,
	resultsForAlgorithmAndMetadata *pb.ResultsForAlgorithmAndMetadataRead,
) (*pb.ResultsForAlgorithmAndMetadata, error) {
	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}

	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)

	metadataFields := resultsForAlgorithmAndMetadata.GetMetadata()
	metadata := make(map[string]any, len(metadataFields))

	for _, m := range metadataFields {
		switch m.GetValue().GetKind().(type) {
		case *structpb.Value_BoolValue:
			metadata[m.Field] = m.Value.GetBoolValue()
		case *structpb.Value_NumberValue:
			metadata[m.Field] = m.Value.GetNumberValue()
		case *structpb.Value_StringValue:
			metadata[m.Field] = m.Value.GetStringValue()
		case *structpb.Value_NullValue:
			metadata[m.Field] = nil
		case *structpb.Value_ListValue:
			metadata[m.Field] = m.Value.GetListValue()
		case *structpb.Value_StructValue:
			metadata[m.Field] = m.Value.GetStructValue()
		}
	}
	metadataJson, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("could not parse metadata as json: %v", err)
	}

	rows, err := qtx.ReadResultsForAlgorithmAndMetadata(
		ctx,
		ReadResultsForAlgorithmAndMetadataParams{
			TimeFrom: pgtype.Timestamp{
				Time:  resultsForAlgorithmAndMetadata.GetTimeFrom().AsTime(),
				Valid: true,
			},
			TimeTo: pgtype.Timestamp{
				Time:  resultsForAlgorithmAndMetadata.GetTimeTo().AsTime(),
				Valid: true,
			},
			MetadataFilter:   metadataJson,
			AlgorithmName:    resultsForAlgorithmAndMetadata.GetAlgorithm().GetName(),
			AlgorithmVersion: resultsForAlgorithmAndMetadata.GetAlgorithm().GetVersion(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not read results for algorithm: %v",
			err,
		)
	}

	resultsPb := pb.ResultsForAlgorithmAndMetadata{
		Results: make([]*pb.ResultsForAlgorithmAndMetadata_ResultsRow, len(rows)),
	}
	var _midpointPb *timestamppb.Timestamp
	if resultsForAlgorithmAndMetadata.GetAlgorithm().GetResultType() == pb.ResultType_VALUE {
		for ii, res := range rows {
			_midpointPb = timestamppb.New(
				res.TimeFrom.Time.Add(res.TimeTo.Time.Sub(res.TimeFrom.Time) / 2),
			)

			resultsPb.Results[ii] = &pb.ResultsForAlgorithmAndMetadata_ResultsRow{
				Time: _midpointPb,
				ResultData: &pb.ResultsForAlgorithmAndMetadata_ResultsRow_SingleValue{
					SingleValue: float32(res.ResultValue.Float64),
				},
			}
		}
	} else if resultsForAlgorithmAndMetadata.GetAlgorithm().GetResultType() == pb.ResultType_ARRAY {
		for ii, res := range rows {
			_midpointPb = timestamppb.New(
				res.TimeFrom.Time.Add(res.TimeTo.Time.Sub(res.TimeFrom.Time) / 2),
			)
			resultsPb.Results[ii] = &pb.ResultsForAlgorithmAndMetadata_ResultsRow{
				Time: _midpointPb,
				ResultData: &pb.ResultsForAlgorithmAndMetadata_ResultsRow_ArrayValues{
					ArrayValues: &pb.FloatArray{
						Values: convertFloat64ToFloat32(res.ResultArray),
					},
				},
			}
		}
	} else if resultsForAlgorithmAndMetadata.GetAlgorithm().GetResultType() == pb.ResultType_STRUCT {
		for ii, res := range rows {
			_midpointPb = timestamppb.New(
				res.TimeFrom.Time.Add(res.TimeTo.Time.Sub(res.TimeFrom.Time) / 2),
			)
			newStruct, err := unmarshalToStruct(res.ResultJson)
			if err != nil {
				return nil, fmt.Errorf("unable to parse struct data for algorithm %v: %v", resultsForAlgorithmAndMetadata.GetAlgorithm(), err)
			}

			resultsPb.Results[ii] = &pb.ResultsForAlgorithmAndMetadata_ResultsRow{
				Time: _midpointPb,
				ResultData: &pb.ResultsForAlgorithmAndMetadata_ResultsRow_StructValue{
					StructValue: newStruct,
				},
			}
		}
	} else {
		return nil, fmt.Errorf("unhandled result type: %v", resultsForAlgorithmAndMetadata.GetAlgorithm().GetResultType())
	}
	return &resultsPb, tx.Commit(ctx)
}

// Annotate a section of time
func (d *Datalayer) Annotate(
	ctx context.Context,
	annotateWrite *pb.AnnotateWrite,
) (*pb.AnnotateResponse, error) {
	tx, err := d.WithTx(ctx)
	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return nil, err
	}

	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	pgTx := tx.(*PgTx)
	qtx := d.queries.WithTx(pgTx.tx)

	// insert annotation
	metadata := annotateWrite.GetMetadata()
	metadataJson, err := metadata.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not marshal provided metadata to json: %v", err)
	}
	annotationId, err := qtx.CreateAnnotation(ctx, CreateAnnotationParams{
		TimeFrom: pgtype.Timestamp{
			Time:  annotateWrite.GetTimeFrom().AsTime(),
			Valid: true,
		},
		TimeTo: pgtype.Timestamp{
			Time:  annotateWrite.GetTimeTo().AsTime(),
			Valid: true,
		},
		Description: pgtype.Text{
			String: annotateWrite.GetDescription(),
			Valid:  true,
		},
		Metadata: metadataJson,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create annotation: %v", err)
	}

	// link annotation to window type
	for _, capturedWindow := range annotateWrite.GetCapturedWindows() {

		err := qtx.LinkAnnotationToWindowType(ctx, LinkAnnotationToWindowTypeParams{
			AnnotationID:  annotationId,
			WindowName:    capturedWindow.GetName(),
			WindowVersion: capturedWindow.GetVersion(),
		})
		if err != nil {
			return nil, fmt.Errorf("could not link annotation to window: %v", err)
		}
	}

	// link annotation to algorithm type
	for _, capturedAlgorithm := range annotateWrite.GetCapturedAlgorithms() {
		err := qtx.LinkAnnotationToAlgorithm(ctx, LinkAnnotationToAlgorithmParams{
			AnnotationID:     annotationId,
			AlgorithmName:    capturedAlgorithm.GetName(),
			AlgorithmVersion: capturedAlgorithm.GetVersion(),
		})
		if err != nil {
			return nil, fmt.Errorf("could not link annotation to algorithm: %v", err)
		}
	}

	return &pb.AnnotateResponse{}, tx.Commit(ctx)
}
