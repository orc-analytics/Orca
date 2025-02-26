package postgresql

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	inte "github.com/predixus/orca/internal"
)

type datalayer struct{}

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
	defer conn.Close(ctx)

	return nil, nil
}
