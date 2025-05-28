package datalayers

import (
	"context"
	"os"
	"testing"

	pb "github.com/predixus/orca/core/protobufs/go"
	"github.com/stretchr/testify/assert"
)

func TestAddAlgorithm(t *testing.T) {
	ctx := context.Background()
	connStr := os.Getenv("ORCA_DATABASE_URL")
	if connStr == "" {
		t.Error("could not find `ORCA_DATABASE_URL` env var")
		return
	}

	dlyr, err := NewClient(ctx, connStr)
	if err != nil {
		t.Error(err)
		return
	}
	tx, err := dlyr.WithTx(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	windowType := pb.WindowType{
		Name:    "TestWindow",
		Version: "1.0.0",
	}

	algo := pb.Algorithm{
		Name:       "TestAlgorithm",
		Version:    "1.0.0",
		WindowType: &windowType,
	}
	proc := pb.ProcessorRegistration{
		Name:                "TestProcessor",
		Runtime:             "Test",
		ConnectionStr:       "Test",
		SupportedAlgorithms: []*pb.Algorithm{&algo},
	}
	err = dlyr.AddAlgorithm(ctx, tx, &algo, &proc)
	assert.NoError(t, err)
}
