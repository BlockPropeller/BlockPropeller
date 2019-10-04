package main

import (
	"os"

	"chainup.dev/chainup"
	"chainup.dev/lib/log"
)

func main() {
	appSrv, closeFn, err := chainup.SetupDatabaseServer()
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
