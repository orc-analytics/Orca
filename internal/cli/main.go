package cli

import (
	"flag"

	li "github.com/predixus/analytics_framework/internal/logger"
	provision "github.com/predixus/analytics_framework/internal/provision_store"
)

func ParseInputs() {
	// first check whether there are any command line arguments
	initPtr := flag.Bool("init-db", false, "Provision the local postgres db")
	flag.Parsed()

	if *initPtr {
		println("Initialising postgres DB")
		err := provision.Provision()
		if err != nil {
			li.Logger.Fatal(err)
		}
		return
	}
}
