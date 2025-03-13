package datalayers

import (
	"context"

	pb "github.com/predixus/orca/protobufs/go"
)

// the interface that all datalayers must implement to be compatible with Orca
type Datalayer interface {
	CreateProcessor(ctx context.Context, proc *pb.ProcessorRegistration) error
	EmitWindow(ctx context.Context, window *pb.Window) error
}
