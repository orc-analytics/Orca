package datalayers

import (
	"context"
	"fmt"

	pb "github.com/predixus/orca/core/protobufs/go"
)

// the interface that all datalayers must implement to be compatible with Orca
type (
	Tx interface {
		Rollback(ctx context.Context)
		Commit(ctx context.Context) error
	}
	Datalayer interface {
		WithTx(ctx context.Context) (Tx, error)
		CreateProcessorAndPurgeAlgos(
			ctx context.Context,
			tx Tx,
			proc *pb.ProcessorRegistration,
		) error
		CreateWindowType(ctx context.Context, qtx Tx, windowType *pb.WindowType) error
		AddAlgorithm(
			ctx context.Context,
			tx Tx,
			algo *pb.Algorithm,
			proc *pb.ProcessorRegistration,
		) error
		AddOverwriteAlgorithmDependency(
			ctx context.Context,
			tx Tx,
			algo *pb.Algorithm,
			proc *pb.ProcessorRegistration,
		) error
		AddOverwriteDataGetter(
			ctx context.Context,
			tx Tx,
			proc *pb.ProcessorRegistration,
		) error
		EmitWindow(ctx context.Context, window *pb.Window) (pb.WindowEmitStatus, error)
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
