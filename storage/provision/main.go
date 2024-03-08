package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "github.com/predixus/analytics_framework/protobufs/go"
)

func protoToPostgresType(field *protoreflect.FieldDescriptor) string {
	kind := (*field).Kind()

	switch kind {
	case protoreflect.BoolKind:
		return "BOOLEAN"
	case protoreflect.EnumKind,
		protoreflect.Int32Kind,
		protoreflect.Sint32Kind,
		protoreflect.Uint32Kind,
		protoreflect.Sfixed32Kind,
		protoreflect.Fixed32Kind:
		return "INT"
	case protoreflect.Int64Kind,
		protoreflect.Sint64Kind,
		protoreflect.Uint64Kind,
		protoreflect.Sfixed64Kind,
		protoreflect.Fixed64Kind:
		return "BIGINT"
	case protoreflect.FloatKind:
		return "REAL"
	case protoreflect.DoubleKind:
		return "DOUBLE PRECISION"
	case protoreflect.StringKind:
		return "TEXT"
	case protoreflect.BytesKind:
		return "BYTEA"
	case protoreflect.MessageKind:
		return "MESSAGE"
	default:
		return "UNKNOWN"
	}
}

func generateTableSchema(
	msg *protoreflect.MessageDescriptor,
	prePendParent bool,
) map[string]string {
	pgFieldMap := make(map[string]string)

	fields := (*msg).Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		fieldName := field.Name()
		fieldType := field.Kind()
		pgType := protoToPostgresType(&field)
		if pgType == "MESSAGE" {
			nestedMessage := field.Message()
			subMap := generateTableSchema(&nestedMessage, true)
			for k, v := range subMap {
				pgFieldMap[k] = v
			}
		} else {
			if prePendParent {
				pgFieldMap[fmt.Sprintf("%s_%s", field.ContainingMessage().Name(), fieldName)] = pgType
			} else {
				pgFieldMap[string(fieldName)] = pgType
			}
		}

		fmt.Println(
			fmt.Sprintf(
				"Field Name: %s - Type: %s - PostgresType: %s",
				fieldName,
				fieldType,
				pgType,
			),
		)
	}
	return pgFieldMap
}

func generateCreateTableStatement(tableName string, tableMap map[string]string) string {
	var columns []string
	for columnName, columnType := range tableMap {
		columns = append(columns, fmt.Sprintf("%s %s", columnName, columnType))
	}
	return fmt.Sprintf("CREATE TABLE %s (%s);", tableName, strings.Join(columns, ", "))
}

func main() {
	godotenv.Load("../../.env")
	var (
		host     = os.Getenv("DB_IP")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)
	fmt.Println(dbname)
	fmt.Println(password)
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	msg := &pb.Epoch{}
	desc := msg.ProtoReflect().Descriptor()
	pgFieldMap := generateTableSchema(&desc, false)
	createStatement := generateCreateTableStatement(string(desc.Name()), pgFieldMap)
	fmt.Println(createStatement)

	_, err = db.Query(createStatement)
	if err != nil {
		log.Fatal(err)
	}
}
