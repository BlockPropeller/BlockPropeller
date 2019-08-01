package main

import (
	"log"
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
	app := chainup.SetupInMemoryApp()

	cmd := AppCmd(app)

	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
