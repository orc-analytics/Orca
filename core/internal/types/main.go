package datalayers

import (
	"context"
	"fmt"

	pb "github.com/orc-analytics/orca/core/protobufs/go"
)

// the interface that all datalayers must implement to be compatible with Orca
type (
	Tx interface {
		Rollback(ctx context.Context)
		Commit(ctx context.Context) error
	}
	Datalayer interface {
		WithTx(ctx context.Context) (Tx, error)

		// Core level operations
		RegisterProcessor(ctx context.Context, proc *pb.ProcessorRegistration) error
		EmitWindow(ctx context.Context, window *pb.Window) (pb.WindowEmitStatus, error)

		// Data level operations
		ReadWindowTypes(ctx context.Context) (*pb.WindowTypes, error)
		ReadAlgorithms(ctx context.Context) (*pb.Algorithms, error)
		ReadProcessors(ctx context.Context) (*pb.Processors, error)
		ReadResultsStats(ctx context.Context) (*pb.ResultsStats, error)
		ReadResultFieldsForAlgorithm(
			ctx context.Context,
			algorithmFieldsRead *pb.AlgorithmFieldsRead,
		) (*pb.AlgorithmFields, error)
		ReadResultsForAlgorithm(
			ctx context.Context,
			resultsForAlgorithm *pb.ResultsForAlgorithmRead,
		) (*pb.ResultsForAlgorithm, error)
		ReadResultsForAlgorithmAndMetadata(
			ctx context.Context,
			resultsForAlgorithmAndMetadata *pb.ResultsForAlgorithmAndMetadataRead,
		) (*pb.ResultsForAlgorithmAndMetadata, error)
		ReadWindows(ctx context.Context, windowsRead *pb.WindowsRead) (*pb.Windows, error)
		ReadDistinctMetadataForWindowType(
			ctx context.Context,
			windowMetadataRead *pb.DistinctMetadataForWindowTypeRead,
		) (*pb.DistinctMetadataForWindowType, error)
		ReadWindowsForMetadata(
			ctx context.Context,
			windowsForMetadataRead *pb.WindowsForMetadataRead,
		) (*pb.WindowsForMetadata, error)
		Annotate(ctx context.Context, annotateWrite *pb.AnnotateWrite) (*pb.AnnotateResponse, error)
	}
)

// custom errors
var (
	AlgorithmExistsUnderDifferentProcessor = fmt.Errorf(
		"algorithm exists under a different processor",
	)
)

type CircularDependencyError struct {
	FromAlgoName      string
	ToAlgoName        string
	FromAlgoVersion   string
	ToAlgoVersion     string
	FromAlgoProcessor string
	ToAlgoProcessor   string
}

func (c *CircularDependencyError) Error() string {
	return fmt.Sprintf(
		"Circular dependency introduced between algortithm %s to %s, with versions %s and %s, of processor(s) %s and %s respectively.",
		c.FromAlgoName,
		c.ToAlgoName,
		c.FromAlgoVersion,
		c.ToAlgoVersion,
		c.FromAlgoProcessor,
		c.ToAlgoProcessor,
	)
}
