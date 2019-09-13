package main

import (
	"fmt"
	"os"

	"chainup.dev/chainup"
	"github.com/urfave/cli"
)

// AppCmd is the top level command wrapping all CLI capabilities of ChainUP.
func AppCmd(app *chainup.App) *cli.App {
	cmd := cli.NewApp()
	cmd.Name = "chainctl"
	cmd.Usage = "ChainCTL is a command line utility created to ease the development of ChainUP."
	cmd.Commands = []cli.Command{
		DeployCmd(app),
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
		fmt.Println(err)
		os.Exit(1)
	}
}
