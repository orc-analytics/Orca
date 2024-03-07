package provision

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/proto"

	pb "github.com/predixus/analytics_framework/protobufs/go"
)

func main() {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database_ip := os.Getenv("DB_IP")
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/public", username, password, database_ip)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

// Example protobuf message
type MyMessage struct {
	Id   int32
	Name string
}

func main_v2() {
	// Create an instance of your protobuf message
	msg := &pb.Epoch{}

	// Get the descriptor for the protobuf message
	desc := proto.MessageReflect(msg).Descriptor()

	// Generate PostgreSQL table schema based on the message descriptor
	schema := generateTableSchema(desc)

	// Print the PostgreSQL table schema
	fmt.Println(schema)
}

// Function to generate PostgreSQL table schema
func generateTableSchema(desc *descriptor.MessageDescriptorProto) string {
	var fields []string

	// Iterate through the fields of the protobuf message
	for _, field := range desc.Field {
		fieldName := field.GetName()
		fieldType := fieldTypeMapping(field.GetType())

		// Add field definition to the schema
		fields = append(fields, fmt.Sprintf("%s %s", fieldName, fieldType))
	}

	// Join all field definitions to form the table schema
	return fmt.Sprintf("CREATE TABLE %s (\n\t%s\n);", desc.GetName(), strings.Join(fields, ",\n\t"))
}

// Function to map protobuf field types to PostgreSQL data types
func fieldTypeMapping(fieldType descriptor.FieldDescriptorProto_Type) string {
	switch fieldType {
	case descriptor.FieldDescriptorProto_TYPE_INT32, descriptor.FieldDescriptorProto_TYPE_INT64:
		return "INTEGER"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "TEXT"
	// Add more mappings for other types as needed
	default:
		return "UNKNOWN"
	}
}
