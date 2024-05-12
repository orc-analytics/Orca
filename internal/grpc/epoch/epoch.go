package grpc_epoch

import (
	"context"

	infc "github.com/predixus/pdb_framework/protobufs/go"
)

// struct level api for the GRPC service
type EpochServiceServer struct {
	infc.UnimplementedEpochServiceServer
}

// CRUD uniary methods on epochs
func (s *EpochServiceServer) RegisterEpoch(
	ctx context.Context,
	epochRequest *infc.EpochRequest,
) (*infc.EpochResponse, error) {
	return nil, nil
}

func (s *EpochServiceServer) DeleteEpoch(
	ctx context.Context,
	epochRequest *infc.EpochRequest,
) (*infc.EpochResponse, error) {
	return nil, nil
}

func (s *EpochServiceServer) ReprocessEpoch(
	ctx context.Context,
	epochRequest *infc.EpochRequest,
) (*infc.EpochResponse, error) {
	return nil, nil
}

func (s *EpochServiceServer) ModifyEpoch(
	ctx context.Context,
	epochRequest *infc.EpochRequest,
) (*infc.EpochResponse, error) {
	return nil, nil
}
