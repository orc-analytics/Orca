package datalayers

import (
	"context"
	"fmt"
	"log/slog"

	psql "github.com/predixus/orca/core/internal/datalayers/postgresql"
	types "github.com/predixus/orca/core/internal/types"
	pb "github.com/predixus/orca/core/protobufs/go"
)

// Platform represents the supported database platforms as the datalayer
type Platform string

// DataLayerConstructor is a function type for creating datalayer clients
type DataLayerConstructor func(ctx context.Context, connStr string) (types.Datalayer, error)

// PlatformInfo holds information about a platform
type PlatformInfo struct {
	Name        Platform
	Constructor DataLayerConstructor
}

// Registry holds all registered platforms
var platformRegistry = make(map[Platform]PlatformInfo)

// Constants for supported platforms
const (
	PostgreSQL Platform = "postgresql"
	// TODO
	// MongoDB  Platform = "mongodb"
	// MySQL    Platform = "mysql"
)

// init registers all available platforms
func init() {
	registerPlatforms()
}

// registerPlatforms adds all available platforms to the registry
func registerPlatforms() {
	platforms := []PlatformInfo{
		{
			Name: PostgreSQL,
			Constructor: func(ctx context.Context, connStr string) (types.Datalayer, error) {
				return psql.NewClient(ctx, connStr)
			},
		},
	}

	for _, platform := range platforms {
		platformRegistry[platform.Name] = platform
	}
}

// GetSupportedPlatforms returns a slice of all supported platforms
func GetSupportedPlatforms() []Platform {
	platforms := make([]Platform, 0, len(platformRegistry))
	for platform := range platformRegistry {
		platforms = append(platforms, platform)
	}
	return platforms
}

// check if the platform is supported
func (p Platform) isValid() bool {
	_, exists := platformRegistry[p]
	return exists
}

// NewDatalayerClient creates a new datalayer client for the specified platform
func NewDatalayerClient(
	ctx context.Context,
	platform Platform,
	connStr string,
) (types.Datalayer, error) {
	platformInfo, exists := platformRegistry[platform]
	if !exists {
		return nil, fmt.Errorf("unsupported platform: %s. Supported platforms: %v",
			platform, GetSupportedPlatforms())
	}

	return platformInfo.Constructor(ctx, connStr)
}

// RegisterProcessor creates and registers a processor with its algorithms
func RegisterProcessor(
	ctx context.Context,
	dlyr types.Datalayer,
	proc *pb.ProcessorRegistration,
) error {
	slog.Debug("creating processor", "protobuf", proc)
	tx, err := dlyr.WithTx(ctx)

	defer tx.Rollback(ctx)

	if err != nil {
		slog.Error("could not start a transaction", "error", err)
		return err
	}

	// register the processor
	err = dlyr.CreateProcessorAndPurgeAlgos(ctx, tx, proc)
	if err != nil {
		slog.Error("could not create processor", "error", err)
		return err
	}

	// add all algorithms first
	for _, algo := range proc.GetSupportedAlgorithms() {
		// add window types
		window_type := algo.GetWindowType()

		err := dlyr.CreateWindowType(ctx, tx, window_type)
		if err != nil {
			slog.Error("could not create window type", "error", err)
			return err
		}

		// create algos
		err = dlyr.AddAlgorithm(ctx, tx, algo, proc)
		if err != nil {
			slog.Error("error creating algorithm", "error", err)
			return err
		}
	}

	// then add the dependencies and associate the processor with all the algos
	for _, algo := range proc.GetSupportedAlgorithms() {
		dependencies := algo.GetDependencies()
		for _, algoDependentOn := range dependencies {
			err := dlyr.AddOverwriteAlgorithmDependency(
				ctx,
				tx,
				algo,
				proc,
			)
			if err != nil {
				slog.Error(
					"could not create algorithm dependency",
					"algorithm",
					algo,
					"depends_on",
					algoDependentOn,
					"error",
					err,
				)
				return err
			}
		}
	}

	return tx.Commit(ctx)
}
