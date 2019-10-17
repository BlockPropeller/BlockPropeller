//+build wireinject

package blockpropeller

import (
	"testing"

	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/database/transaction"
	"blockpropeller.dev/blockpropeller/httpserver"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/provision"
	"blockpropeller.dev/lib/log"
	"github.com/google/wire"
)

var testAppSet = wire.NewSet(
	ProvideTestConfigProvider,
	log.NewTestingLogger,
	wire.Bind(new(log.Logger), new(*log.TestingLogger)),

	transaction.NewInMemoryTransactionContext,
	wire.Bind(new(transaction.TxContext), new(*transaction.InMemoryTxContext)),

	account.NewInMemoryRepository,
	wire.Bind(new(account.Repository), new(*account.InMemoryRepository)),

	provision.NewInMemoryJobRepository,
	wire.Bind(new(provision.JobRepository), new(*provision.InMemoryJobRepository)),

	infrastructure.NewInMemoryServerRepository,
	wire.Bind(new(infrastructure.ServerRepository), new(*infrastructure.InMemoryServerRepository)),

	infrastructure.NewInMemoryDeploymentRepository,
	wire.Bind(new(infrastructure.DeploymentRepository), new(*infrastructure.InMemoryDeploymentRepository)),

	infrastructure.NewInMemoryProviderSettingsRepository,
	wire.Bind(new(infrastructure.ProviderSettingsRepository), new(*infrastructure.InMemoryProviderSettingsRepository)),

	AppSet,
)

// SetupTestApp constructs an in-memory variant of the StateMachine handling Server state transitions.
func SetupTestApp(t *testing.T) *App {
	panic(wire.Build(
		testAppSet,
	))
}

// SetupTestServer constructs a testing variant of the BlockPropeller Server.
func SetupTestServer(t *testing.T) (*AppServer, func(), error) {
	panic(wire.Build(
		testAppSet,

		httpserver.Set,
		NewAppServer,
	))
}
