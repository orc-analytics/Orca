package internal

import (
	"context"

	pb "github.com/predixus/orca/protobufs/go"
)

// the interface that all datalayers must implement to be compatible
type Datalayer interface {
	AddProcessor(ctx context.Context, proc *pb.ProcessorRegistration) error
}
