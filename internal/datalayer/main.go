package datalayer

import (
	"log"

	"github.com/uptrace/bun"

	pb "github.com/predixus/analytics_framework/protobufs/go"
)

var MetaDB *bun.DB

func WriteEpoch(epoch *pb.Epoch) {
	var query_string string

	_, err := MetaDB.Query(query_string)
	if err != nil {
		log.Fatal(err)
	}
}

//
// func ReadEpoch(epoch *pb.Epoch) {
//   ...
// }
