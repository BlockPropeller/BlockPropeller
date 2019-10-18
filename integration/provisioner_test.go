package integration_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"blockpropeller.dev/blockpropeller"
	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/binance"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/provision"
	"blockpropeller.dev/lib/test"
	"github.com/blang/semver"

	_ "blockpropeller.dev/blockpropeller/terraform/cloudprovider/digitalocean"
)

func TestProvisioningJob(t *testing.T) {
	test.Integration(t)

	app := blockpropeller.SetupTestApp(t)

	acc := createTestAccount(t, app)

	provider := infrastructure.NewProviderSettings(
		acc.ID, "Test Provider", infrastructure.ProviderDigitalOcean, app.Config.DigitalOcean.AccessToken)

	server, err := infrastructure.NewServerBuilder(acc.ID).
		Provider(provider.Type).
		Build()
	test.CheckErr(t, "build server spec", err)

	job, err := provision.NewJobBuilder(acc.ID).
		Provider(provider).
		Server(server).
		Deployment(binance.NewNodeDeployment(
			binance.NetworkTest,
			binance.TypeLightNode,
			semver.MustParse("0.6.1"),
		)).
		Build()
	test.CheckErr(t, "build job spec", err)

	err = app.JobScheduler.Schedule(context.TODO(), job)
	test.CheckErr(t, "schedule job", err)

	err = app.Provisioner.Provision(context.TODO(), job)
	defer func() {
		if job.Server == nil || job.Server.WorkspaceSnapshot == nil {
			return
		}

		// Destroy infrastructure created for the test.
		err = app.Provisioner.Undo(context.Background(), job)
		test.CheckErr(t, "undo infrastructure", err)
	}()
	test.CheckErr(t, "run deploy command", err)

	srv, err := app.ServerRepository.Find(context.TODO(), server.ID)
	test.CheckErr(t, "find requested server", err)

	test.AssertStringsEqual(t, "sever provisioning state",
		srv.State.String(), infrastructure.ServerStateOk.String())
	test.AssertStringsEqual(t, "server provider",
		srv.Provider.String(), infrastructure.ProviderDigitalOcean.String())

	test.AssertIntsEqual(t, "server has deployment", len(srv.Deployments), 1)
	test.AssertStringsEqual(t, "deployment ready",
		srv.Deployments[0].State.String(), infrastructure.DeploymentStateOk.String())

	// Test binance chain node is reachable.
	time.Sleep(5 * time.Second)

	err = infrastructure.CheckHealth(srv, srv.Deployments[0])
	test.CheckErr(t, "check node health", err)
}

func createTestAccount(t *testing.T, app *blockpropeller.App) *account.Account {
	acc := account.NewAccount(account.Email(fmt.Sprintf("test-%d@example.com", rand.Int())), "")

	err := app.AccountRepository.Create(context.Background(), acc)
	test.CheckErr(t, "create account", err)

	return acc
}
