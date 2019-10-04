//+build wireinject

package chainup

import (
	"chainup.dev/chainup/account"
	"chainup.dev/chainup/database/transaction"
	"chainup.dev/chainup/httpserver"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"chainup.dev/lib/log"
	"chainup.dev/lib/server"
	"github.com/google/wire"
)

var inMemAppSet = wire.NewSet(
	ProvideFileConfigProvider,

	log.NewConsoleLogger,
	wire.Bind(new(log.Logger), new(*log.ConsoleLogger)),

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

// SetupInMemoryApp constructs an in-memory variant of the ChainUP App.
func SetupInMemoryApp() *App {
	panic(wire.Build(
		inMemAppSet,
	))
}

// SetupInMemoryServer constructs an in-memory backed variant of the ChainUP Server.
func SetupInMemoryServer() (*server.Server, func(), error) {
	panic(wire.Build(
		inMemAppSet,

		httpserver.Set,
	))
}
