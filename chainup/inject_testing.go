//+build wireinject

package chainup

import (
	"testing"

	"chainup.dev/chainup/account"
	"chainup.dev/chainup/database/transaction"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"chainup.dev/lib/log"
	"github.com/google/wire"
)

// SetupTestApp constructs an in-memory variant of the StateMachine handling Server state transitions.
func SetupTestApp(t *testing.T) *App {
	panic(wire.Build(
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
	))
}
