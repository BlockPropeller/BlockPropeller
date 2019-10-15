package account

import (
	"blockpropeller.dev/blockpropeller"
	"github.com/urfave/cli"
)

// Cmd is an umbrella command for operations targeted at a specific account.
func Cmd(app *blockpropeller.App) cli.Command {
	return cli.Command{
		Name:  "account",
		Usage: "Account related commands",
		Subcommands: []cli.Command{
			listCmd(app),
			createCmd(app),
		},
	}
}
