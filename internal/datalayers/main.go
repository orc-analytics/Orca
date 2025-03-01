package datalayers

import (
	"context"
	"fmt"
	"log/slog"

	inte "github.com/predixus/orca/internal"
	psql "github.com/predixus/orca/internal/datalayers/postgresql"
)

// represents the supported database platforms as the datalayer
type Platform string

const (
	PostgreSQL Platform = "postgresql"
)

// check if the platform is supported
func (p Platform) IsValid() bool {
	switch p {
	case PostgreSQL:
		return true
	default:
		return false
	}
}

func NewDatalayerClient(
	ctx context.Context,
	platform Platform,
	connStr string,
) (inte.Datalayer, error) {
	if !platform.IsValid() {
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}

	switch platform {
	case PostgreSQL:
		return psql.NewClient(ctx, &connStr)
	default:
		slog.Error("attempted to access unsuported platform", "platform", platform)
		return nil, fmt.Errorf("platform not implemented: %s", platform)
	}
}
