package server

import (
	"context"
	"os"

	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/blockpropeller/cmd/blockctl/util/localauth"
	"blockpropeller.dev/lib/log"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func listCmd(app *blockpropeller.App) cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "List all servers",
		Action: func(c *cli.Context) {
			acc := localauth.Account

			servers, err := app.ServerRepository.List(context.Background(), acc.ID)
			if err != nil {
				log.ErrorErr(err, "failed listing servers")
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "State", "Provider", "IP"})

			for _, srv := range servers {
				table.Append([]string{
					srv.ID.String(),
					srv.Name,
					srv.State.String(),
					srv.Provider.String(),
					srv.IPAddress,
				})
			}

			table.Render()
		},
	}
}
