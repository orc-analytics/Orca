package datalayer

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	pb "github.com/predixus/analytics_framework/protobufs/go"
)

var StorageDB *sql.DB

type Epoch pb.Epoch

func ConnectDB() {
	// load in variables
	godotenv.Load()
	var (
		host     = os.Getenv("DB_IP")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)

	// start a connection to the storage db
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	StorageDB = db
}

// func (epoch Epoch) WriteEpoch() {
//   epoch.Type.Version.if
// 	insertSmt := fmt.Sprintf("INSERT INTO epoch (epoch_start, epoch_end, origin, type, ) values ();")
//    e
//
// 	StorageDB.Query(insertSmt)
// }

//
// func ReadEpoch(epoch *pb.Epoch) {
//   ...
// }

// func (epoch *pb.Epoch) UpdateEpoch(ctx context.Context) {
//   ...
// }

// func (epoch *pb.Epoch) UpdateEpoch(ctx context.Context){
//   ...
// }
