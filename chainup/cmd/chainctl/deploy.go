package main

import (
	"context"

	"chainup.dev/chainup"
	"chainup.dev/chainup/binance"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/lib/log"
	"github.com/urfave/cli"
)

// DeployCmd is a command for creating new infrastructure and deploying a Binance Chain node on top of it.
//
// This command serves as an MVP for the infrastructure and provisioning of ChainUP.
func DeployCmd(app *chainup.App) cli.Command {
	return cli.Command{
		Name:  "deploy",
		Usage: "Deploy a new Binance Chain node.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "network",
				Usage: "Network you wish to join with your node.",
				Value: binance.NetworkTest.String(),
			},
			cli.StringFlag{
				Name:  "mode",
				Usage: "Mode in which to deploy the Binance Chain node.",
				Value: binance.TypeLightNode.String(),
			},
			cli.StringFlag{
				Name:  "provider",
				Usage: "Cloud provider to use for provisioning infrastructure.",
				Value: infrastructure.ProviderDigitalOcean.String(),
			},
		},
		Action: func(c *cli.Context) {
			network := binance.NewNetwork(c.String("network"))
			if !network.IsValid() {
				log.Error("Invalid network flag.", log.Fields{
					"valid_networks": binance.ValidNetworks,
				})
				return
			}

			nodeType := binance.NewNodeType(c.String("mode"))
			if !nodeType.IsValid() {
				log.Error("Invalid node type flag.", log.Fields{
					"valid_types": binance.ValidNodeTypes,
				})
				return
			}

			providerType := infrastructure.NewProviderType(c.String("provider"))
			if !providerType.IsValid() {
				log.Error("Invalid provider type flag.", log.Fields{
					"valid_types": infrastructure.ValidProviders,
				})
				return
			}

			log.Info("Starting provisioning process...", log.Fields{
				"network":       network.String(),
				"node_type":     nodeType.String(),
				"provider_type": providerType.String(),
			})

			server := infrastructure.NewServerBuilder().Build()

			err := app.Provisioner.Provision(context.Background(), server)
			if err != nil {
				log.ErrorErr(err, "Failed running server state machine")
				return
			}

			log.Info("Finished provisioning server", log.Fields{
				"id":    server.ID,
				"state": server.State,
			})
		},
	}
}
