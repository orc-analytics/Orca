package datalayer_provision

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "github.com/predixus/analytics_framework/protobufs/go"
)

func TestTableSchemaCorrect(t *testing.T) {
	assert := assert.New(t)
	var epoch protoreflect.ProtoMessage = protoreflect.ProtoMessage(&pb.Epoch{
		EpochStart: "test_start",
		EpochEnd:   "test_end",
		Origin: &pb.Origin{
			Name: "test_origin",
		},
	})
	schema := generateTableSchema(&epoch)

	assert.Equal(t, schema["origin"], "name")
}
