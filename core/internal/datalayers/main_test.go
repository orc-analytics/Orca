package datalayers

import (
	"context"
	"testing"

	pb "github.com/predixus/orca/core/protobufs/go"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestAddAlgorithm(t *testing.T) {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
	)
	connStr, err := postgresContainer.ConnectionString(ctx)

	assert.NoError(t, err)

	err = root.MigrateDatalayer("postgresql", connStr)
	assert.NoError(t, err)

	dlyr, err := NewDatalayerClient(ctx, "postgresql", connStr)
	assert.NoError(t, err)

	tx, err := dlyr.WithTx(ctx)
	defer tx.Rollback(ctx)
	assert.NoError(t, err)

	windowType := pb.WindowType{
		Name:    "TestWindow",
		Version: "1.0.0",
	}

	algo := pb.Algorithm{
		Name:       "TestAlgorithm",
		Version:    "1.0.0",
		WindowType: &windowType,
	}
	proc := pb.ProcessorRegistration{
		Name:                "TestProcessor",
		Runtime:             "Test",
		ConnectionStr:       "Test",
		SupportedAlgorithms: []*pb.Algorithm{&algo},
	}
	err = dlyr.AddAlgorithm(ctx, tx, &algo, &proc)
	assert.NoError(t, err)
}
