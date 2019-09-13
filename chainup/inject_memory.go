//+build wireinject

package chainup

import (
	"chainup.dev/chainup/database/transaction"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"chainup.dev/lib/log"
	"github.com/google/wire"
)

// SetupInMemoryApp constructs an in-memory variant of the ChainUP App.
func SetupInMemoryApp() *App {
	panic(wire.Build(
		ProvideFileConfigProvider,

		log.NewConsoleLogger,
		wire.Bind(new(log.Logger), new(*log.ConsoleLogger)),

		transaction.NewInMemoryTransactionContext,
		wire.Bind(new(transaction.TxContext), new(*transaction.InMemoryTxContext)),

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
