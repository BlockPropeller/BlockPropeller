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

	server := infrastructure.NewServerBuilder().Build()

	err := app.Provisioner.Provision(context.Background(), server)
	test.CheckErr(t, "run deploy command", err)
}
