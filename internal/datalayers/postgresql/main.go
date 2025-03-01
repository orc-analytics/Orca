package postgresql

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	inte "github.com/predixus/orca/internal"
	pb "github.com/predixus/orca/protobufs/go"
)

type datalayer struct {
	queries *Queries
	closeFn func(context.Context) error
}

// generate a new client for the postgres datalayer
func NewClient(ctx context.Context, connStr *string) (inte.Datalayer, error) {
	var psqlStr string

	if connStr != nil {
		psqlStr = *connStr
	} else { // sensible default for local dev
		psqlStr = "postgresql://orca:orca@localhost:5432/orca?sslmode=verify-full"
	}

	conn, err := pgx.Connect(ctx, psqlStr)
	if err != nil {
		slog.Error("Issue connecting to postgres", "error", err)
		return nil, err
	}

	queries := New(conn)
	return &datalayer{
		queries: queries,
		closeFn: conn.Close,
	}, nil
}

// AddProcessor add a processor to the Orca server
func (d *datalayer) AddProcessor(ctx context.Context, proc *pb.ProcessorRegistration) error {
	runtime := proc.GetRuntime()
	name := proc.GetName()

	if runtime == "" {
		msg := "Runtime not provided"
		slog.Error(msg, "protobuf", proc)
		return fmt.Errorf(msg)
	}

	if name == "" {
		msg := "Processor name not provided"
		slog.Error("msg", "protobuf", proc)
		return fmt.Errorf(msg)
	}

	d.queries.AddProcessor(ctx, AddProcessorParams{
		Name:    name,
		Runtime: runtime,
	})
	return nil
}
