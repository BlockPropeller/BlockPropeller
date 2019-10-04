package auth

import (
	"chainup.dev/chainup"
	"github.com/urfave/cli"
)

// Cmd is an umbrella command for admin operations.
func Cmd(app *chainup.App) cli.Command {
	return cli.Command{
		Name:  "auth",
		Usage: "Authentication commands",
		Subcommands: []cli.Command{
			loginCmd(app),
			whoamiCmd(app),
			logoutCmd(),
		},
	}
}
