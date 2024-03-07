package api

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/predixus/analytics_framework/src/api/epoch"
	li "github.com/predixus/analytics_framework/src/logger"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	li.Logger.Printf("Recieved Request:")
	w.WriteHeader(http.StatusOK)
}

func GenerateRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/register-epoch", epoch.EpochHandler)

	return r
}

func StartHTTPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	// Start the HTTP API
	var wait time.Duration
	flag.DurationVar(
		&wait,
		"graceful-timeout",
		time.Second*15,
		"The duration for which the server gracefull waits for existing connections to finish.",
	)
	flag.Parse()

	// Route definitions
	r := GenerateRouter()
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "4040"
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// setup logging
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// create a channel to listen for interrupts
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// block until the signal is recieved
	<-c

	// create a deadline
	// this will not block if there are no connections
	// if there are, it will wait untill the timeout completes
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("Shutting Down")
}
