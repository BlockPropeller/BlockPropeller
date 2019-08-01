package integration_test

import (
	"context"
	"testing"

	"chainup.dev/chainup"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/lib/test"
)

func TestProvisioning(t *testing.T) {
	test.Integration(t)

	app := chainup.SetupInMemoryApp()

	srvRequest := infrastructure.NewServerBuilder().Build()

	err := app.Provisioner.Provision(context.Background(), srvRequest)
	test.CheckErr(t, "run deploy command", err)

	srv, err := app.ServerRepository.Find(srvRequest.ID)
	test.CheckErr(t, "find requested server", err)

	test.AssertBoolEqual(t, "sever provisioning state",
		srv.State.IsSuccessful, true)
	test.AssertStringsEqual(t, "server provider",
		srv.Provider.String(), infrastructure.ProviderDigitalOcean.String())
}
