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
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)

	assert.NoError(t, err)
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")

	assert.NoError(t, err)

	err = MigrateDatalayer("postgresql", connStr)

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
	proc_1 := pb.ProcessorRegistration{
		Name:                "TestProcessor",
		Runtime:             "Test",
		ConnectionStr:       "Test",
		SupportedAlgorithms: []*pb.Algorithm{&algo},
	}

	// 1. register a processor
	err = dlyr.CreateProcessorAndPurgeAlgos(ctx, tx, &proc_1)
	assert.NoError(t, err)

	// 2. register the window type
	err = dlyr.CreateWindowType(ctx, tx, &windowType)
	assert.NoError(t, err)

	// 3. add an algorithm
	err = dlyr.AddAlgorithm(ctx, tx, &algo, &proc_1)
	assert.NoError(t, err)

	// 4. add the same algorithm again
	// adding the same algorithm again should cause no issues.
	err = dlyr.AddAlgorithm(ctx, tx, &algo, &proc_1)
	assert.NoError(t, err)

	// 5. register a different processor
	proc_2 := pb.ProcessorRegistration{
		Name:                "TestProcessor2",
		Runtime:             "Test",
		ConnectionStr:       "Test",
		SupportedAlgorithms: []*pb.Algorithm{&algo},
	}
	err = dlyr.CreateProcessorAndPurgeAlgos(ctx, tx, &proc_2)
	assert.NoError(t, err)

	// 6. register the same algorithm (by name and version), but with a new processor
	// should have 0 issues as the algorithm + processor pair is unique
	err = dlyr.AddAlgorithm(ctx, tx, &algo, &proc_2)
	assert.NoError(t, err)
}
