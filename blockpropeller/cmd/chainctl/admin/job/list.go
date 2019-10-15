package job

import (
	"context"
	"os"
	"time"

	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/blockpropeller/cmd/chainctl/util/localauth"
	"blockpropeller.dev/lib/log"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func listCmd(app *blockpropeller.App) cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "List all jobs",
		Action: func(c *cli.Context) {
			acc := localauth.Account

			jobs, err := app.JobRepository.List(context.Background(), acc.ID)
			if err != nil {
				log.ErrorErr(err, "failed listing jobs")
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetAutoMergeCells(true)
			table.SetRowLine(true)
			table.SetHeader([]string{
				"ID",
				"State",
				"",
				"",
			})

			for _, job := range jobs {
				table.Append([]string{
					job.ID.String(),
					job.State.String(),
					"Provider: " + job.ProviderSettingsID.String(),
					"Created: " + job.CreatedAt.Format(time.Stamp),
				})
				table.Append([]string{
					job.ID.String(),
					job.State.String(),
					"Server: " + job.ServerID.String(),
					"",
				})
				table.Append([]string{
					job.ID.String(),
					job.State.String(),
					"Deployment: " + job.DeploymentID.String(),
					"Finished: " + job.FinishedAt.Format(time.Stamp),
				})
			}

			table.Render()
		},
	}
}
