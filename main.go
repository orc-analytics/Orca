package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	api "github.com/predixus/analytics_framework/src/api"
	grpc "github.com/predixus/analytics_framework/src/grpc"
)

func main() {
	grpc.StartGRPCServer()

	var wait time.Duration
	flag.DurationVar(
		&wait,
		"graceful-timeout",
		time.Second*15,
		"The duration for which the server gracefull waits for existing connections to finish.",
	)
	flag.Parse()

	// Route definitions
	r := api.GenerateRouter()
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
	os.Exit(0)
}
