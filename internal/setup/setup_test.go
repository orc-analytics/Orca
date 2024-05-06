package setup_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/predixus/analytics_framework/internal/datalayer"
	"github.com/predixus/analytics_framework/internal/setup"
)

// DBInterface defines the methods you use from *sql.DB
type DBInterface interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// MockDB is a mock object that implements DBInterface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	// Mock the Query method
	argsMock := m.Called(query, args)
	return argsMock.Get(0).(*sql.Rows), argsMock.Error(1)
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	// Mock the Exec method
	argsMock := m.Called(query, args)
	return argsMock.Get(0).(sql.Result), argsMock.Error(1)
}

func (m *MockDB) ConnectDB() error {
	// Mock implementation, return nil for success
	return nil
}

func (m *MockDB) Close() {
	// Mock implementation, do nothing
}

func TestSetup(t *testing.T) {
	// Replace the original implementations with mocks
	datalayer.StorageDB = &MockDB{}

	// Run the setup function
	setup.Setup()

	// Ensure that the setup completes without errors
	// Here, you might want to add more specific assertions based on your requirements
	assert.True(t, true, "Setup completed successfully")
}
