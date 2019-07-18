package main

import (
	"chainup.dev/chainup"
	"chainup.dev/chainup/binance"
	"chainup.dev/chainup/provision"
	"chainup.dev/lib/log"
	"github.com/urfave/cli"
)

// DeployCmd is a command for creating new infrastructure and deploying a Binance Chain node on top of it.
//
// This command serves as an MVP for the infrastructure and provisioning of ChainUP.
func DeployCmd() cli.Command {
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
				Value: chainup.ProviderDigitalOcean.String(),
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

			providerType := chainup.NewProviderType(c.String("provider"))
			if !providerType.IsValid() {
				log.Error("Invalid provider type flag.", log.Fields{
					"valid_types": chainup.ValidProviders,
				})
				return
			}

			log.Info("Starting provisioning process...", log.Fields{
				"network":       network.String(),
				"node_type":     nodeType.String(),
				"provider_type": providerType.String(),
			})

			job := provision.NewJobBuilder().Build()

			provisioner := chainup.SetupInMemoryProvisioner()

			//@TODO: Create resource creation request for machines that need to be created and services that need to be running on top.
			//@TODO: Kick-off the provisioning process.
			//@TODO: Wait for the process to complete and return the results to the user.
			err := provisioner.WaitFor(job)
			if err != nil {
				log.ErrorErr(err, "Failed running provisioner job")
				return
			}

			log.Info("Finished provisioning job", log.Fields{
				"id":    job.ID,
				"state": job.State,
			})
		},
	}
}
