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
type DB interface {
	Connect() error
	Close() error
}

// DBConnection is a struct for connecting to the database
type Db struct {
	DB *sql.DB
}

// ConnectDB connects to the database and returns a db instance
func (c *Db) Connect() error {
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

	c.DB = db
	return nil
}

func (c *Db) Close() error {
	if c.DB == nil {
		return nil
	}

	err := c.DB.Close()
	c.DB = nil
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
