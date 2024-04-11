package datalayer_provision

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
	msg *protoreflect.ProtoMessage,
) map[string]string {
	desc := (*msg).ProtoReflect().Descriptor()
	pgFieldMap := make(map[string]string)

	stack := []*protoreflect.MessageDescriptor{&desc}
	prependParent := []bool{false}

	for len(stack) > 0 {
		// get the current message
		currMsg := stack[len(stack)-1]
		prependCurrentMessage := prependParent[len(stack)-1]

		// move the stack forward
		stack = stack[:len(stack)-1]
		prependParent = prependParent[:len(stack)-1]

		fields := (*currMsg).Fields()
		for i := 0; i < fields.Len(); i++ {
			field := fields.Get(i)
			fieldName := field.Name()
			pgType := protoToPostgresType(&field)

			if pgType == "UNKNOWN" {
				log.Fatal(
					fmt.Sprintf(
						"Unsure how to translate pbuf type to PG type: %s",
						field.Kind().String(),
					),
				)
			}

			if pgType == "MESSAGE" {
				// we've encoutered a message. Add it to the stack to handle later
				nestedMessage := field.Message()
				stack = append(stack, &nestedMessage)
				prependParent = append(prependParent, true)
			} else {
				if prependCurrentMessage {
					pgFieldMap[fmt.Sprintf("%s_%s", field.ContainingMessage().Name(), fieldName)] = pgType
				} else {
					pgFieldMap[string(fieldName)] = pgType
				}
			}
		}
	}

	return pgFieldMap
}

func generateAlterTableStatement(
	msg *protoreflect.ProtoMessage,
	tableMap map[string]string,
) string {
	var queryRows []string
	for columnName, columnType := range tableMap {
		queryRows = append(
			queryRows,
			fmt.Sprintf(
				"ADD COLUMN IF NOT EXISTS %s %s",
				strings.ToLower(columnName),
				columnType,
			),
		)
	}
	return fmt.Sprintf(
		"ALTER TABLE %s %s;",
		strings.ToLower(string((*msg).ProtoReflect().Descriptor().Name())),
		strings.Join(queryRows, ", "),
	)
}

func generateCreateTableStatement(
	msg *protoreflect.ProtoMessage,
	tableMap map[string]string,
) string {
	var queryColumns []string
	for columnName, columnType := range tableMap {
		queryColumns = append(
			queryColumns,
			fmt.Sprintf("%s %s", strings.ToLower(columnName), columnType),
		)
	}
	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (%s);",
		strings.ToLower(string((*msg).ProtoReflect().Descriptor().Name())),
		strings.Join(queryColumns, ", "),
	)
}

func main() {
	godotenv.Load()
	var (
		host     = os.Getenv("DB_IP")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)

	messagesToTabularise := make([]protoreflect.ProtoMessage, 3)
	messagesToTabularise[0] = &pb.Epoch{}
	messagesToTabularise[1] = &pb.Algorithm{}
	messagesToTabularise[2] = &pb.Pipeline{}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, msg := range messagesToTabularise {

		pgFieldMap := generateTableSchema(&msg)

		createStatement := generateCreateTableStatement(&msg, pgFieldMap)
		alterStatement := generateAlterTableStatement(&msg, pgFieldMap)
		_, err = db.Query(createStatement)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Query(alterStatement)
		if err != nil {
			log.Fatal(err)
		}
	}
}
