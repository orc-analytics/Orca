package epoch

import (
	"log"
	"net/http"
)

func EpochHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Registering Epoch")
	w.WriteHeader(http.StatusOK)
}
