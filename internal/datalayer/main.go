package datalayer

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	li "github.com/predixus/analytics_framework/internal/logger"
	pb "github.com/predixus/analytics_framework/protobufs/go"
)

// Epoch represents an epoch instance
type Epoch pb.Epoch

// DBConnector abstracts the database connection functionality.
type DBConnector interface {
	Connect() *sql.DB
	Close() error
}

// DBConnection is a struct for connecting to the database
type DbConnector struct{}

// ConnectDB connects to the database and returns a db instance
func (c *DbConnector) Connect() *sql.DB {
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
		li.Logger.Fatalf("Could not open DB: %v", err)
	}

	err = db.Ping()
	if err != nil {
		li.Logger.Fatalf("Could not talk to the DB: %v", err)
	}
	return db
}

func (c *DbConnector) Close() error {
	err := c.Close()
	return err
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
