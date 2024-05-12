package epoch

import (
	"net/http"

	li "github.com/predixus/pdb_framework/internal/logger"
)

func EpochHandler(w http.ResponseWriter, r *http.Request) {
	li.Logger.Println("Registering Epoch")
	w.WriteHeader(http.StatusOK)
}
