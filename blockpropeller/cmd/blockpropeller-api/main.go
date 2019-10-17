package main

import (
	"context"
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

	err = appSrv.Start(context.Background())
	if err != nil {
		log.ErrorErr(err, "Failed running app server")
		os.Exit(1)
	}
}
