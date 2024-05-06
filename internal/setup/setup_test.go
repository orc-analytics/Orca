package setup_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	setup "github.com/predixus/analytics_framework/internal/setup"
)

type MockDB struct{}

func (d *MockDB) Connect() error { return nil }
func (d *MockDB) Close() error   { return nil }

type MockAPI struct{}

func (d *MockAPI) Start(wg *sync.WaitGroup) {
	wg.Done()
}

type MockGrpc struct{}

func (d *MockGrpc) Start(wg *sync.WaitGroup) {
	wg.Done()
}

func TestSetupCompletes(t *testing.T) {
	db := MockDB{}
	api := MockAPI{}
	grpc := MockGrpc{}

	err := setup.Setup(&db, &grpc, &api)
	assert.Nil(t, err, "No error returned on setup")
	assert.True(t, true, "Setup completed successfully")
}
