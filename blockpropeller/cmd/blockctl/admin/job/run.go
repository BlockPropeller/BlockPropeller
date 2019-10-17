package job

import (
	"context"

	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/blockpropeller/binance"
	"blockpropeller.dev/blockpropeller/cmd/blockctl/util/localauth"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/provision"
	"blockpropeller.dev/lib/log"
	"github.com/blang/semver"
	"github.com/urfave/cli"
)

// runCmd is a command for creating new infrastructure and deploying a Binance Chain node on top of it.
//
// This command serves as an MVP for the infrastructure and provisioning of BlockPropeller.
func runCmd(app *blockpropeller.App) cli.Command {
	return cli.Command{
		Name:  "run",
		Usage: "Run a new Binance Chain node job.",
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
			cli.StringFlag{
				Name:  "key",
				Usage: "Cloud provider access key to use for provisioning infrastructure.",
			},
		},
		Action: func(c *cli.Context) {
			acc := localauth.Account

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

			providerKey := app.Config.DigitalOcean.AccessToken
			if c.String("key") != "" {
				providerKey = c.String("key")
			}
			if providerKey == "" {
				log.Error("Invalid provider key flag. The key must not be empty.")
				return
			}

			log.Info("Starting provisioning process...", log.Fields{
				"network":       network.String(),
				"node_type":     nodeType.String(),
				"provider_type": providerType.String(),
			})

			provider := infrastructure.NewProviderSettings(acc.ID, "Provider Settings", providerType, providerKey)
			err := app.ProviderSettingsRepository.Create(context.TODO(), provider)
			if err != nil {
				log.ErrorErr(err, "Failed saving provider settings.")
				return
			}

			srv, err := infrastructure.NewServerBuilder(acc.ID).
				Provider(provider.Type).
				Build()

			job, err := provision.NewJobBuilder(acc.ID).
				Provider(provider).
				Server(srv).
				Deployment(binance.NewNodeDeployment(
					network,
					nodeType,
					semver.MustParse("0.6.1"),
				)).
				Build()

			err = app.JobScheduler.Schedule(context.TODO(), job)
			if err != nil {
				log.ErrorErr(err, "Failed scheduling job")
				return
			}

			err = app.Provisioner.Provision(context.Background(), job)
			if err != nil {
				log.ErrorErr(err, "Failed running server state machine")
				return
			}

			log.Info("Finished provisioning job", log.Fields{
				"id":    job.ID,
				"state": job.State,
			})
		},
	}
}
