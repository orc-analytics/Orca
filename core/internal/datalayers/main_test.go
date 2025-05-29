package datalayers

import (
	"context"
	"testing"

	pb "github.com/predixus/orca/core/protobufs/go"

	types "github.com/predixus/orca/core/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupPg(t *testing.T, ctx context.Context) string {
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
	return connStr
}

func TestAddAlgorithm(t *testing.T) {
	ctx := context.Background()

	// TODO - parametrise over datalayers
	connStr := setupPg(t, ctx)

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
	// should have no issues as the algorithm + processor pair is unique
	err = dlyr.AddAlgorithm(ctx, tx, &algo, &proc_2)
	assert.NoError(t, err)
}

func TestCircularDependency(t *testing.T) {
	ctx := context.Background()

	// TODO - parametrise over datalayers
	connStr := setupPg(t, ctx)

	dlyr, err := NewDatalayerClient(ctx, "postgresql", connStr)
	assert.NoError(t, err)

	tx, err := dlyr.WithTx(ctx)
	defer tx.Rollback(ctx)
	assert.NoError(t, err)

	windowType := pb.WindowType{
		Name:    "TestWindow",
		Version: "1.0.0",
	}

	algo1 := pb.Algorithm{
		Name:       "TestAlgorithm1",
		Version:    "1.0.0",
		WindowType: &windowType,
	}

	algo2 := pb.Algorithm{
		Name:       "TestAlgorithm2",
		Version:    "1.0.0",
		WindowType: &windowType,
	}

	proc := pb.ProcessorRegistration{
		Name:                "TestProcessor",
		Runtime:             "Test",
		ConnectionStr:       "Test",
		SupportedAlgorithms: []*pb.Algorithm{&algo1, &algo2},
	}

	// 1. register a processor
	err = dlyr.CreateProcessorAndPurgeAlgos(ctx, tx, &proc)
	assert.NoError(t, err)

	// 2. register the window type
	err = dlyr.CreateWindowType(ctx, tx, &windowType)
	assert.NoError(t, err)

	// 3. add an algorithm
	err = dlyr.AddAlgorithm(ctx, tx, &algo1, &proc)
	assert.NoError(t, err)

	// 4. add another algorithm
	err = dlyr.AddAlgorithm(ctx, tx, &algo2, &proc)
	assert.NoError(t, err)

	// so far so good

	// 5. add a dependency between algorithm 1 and algorithm 2
	algo1.Dependencies = []*pb.AlgorithmDependency{
		{
			Name:             "TestAlgorithm2",
			Version:          "1.0.0",
			ProcessorName:    "TestProcessor",
			ProcessorRuntime: "Test",
		},
	}

	err = dlyr.AddOverwriteAlgorithmDependency(
		ctx,
		tx,
		&algo1,
		&proc,
	)
	assert.NoError(t, err)

	// 6. now add a dependency between 2 and 1. This should raise a circular error
	algo2.Dependencies = []*pb.AlgorithmDependency{
		{
			Name:             "TestAlgorithm1",
			Version:          "1.0.0",
			ProcessorName:    "TestProcessor",
			ProcessorRuntime: "Test",
		},
	}

	err = dlyr.AddOverwriteAlgorithmDependency(
		ctx,
		tx,
		&algo2,
		&proc,
	)
	assert.ErrorIs(t, err, types.CircularDependencyFound)
}

func TestValidDependenciesBetweenProcessors(t *testing.T) {
	ctx := context.Background()

	// TODO - parametrise over datalayers
	connStr := setupPg(t, ctx)
	dlyr, err := NewDatalayerClient(ctx, "postgresql", connStr)

	assert.NoError(t, err)

	tx, err := dlyr.WithTx(ctx)
	defer tx.Rollback(ctx)
	assert.NoError(t, err)

	windowType := pb.WindowType{
		Name:    "TestWindow",
		Version: "1.0.0",
	}

	algo1 := pb.Algorithm{
		Name:       "TestAlgorithm1",
		Version:    "1.0.0",
		WindowType: &windowType,
	}
	algo2 := pb.Algorithm{
		Name:       "TestAlgorithm2",
		Version:    "1.0.0",
		WindowType: &windowType,
	}

	algo3 := pb.Algorithm{
		Name:       "TestAlgorithm3",
		Version:    "1.0.0",
		WindowType: &windowType,
	}

	algo4 := pb.Algorithm{
		Name:       "TestAlgorithm4",
		Version:    "1.0.0",
		WindowType: &windowType,
	}

	proc := pb.ProcessorRegistration{
		Name:                "TestProcessor",
		Runtime:             "Test",
		ConnectionStr:       "Test",
		SupportedAlgorithms: []*pb.Algorithm{&algo1, &algo2},
	}

	algo3.Dependencies = []*pb.AlgorithmDependency{
		{
			Name:             algo1.Name,
			Version:          algo1.Version,
			ProcessorName:    proc.GetName(),
			ProcessorRuntime: proc.GetRuntime(),
		},
		{
			Name:             algo2.Name,
			Version:          algo2.Version,
			ProcessorName:    proc.GetName(),
			ProcessorRuntime: proc.GetRuntime(),
		},
	}

	algo4.Dependencies = []*pb.AlgorithmDependency{
		{
			Name:             algo3.Name,
			Version:          algo3.Version,
			ProcessorName:    proc.GetName(),
			ProcessorRuntime: proc.GetRuntime(),
		},
	}

	// 1. register a processor
	err = dlyr.CreateProcessorAndPurgeAlgos(ctx, tx, &proc)
	assert.NoError(t, err)

	// 2. register the window type
	err = dlyr.CreateWindowType(ctx, tx, &windowType)
	assert.NoError(t, err)

	// 3. add algorithms
	err = dlyr.AddAlgorithm(ctx, tx, &algo1, &proc)
	assert.NoError(t, err)
	err = dlyr.AddAlgorithm(ctx, tx, &algo2, &proc)
	assert.NoError(t, err)
	err = dlyr.AddAlgorithm(ctx, tx, &algo3, &proc)
	assert.NoError(t, err)
	err = dlyr.AddAlgorithm(ctx, tx, &algo4, &proc)
	assert.NoError(t, err)

	err = dlyr.AddOverwriteAlgorithmDependency(
		ctx,
		tx,
		&algo3,
		&proc,
	)
	assert.NoError(t, err)

	err = dlyr.AddOverwriteAlgorithmDependency(
		ctx,
		tx,
		&algo4,
		&proc,
	)
	assert.NoError(t, err)
}
