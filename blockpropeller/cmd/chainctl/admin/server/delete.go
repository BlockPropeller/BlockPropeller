package server

import (
	"context"

	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/lib/log"
	"github.com/urfave/cli"
)

func deleteCmd(app *blockpropeller.App) cli.Command {
	return cli.Command{
		Name:  "delete",
		Usage: "Delete a specified server",
		Action: func(c *cli.Context) {
			if !c.Args().Present() {
				log.Error("please enter a server ID")
				return
			}

			ctx := context.Background()
			srvID := infrastructure.ServerIDFromString(c.Args().First())

			srv, err := app.ServerRepository.Find(ctx, srvID)
			if err != nil {
				log.ErrorErr(err, "could not find server by ID", log.Fields{
					"server_id": srvID,
				})
				return
			}

			err = app.Provisioner.ServerDestroyer.Destroy(ctx, srv)
			if err != nil {
				log.ErrorErr(err, "could not destroy server")
				return
			}

			log.Info("Server successfully deleted.")
		},
	}
}
