package main

import (
	"fmt"
	"os"

	"chainup.dev/chainup"
	"chainup.dev/chainup/cmd/chainctl/admin"
	"chainup.dev/chainup/cmd/chainctl/auth"
	"chainup.dev/chainup/cmd/chainctl/util/localauth"
	"chainup.dev/lib/log"
	"github.com/urfave/cli"
)

// AppCmd is the top level command wrapping all CLI capabilities of ChainUP.
func AppCmd(app *chainup.App) *cli.App {
	cmd := cli.NewApp()
	cmd.Name = "chainctl"
	cmd.Usage = "ChainCTL is a command line utility created to ease the development of ChainUP."
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
	app, closeFn, err := chainup.SetupDatabaseApp()
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
