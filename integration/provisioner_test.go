package integration_test

import (
	"context"
	"testing"

	"chainup.dev/chainup"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"chainup.dev/lib/test"
)

func TestProvisioningJob(t *testing.T) {
	test.Integration(t)

	app := chainup.SetupTestApp()

	provider := infrastructure.NewProviderSettings(
		infrastructure.ProviderDigitalOcean, app.Config.DigitalOcean.AccessToken)

	job, err := provision.NewJobBuilder().
		Provider(provider).
		Build()
	test.CheckErr(t, "build job spec", err)

	err = app.Provisioner.Provision(context.Background(), job.Server)
	test.CheckErr(t, "run deploy command", err)

	srv, err := app.ServerRepository.Find(job.Server.ID)
	test.CheckErr(t, "find requested server", err)

	test.AssertBoolEqual(t, "sever provisioning state",
		srv.State.IsSuccessful, true)
	test.AssertStringsEqual(t, "server provider",
		srv.Provider.String(), infrastructure.ProviderDigitalOcean.String())
}
