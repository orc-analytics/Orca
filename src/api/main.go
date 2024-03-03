package api

import (
	"net/http"

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
