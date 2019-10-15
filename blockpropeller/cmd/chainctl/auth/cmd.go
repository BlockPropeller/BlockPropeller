package auth

import (
	"blockpropeller.dev/blockpropeller"
	"github.com/urfave/cli"
)

// Cmd is an umbrella command for admin operations.
func Cmd(app *blockpropeller.App) cli.Command {
	return cli.Command{
		Name:  "auth",
		Usage: "Authentication commands",
		Subcommands: []cli.Command{
			registerCmd(app),
			loginCmd(app),
			whoamiCmd(app),
			logoutCmd(),
		},
	}
}
