package internal

import (
	pb "github.com/predixus/orca/protobufs/go"
)

// the interface that all datalayers must implement to be compatible
type Datalayer interface {
	AddProcessor(pb.ProcessorRegistration) error
}
