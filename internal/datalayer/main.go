package datalayer

import (
	"database/sql"
	"log"

	pb "github.com/predixus/analytics_framework/protobufs/go"
)

var MetaDB *sql.DB

func WriteEpoch(epoch *pb.Epoch) {
	err := MetaDB.Query(query_string)
	if err != nil {
		log.Fatal(err)
	}
}

//
// func ReadEpoch(epoch *pb.Epoch) {
//   ...
// }
