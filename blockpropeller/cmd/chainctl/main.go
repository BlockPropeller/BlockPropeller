package main

import (
	"fmt"
	"os"

	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/blockpropeller/cmd/chainctl/admin"
	"blockpropeller.dev/blockpropeller/cmd/chainctl/auth"
	"blockpropeller.dev/blockpropeller/cmd/chainctl/util/localauth"
	"blockpropeller.dev/lib/log"
	"github.com/urfave/cli"
)

// AppCmd is the top level command wrapping all CLI capabilities of BlockPropeller.
func AppCmd(app *blockpropeller.App) *cli.App {
	cmd := cli.NewApp()
	cmd.Name = "chainctl"
	cmd.Usage = "ChainCTL is a command line utility created to ease the development of BlockPropeller."
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

	cmd := AppCmd(app)

	err = cmd.Run(os.Args)
	if err != nil {
		log.ErrorErr(err, "Failed running command")
		os.Exit(1)
	}
}
