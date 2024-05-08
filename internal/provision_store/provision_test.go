package datalayer_provision_test

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"

	prov "github.com/predixus/analytics_framework/internal/provision_store"
)

func TestGenerateCreateTableStatement(t *testing.T) {
	msg := mockProtoMessage{ID: 1234, Name: "Mock"}

	tableMap := map[string]string{
		"ID":   "INT",
		"Name": "VARCHAR(255)",
		"Age":  "INT",
	}

	expectedTableName := "mockmessagename"
	expectedQueryColumns := "id INT, name VARCHAR(255), age INT"
	expectedQuery := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (%s);",
		expectedTableName,
		expectedQueryColumns,
	)

	result := prov.GenerateCreateTableStatement(&msg, tableMap)

	if result != expectedQuery {
		t.Errorf("Unexpected query generated: got %s, want %s", result, expectedQuery)
	}
}

// Mock ProtoMessage implementation for testing
type mockProtoMessage struct {
	ID   int32
	Name string
}

func (m *mockProtoMessage) ProtoReflect() protoreflect.Message {
	return nil // Return a mock protoreflect.Message for testing
}
