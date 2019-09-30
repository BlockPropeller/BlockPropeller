package server

import (
	"context"
	"os"

	"chainup.dev/chainup"
	"chainup.dev/lib/log"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func listCmd(app *chainup.App) cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "List all servers",
		Action: func(c *cli.Context) {
			servers, err := app.ServerRepository.List(context.Background())
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
					srv.IPAddress.String(),
				})
			}

			table.Render()
		},
	}
}
