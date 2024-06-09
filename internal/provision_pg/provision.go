package datalayer_provision

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/reflect/protoreflect"

	li "github.com/predixus/pdb_framework/internal/logger"
	pb "github.com/predixus/pdb_framework/protobufs/go"
)

func RemoveIndex[T any](s []*T, index int) ([]*T, error) {
	if index < 0 || index >= len(s) {
		return nil, errors.New("Index is greater than number of elements")
	}

	ret := make([]*T, 0)
	ret = append(ret, s[:index]...)

	return append(ret, s[index+1:]...), nil
}

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

func generateColumnTypes(
	msg *protoreflect.ProtoMessage,
) map[string]string {
	var err error
	desc := (*msg).ProtoReflect().Descriptor()

	// the map from protobuf type to postgres type
	pgFieldMap := make(map[string]string)

	// stack of messages to convert - initially store the high level message
	stack := []*protoreflect.MessageDescriptor{&desc}
	tmpStr := string(desc.Name())
	fieldNames := []*string{&tmpStr}

	for len(stack) > 0 {
		// get the current message
		currMsg := stack[0]
		currFieldName := fieldNames[0]

		fields := (*currMsg).Fields()
		for i := 0; i < fields.Len(); i++ {
			field := fields.Get(i)
			fieldName := field.Name()
			strFieldName := string(fieldName)
			pgType := protoToPostgresType(&field)

			if pgType == "UNKNOWN" {
				log.Fatalf("Unsure how to translate pbuf type to PG type: %s",
					field.Kind().String(),
				)
			}
			index_name := fmt.Sprintf("%s_%s", *currFieldName, strFieldName)

			if pgType == "MESSAGE" {
				// we've encoutered a message. Add it to the stack to handle later
				nestedMessage := field.Message()
				stack = append(stack, &nestedMessage)
				fieldNames = append(fieldNames, &index_name)
			} else {
				pgFieldMap[index_name] = pgType
			}
		}
		// remove the current message from the stack
		stack, err = RemoveIndex(stack, 0)
		if err != nil {
			log.Fatal(
				"Attempted to remove an index in the message stack that does not exist. Aborting.",
			)
		}

		fieldNames, err = RemoveIndex(fieldNames, 0)
		if err != nil {
			log.Fatal(
				"Attempted to remove an index in the field name stack that does not exist. Aborting.",
			)
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

func Provision() error {
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

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		li.Logger.Error(err)
		return err
	}
	defer db.Close()

	for _, msg := range messagesToTabularise {

		pgFieldMap := generateColumnTypes(&msg)

		// for creating a brand new set of tables
		createStatement := generateCreateTableStatement(&msg, pgFieldMap)

		// for altering a table if there are new fields.
		alterStatement := generateAlterTableStatement(&msg, pgFieldMap)
		_, err = db.Query(createStatement)
		if err != nil {
			li.Logger.Error(err)
			return err
		}

		_, err = db.Query(alterStatement)
		if err != nil {
			li.Logger.Error(err)
			return err
		}

	}
	return nil
}
