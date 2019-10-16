package main

import (
	"os"

	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/lib/log"
)

func main() {
	appSrv, closeFn, err := blockpropeller.SetupDatabaseServer()
	if err != nil {
		log.ErrorErr(err, "Failed setting up application server")
		os.Exit(1)
	}
	defer closeFn()

	err = appSrv.Start()
	if err != nil {
		log.ErrorErr(err, "Failed running HTTP server")
		os.Exit(1)
	}
}
