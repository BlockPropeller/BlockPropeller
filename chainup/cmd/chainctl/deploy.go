package main

import (
	"chainup.dev/chainup/binance"
	"chainup.dev/lib/log"
	"github.com/urfave/cli"
)

// DeployCmd is a command for creating new infrastructure and deploying a Binance Chain node on top of it.
//
// This command serves as an MVP for the infrastructure and provisioning of ChainUP.
func DeployCmd() cli.Command {
	return cli.Command{
		Name:        "deploy",
		Description: "Deploy a new Binance Chain node.",
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
		},
		Action: func(c *cli.Context) {
			network := binance.NewNetwork(c.String("network"))
			if !network.IsValid() {
				log.Error("Invalid network flag.", log.Fields{
					"valid_networks": binance.ValidNetworks,
				})
			}

			nodeType := binance.NewNodeType(c.String("mode"))
			if !nodeType.IsValid() {
				log.Error("Invalid node type flag.", log.Fields{
					"valid_types": binance.ValidNodeTypes,
				})
			}

			log.Info("Starting provisioning process...")
		},
	}
}
