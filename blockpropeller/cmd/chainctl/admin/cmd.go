package admin

import (
	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/blockpropeller/cmd/chainctl/admin/account"
	"blockpropeller.dev/blockpropeller/cmd/chainctl/admin/job"
	"blockpropeller.dev/blockpropeller/cmd/chainctl/admin/server"
	"blockpropeller.dev/blockpropeller/cmd/chainctl/util/localauth"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// Cmd is an umbrella command for admin operations.
func Cmd(app *blockpropeller.App) cli.Command {
	return cli.Command{
		Name:  "admin",
		Usage: "BlockPropeller administration commands",
		Before: func(ctx *cli.Context) error {
			acc := localauth.Account
			if acc == nil {
				return errors.New("access denied: missing authenticated account")
			}

			return nil
		},
		Subcommands: []cli.Command{
			account.Cmd(app),
			server.Cmd(app),
			job.Cmd(app),
		},
	}
}
