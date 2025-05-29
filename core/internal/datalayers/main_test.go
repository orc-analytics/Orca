package datalayers

import (
	"context"
	"os"
	"testing"

	pb "github.com/predixus/orca/core/protobufs/go"

	types "github.com/predixus/orca/core/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	testConnStr string
	testCtx     context.Context
)

func TestMain(m *testing.M) {
	var cleanup func()
	testCtx = context.Background()
	testConnStr, cleanup = setupPgOnce(testCtx)

	// Run all tests
	code := m.Run()

	// Cleanup
	cleanup()
	os.Exit(code)
}

func setupPgOnce(ctx context.Context) (string, func()) {
	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		panic("Failed to start postgres container: " + err.Error())
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic("Failed to get connection string: " + err.Error())
	}

	err = MigrateDatalayer("postgresql", connStr)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	cleanup := func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			// Log error but don't panic during cleanup
			println("Failed to terminate postgres container:", err.Error())
		}
	}

	return connStr, cleanup
}

func getCleanTx(t *testing.T, ctx context.Context) (types.Datalayer, types.Tx) {
	dlyr, err := NewDatalayerClient(ctx, "postgresql", testConnStr)
	assert.NoError(t, err)

	tx, err := dlyr.WithTx(ctx)
	assert.NoError(t, err)

	return dlyr, tx
}

func TestAddAlgorithm(t *testing.T) {
	dlyr, tx := getCleanTx(t, testCtx)
	defer tx.Rollback(testCtx)

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
	err := dlyr.CreateProcessorAndPurgeAlgos(testCtx, tx, &proc_1)
	assert.NoError(t, err)

	// 2. register the window type
	err = dlyr.CreateWindowType(testCtx, tx, &windowType)
	assert.NoError(t, err)

	// 3. add an algorithm
	err = dlyr.AddAlgorithm(testCtx, tx, &algo, &proc_1)
	assert.NoError(t, err)

	// 4. add the same algorithm again
	// adding the same algorithm again should cause no issues.
	err = dlyr.AddAlgorithm(testCtx, tx, &algo, &proc_1)
	assert.NoError(t, err)

	// 5. register a different processor
	proc_2 := pb.ProcessorRegistration{
		Name:                "TestProcessor2",
		Runtime:             "Test",
		ConnectionStr:       "Test",
		SupportedAlgorithms: []*pb.Algorithm{&algo},
	}
	err = dlyr.CreateProcessorAndPurgeAlgos(testCtx, tx, &proc_2)
	assert.NoError(t, err)

	// 6. register the same algorithm (by name and version), but with a new processor
	// should have no issues as the algorithm + processor pair is unique
	err = dlyr.AddAlgorithm(testCtx, tx, &algo, &proc_2)
	assert.NoError(t, err)
}

func TestCircularDependency(t *testing.T) {
	dlyr, tx := getCleanTx(t, testCtx)
	defer tx.Rollback(testCtx)

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
	err := dlyr.CreateProcessorAndPurgeAlgos(testCtx, tx, &proc)
	assert.NoError(t, err)

	// 2. register the window type
	err = dlyr.CreateWindowType(testCtx, tx, &windowType)
	assert.NoError(t, err)

	// 3. add an algorithm
	err = dlyr.AddAlgorithm(testCtx, tx, &algo1, &proc)
	assert.NoError(t, err)

	// 4. add another algorithm
	err = dlyr.AddAlgorithm(testCtx, tx, &algo2, &proc)
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
		testCtx,
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
		testCtx,
		tx,
		&algo2,
		&proc,
	)
	assert.ErrorIs(t, err, types.CircularDependencyFound)
}

func TestValidDependenciesBetweenProcessors(t *testing.T) {
	dlyr, tx := getCleanTx(t, testCtx)
	defer tx.Rollback(testCtx)

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
	err := dlyr.CreateProcessorAndPurgeAlgos(testCtx, tx, &proc)
	assert.NoError(t, err)

	// 2. register the window type
	err = dlyr.CreateWindowType(testCtx, tx, &windowType)
	assert.NoError(t, err)

	// 3. add algorithms
	err = dlyr.AddAlgorithm(testCtx, tx, &algo1, &proc)
	assert.NoError(t, err)
	err = dlyr.AddAlgorithm(testCtx, tx, &algo2, &proc)
	assert.NoError(t, err)
	err = dlyr.AddAlgorithm(testCtx, tx, &algo3, &proc)
	assert.NoError(t, err)
	err = dlyr.AddAlgorithm(testCtx, tx, &algo4, &proc)
	assert.NoError(t, err)

	err = dlyr.AddOverwriteAlgorithmDependency(
		testCtx,
		tx,
		&algo3,
		&proc,
	)
	assert.NoError(t, err)

	err = dlyr.AddOverwriteAlgorithmDependency(
		testCtx,
		tx,
		&algo4,
		&proc,
	)
	assert.NoError(t, err)
}
