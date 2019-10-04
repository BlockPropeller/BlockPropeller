package admin

import (
	"chainup.dev/chainup"
	"chainup.dev/chainup/cmd/chainctl/admin/account"
	"chainup.dev/chainup/cmd/chainctl/admin/job"
	"chainup.dev/chainup/cmd/chainctl/admin/server"
	"github.com/urfave/cli"
)

// Cmd is an umbrella command for admin operations.
func Cmd(app *chainup.App) cli.Command {
	return cli.Command{
		Name:  "admin",
		Usage: "ChainUP administration commands",
		Subcommands: []cli.Command{
			account.Cmd(app),
			server.Cmd(app),
			job.Cmd(app),
		},
	}
}
