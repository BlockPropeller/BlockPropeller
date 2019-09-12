package integration_test

import (
	"context"
	"testing"
	"time"

	"chainup.dev/chainup"
	"chainup.dev/chainup/binance"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"chainup.dev/lib/test"
	"github.com/blang/semver"
)

func TestProvisioningJob(t *testing.T) {
	test.Integration(t)

	app := chainup.SetupTestApp(t)

	provider := infrastructure.NewProviderSettings(
		infrastructure.ProviderDigitalOcean, app.Config.DigitalOcean.AccessToken)

	server, err := infrastructure.NewServerBuilder().
		Provider(provider.Type).
		Build()
	test.CheckErr(t, "build server spec", err)

	job, err := provision.NewJobBuilder().
		Provider(provider).
		Server(server).
		Deployment(binance.NewNodeDeployment(
			binance.NetworkTest,
			binance.TypeLightNode,
			semver.MustParse("0.6.1"),
		)).
		Build()
	test.CheckErr(t, "build job spec", err)

	err = app.Provisioner.Provision(context.Background(), job)
	defer func() {
		if job.WorkspaceSnapshot == nil {
			return
		}

		// Destroy infrastructure created for the test.
		err = app.Provisioner.Undo(context.Background(), job)
		test.CheckErr(t, "undo infrastructure", err)
	}()
	test.CheckErr(t, "run deploy command", err)

	//@TODO: Check for persistence later.
	//srv, err := app.ServerRepository.Find(job.Server.ID)
	//test.CheckErr(t, "find requested server", err)
	srv := job.Server

	test.AssertBoolEqual(t, "sever provisioning state",
		srv.State.IsSuccessful, true)
	test.AssertStringsEqual(t, "server provider",
		srv.Provider.String(), infrastructure.ProviderDigitalOcean.String())

	test.AssertIntsEqual(t, "server has deployment", len(srv.Deployments), 1)
	test.AssertStringsEqual(t, "deployment ready",
		srv.Deployments[0].State.String(), infrastructure.DeploymentStateRunning.String())

	// Test binance chain node is reachable.
	time.Sleep(5 * time.Second)

	err = infrastructure.CheckHealth(srv, srv.Deployments[0])
	test.CheckErr(t, "check node health", err)
}
