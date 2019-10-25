package main

import (
	"context"
	"os"

	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/lib/log"

	_ "blockpropeller.dev/blockpropeller/terraform/cloudprovider/digitalocean"
)

func main() {
	appSrv, closeFn, err := blockpropeller.SetupDatabaseServer()
	if err != nil {
		log.ErrorErr(err, "Failed setting up application server")
		os.Exit(1)
	}
	defer closeFn()

	appSrv.App.InitGlobal()

	err = appSrv.Start(context.Background())
	if err != nil {
		log.ErrorErr(err, "Failed running app server")
		os.Exit(1)
	}
}
