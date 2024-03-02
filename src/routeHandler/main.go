package routeHandler

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/predixus/analytics_framework/src/routeHandler/epoch"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved Request:")
	w.WriteHeader(http.StatusOK)
}

func GenerateRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/register-epoch", epoch.EpochHandler)

	return r
}
