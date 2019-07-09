package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "chainctl"
	app.Usage = "ChainCTL is a command line utility created to ease the development of ChainUP."
	app.Commands = []cli.Command{
		DeployCmd(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
