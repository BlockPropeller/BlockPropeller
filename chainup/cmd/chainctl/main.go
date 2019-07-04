package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"chainup.dev/chainup"
)

func main() {
	app := cli.NewApp()
	app.Name = "chainctl"
	app.Usage = "ChainCTL is a command line utility created to ease the development of ChainUP."
	app.Action = func(c *cli.Context) error {
		container := chainup.NewContainer("chainup/binance-fullnode-prod:0.5.8")

		fmt.Println(container.String())

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
