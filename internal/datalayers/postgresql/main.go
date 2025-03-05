package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	pb "github.com/predixus/orca/protobufs/go"
)

type Datalayer struct {
	queries *Queries
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

	queries := New(conn)
	return &Datalayer{
		queries: queries,
		closeFn: conn.Close,
	}, nil
}

// AddProcessor add a processor to the Orca server
func (d *Datalayer) AddProcessor(ctx context.Context, proc *pb.ProcessorRegistration) error {
	runtime := proc.GetRuntime()
	name := proc.GetName()

	if runtime == "" {
		msg := "Runtime not provided"
		slog.Error(msg, "protobuf", proc)
		return fmt.Errorf(msg)
	}

	if name == "" {
		msg := "Processor name not provided"
		slog.Error(msg, "protobuf", proc)
		return fmt.Errorf(msg)
	}

	d.queries.AddProcessor(ctx, AddProcessorParams{
		Name:    name,
		Runtime: runtime,
	})
	return nil
}

func (d *Datalayer) RegisterWindow(ctx context.Context, window *pb.Window) error {
	from := window.GetFrom()
	to := window.GetTo()
	windowName := window.GetName()
	origin := window.GetOrigin()

	if from == 0 || to == 0 {
		msg := "`from` or `to` can't be equal to 0"
		slog.Error(msg, "from", from, "to", to)
		return fmt.Errorf(msg)
	}

	slog.Info("Registering window", "from", from, "to", to, "window name", windowName)
	_, err := d.queries.RegisterWindow(ctx, RegisterWindowParams{
		WindowName: windowName,
		TimeFrom:   int64(from),
		TimeTo:     int64(to),
		Origin:     origin,
	})
	if err != nil {
		slog.Error("Issue registering window", "error", err.Error())
		return err
	}
	return nil
}
