package server

import (
	"blockpropeller.dev/blockpropeller"
	"github.com/urfave/cli"
)

// Cmd is an umbrella command for operations targeted at a specific server.
func Cmd(app *blockpropeller.App) cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "Server related commands",
		Subcommands: []cli.Command{
			listCmd(app),
			deleteCmd(app),
			dumpKeyCmd(app),
		},
	}
}
