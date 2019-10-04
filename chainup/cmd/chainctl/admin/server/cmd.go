package server

import (
	"chainup.dev/chainup"
	"github.com/urfave/cli"
)

// Cmd is an umbrella command for operations targeted at a specific server.
func Cmd(app *chainup.App) cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "Server related commands",
		Subcommands: []cli.Command{
			listCmd(app),
			deleteCmd(app),
		},
	}
}
