package main

import (
	"fmt"
	"os"

	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/blockpropeller/cmd/blockctl/admin"
	"blockpropeller.dev/blockpropeller/cmd/blockctl/auth"
	"blockpropeller.dev/blockpropeller/cmd/blockctl/util/localauth"
	"blockpropeller.dev/blockpropeller/encryption"
	"blockpropeller.dev/lib/log"
	"github.com/urfave/cli"

	_ "blockpropeller.dev/blockpropeller/terraform/cloudprovider/digitalocean"
)

// AppCmd is the top level command wrapping all CLI capabilities of BlockPropeller.
func AppCmd(app *blockpropeller.App) *cli.App {
	cmd := cli.NewApp()
	cmd.Name = "blockctl"
	cmd.Usage = "BlockCTL is a command line utility created to ease the development of BlockPropeller."
	cmd.Before = func(*cli.Context) error {
		localauth.Authenticate(app)
		return nil
	}
	cmd.Commands = []cli.Command{
		auth.Cmd(app),
		admin.Cmd(app),
	}

	return cmd
}

func main() {
	app, closeFn, err := blockpropeller.SetupDatabaseApp()
	if err != nil {
		fmt.Printf("failed setting up database app: %s\n", err)
		os.Exit(1)
	}
	defer closeFn()

	log.SetGlobal(app.Logger)
	encryption.Init(app.Config.Encryption.Secret)

	cmd := AppCmd(app)

	err = cmd.Run(os.Args)
	if err != nil {
		log.ErrorErr(err, "Failed running command")
		os.Exit(1)
	}
}
