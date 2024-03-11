package main

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "github.com/predixus/analytics_framework/protobufs/go"
)

func TestProtoToPostgresType(t *testing.T) {
	// Test cases for protoToPostgresType function
	tests := []struct {
		name     string
		field    *protoreflect.FieldDescriptor
		expected string
	}{
		{
			name: "Bool",
			field: &protoreflect.FieldDescriptor{
				FieldDescriptorProto: &protoreflect.FieldDescriptorProto{
					Type: &protoreflect.FieldDescriptorProto_Type{
						Kind: &protoreflect.FieldDescriptorProto_Type_Bool{Bool: true},
					},
				},
			},
			expected: "BOOLEAN",
		},
		{
			name: "Enum",
			field: &protoreflect.FieldDescriptor{
				FieldDescriptorProto: &protoreflect.FieldDescriptorProto{
					Type: &protoreflect.FieldDescriptorProto_Type{
						Kind: &protoreflect.FieldDescriptorProto_Type_Enum{
							Enum: &protoreflect.EnumDescriptorProto{},
						},
					},
				},
			},
			expected: "INT",
		},
		// Add more test cases for other types
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := protoToPostgresType(tc.field)
			if got != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestGenerateTableSchema(t *testing.T) {
	// Test cases for generateTableSchema function
	tests := []struct {
		name     string
		msg      proto.Message
		expected map[string]string
	}{
		{
			name: "Epoch",
			msg:  &pb.Epoch{},
			expected: map[string]string{
				"id":   "INT",
				"name": "TEXT",
				// Add more fields as needed
			},
		},
		// Add more test cases for other messages
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			desc := tc.msg.ProtoReflect().Descriptor()
			got := generateTableSchema(&desc, false)
			if !equalMaps(got, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestGenerateAlterTableStatement(t *testing.T) {
	// Test cases for generateAlterTableStatement function
	tests := []struct {
		name     string
		table    string
		tableMap map[string]string
		expected string
	}{
		{
			name:  "Epoch",
			table: "epoch",
			tableMap: map[string]string{
				"id":   "INT",
				"name": "TEXT",
			},
			expected: `ALTER TABLE "epoch" ADD COLUMN IF NOT EXISTS "id" INT, ADD COLUMN IF NOT EXISTS "name" TEXT;`,
		},
		// Add more test cases as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := generateAlterTableStatement(tc.table, tc.tableMap)
			if got != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestGenerateCreateTableStatement(t *testing.T) {
	// Test cases for generateCreateTableStatement function
	tests := []struct {
		name     string
		table    string
		tableMap map[string]string
		expected string
	}{
		{
			name:  "Epoch",
			table: "epoch",
			tableMap: map[string]string{
				"id":   "INT",
				"name": "TEXT",
			},
			expected: `CREATE TABLE IF NOT EXISTS "epoch" ("id" INT, "name" TEXT);`,
		},
		// Add more test cases as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := generateCreateTableStatement(tc.table, tc.tableMap)
			if got != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func equalMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
