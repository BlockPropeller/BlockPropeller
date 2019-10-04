package account

import (
	"context"
	"os"
	"time"

	"chainup.dev/chainup"
	"chainup.dev/lib/log"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func listCmd(app *chainup.App) cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "List all accounts",
		Action: func(c *cli.Context) {
			accounts, err := app.AccountRepository.List(context.Background())
			if err != nil {
				log.ErrorErr(err, "failed listing accounts")
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Email", "CreatedAt"})

			for _, acc := range accounts {
				table.Append([]string{
					acc.ID.String(),
					acc.Email.String(),
					acc.CreatedAt.Format(time.Stamp),
				})
			}

			table.Render()
		},
	}
}
