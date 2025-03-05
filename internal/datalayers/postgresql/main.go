package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	pb "github.com/predixus/orca/protobufs/go"
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

	procRuntime := proc.GetRuntime()
	procName := proc.GetName()
	procAlgos := proc.GetSupportedAlgorithms()

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
		ConnectionString: "", // TODO: Get connection string
	})
	if err != nil {
		slog.Error("could not create processor", "error", err)
		return err
	}

	for _, algo := range procAlgos {
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

		// create the dependencies
		dependencies := algo.GetDependencies()
		for _, algoDependentOn := range dependencies {
			err := qtx.CreateAlgorithmDependency(ctx, CreateAlgorithmDependencyParams{
				FromAlgorithmName:    algo.Name,
				FromAlgorithmVersion: algo.Version,
				FromProcessorName:    procName,
				FromProcessorRuntime: procRuntime,
				ToAlgorithmName:      algoDependentOn.GetName(),
				ToAlgorithmVersion:   algoDependentOn.GetVersion(),
				ToProcessorName:      algoDependentOn.GetProcessorName(),
				ToProcessorRuntime:   algoDependentOn.GetProcessorRuntime(),
			})
			if err != nil {
				slog.Error(
					"cloud not create algotrithm dependency",
					"algorithm",
					algo,
					"depends_on",
					algoDependentOn,
				)
			}
		}

	}

	// associate the processor with all the algos
	for _, algo := range procAlgos {
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

func (d *Datalayer) RegisterWindowType(ctx context.Context, windowType *pb.WindowType) error {
	slog.Info("registering window type", "windowType", windowType)

	err := d.queries.CreateWindowType(ctx, CreateWindowTypeParams{
		Name:    windowType.GetName(),
		Version: windowType.GetVersion(),
	})
	if err != nil {
		slog.Error("Issue registering window", "error", err.Error())
		return err
	}
	return nil
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
	return nil
}
