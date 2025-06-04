package datalayers

import (
	"context"
	"fmt"
	"log/slog"

	psql "github.com/predixus/orca/core/internal/datalayers/postgresql"
	types "github.com/predixus/orca/core/internal/types"
)

// represents the supported database platforms as the datalayer
type Platform string

const (
	PostgreSQL Platform = "postgresql"
)

// check if the platform is supported
func (p Platform) isValid() bool {
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
) (types.Datalayer, error) {
	if !platform.isValid() {
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}

	switch platform {
	case PostgreSQL:
		return psql.NewClient(ctx, connStr)
	default:
		slog.Error("datalayer not supported", "platform", platform)
		return nil, fmt.Errorf("platform not implemented: %s", platform)
	}
}
