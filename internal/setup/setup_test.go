package setup_test

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	setup "github.com/predixus/analytics_framework/internal/setup"
)

type MockDB struct{}

func (d *MockDB) Connect() *sql.DB { return nil }
func (d *MockDB) Close() error     { return nil }

type MockAPI struct{}

func (d *MockAPI) Start(wg *sync.WaitGroup) {}

type MockGrpc struct{}

func (d *MockGrpc) Start(wg *sync.WaitGroup) {}

func TestSetupCompletes(t *testing.T) {
	db := MockDB{}
	api := MockAPI{}
	grpc := MockGrpc{}

	setup.Setup(&db, &grpc, &api)
	assert.True(t, true, "Setup completed successfully")
}
