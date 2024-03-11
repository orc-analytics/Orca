package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	api "github.com/predixus/analytics_framework/internal/api"
	datalayer "github.com/predixus/analytics_framework/internal/datalayer"
	grpc "github.com/predixus/analytics_framework/internal/grpc"
)

func main() {
	// load in variables
	godotenv.Load()
	var (
		host     = os.Getenv("DB_IP")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)

	// start a connection to the MetaDB
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s ",
		host, port, user, password, dbname)
	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr)))
	datalayer.MetaDB = bun.NewDB(db, pgdialect.New())
	defer datalayer.MetaDB.Close()

	// start a waitgroup
	var wg sync.WaitGroup

	// Start the gRPC framework server
	wg.Add(1)
	go func() {
		grpc.StartGRPCServer(&wg)
	}()

	// Start the API - TODO: kick out into seperate module
	wg.Add(1)
	go func() {
		api.StartHTTPServer(&wg)
	}()

	wg.Wait()
}
