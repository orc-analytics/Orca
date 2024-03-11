package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"

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
	var err error
	datalayer.MetaDB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to open connection to the DB: %v", datalayer.MetaDB)
	}
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
