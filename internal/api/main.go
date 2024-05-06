package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type HTTPServer interface {
	Start(wg *sync.WaitGroup)
}

type HttpServer struct {
	*http.Server
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Request:")
	w.WriteHeader(http.StatusOK)
}

func generateRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	return r
}

func (s *HttpServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	// Route definitions
	r := generateRouter()
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "4040"
	}
	s.Server = &http.Server{
		Addr:         fmt.Sprintf("localhost:%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// setup logging
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	// create a channel to listen for interrupts
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// block until the signal is received
	<-c

	// create a deadline
	// this will not block if there are no connections
	// if there are, it will wait until the timeout completes
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server Shutdown: %v", err)
	}

	log.Println("HTTP: Shutting Down")
}
