package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/predixus/analytics_framework/src/routeHandler"
)

func main() {
	var wait time.Duration
	flag.DurationVar(
		&wait,
		"graceful-timeout",
		time.Second*15,
		"The duration for which the server gracefull waits for existing connections to finish.",
	)
	flag.Parse()

	// Route definitions
	r := routeHandler.GenerateRouter()
	srv := &http.Server{
		Addr:         "localhost:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

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
