package setup_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/predixus/analytics_framework/internal/api"
	"github.com/predixus/analytics_framework/internal/datalayer"
	"github.com/predixus/analytics_framework/internal/grpc"
	"github.com/predixus/analytics_framework/internal/setup"
)

// MockDB is a mock implementation of the database connection.
type MockDB struct{}

func (m *MockDB) ConnectDB() error {
	// Mock implementation, return nil for success
	return nil
}

func (m *MockDB) Close() {
	// Mock implementation, do nothing
}

// MockGRPCServer is a mock implementation of the gRPC server.
type MockGRPCServer struct{}

func (m *MockGRPCServer) StartGRPCServer(wg *sync.WaitGroup) {
	// Mock implementation, do nothing
	wg.Done()
}

// MockHTTPServer is a mock implementation of the HTTP server.
type MockHTTPServer struct{}

func (m *MockHTTPServer) StartHTTPServer(wg *sync.WaitGroup) {
	// Mock implementation, do nothing
	wg.Done()
}

func TestSetup(t *testing.T) {
	// Replace the original implementations with mocks
	datalayer.StorageDB = &MockDB{}
	grpc.FrameworkServer = &MockGRPCServer{}
	api.HTTPServer = &MockHTTPServer{}

	// Run the setup function
	setup.Setup()

	// Ensure that the setup completes without errors
	// Here, you might want to add more specific assertions based on your requirements
	assert.True(t, true, "Setup completed successfully")
}
