//+build wireinject

package chainup

import (
	"chainup.dev/chainup/account"
	"chainup.dev/chainup/database"
	"chainup.dev/chainup/httpserver"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"chainup.dev/lib/log"
	"chainup.dev/lib/server"
	"github.com/google/wire"
)

var dbAppSet = wire.NewSet(
	ProvideFileConfigProvider,

	log.NewConsoleLogger,
	wire.Bind(new(log.Logger), new(*log.ConsoleLogger)),

	database.Set,

	database.NewAccountRepository,
	wire.Bind(new(account.Repository), new(*database.AccountRepository)),

	database.NewJobRepository,
	wire.Bind(new(provision.JobRepository), new(*database.JobRepository)),

	database.NewServerRepository,
	wire.Bind(new(infrastructure.ServerRepository), new(*database.ServerRepository)),

	database.NewDeploymentRepository,
	wire.Bind(new(infrastructure.DeploymentRepository), new(*database.DeploymentRepository)),

	database.NewProviderSettingsRepository,
	wire.Bind(new(infrastructure.ProviderSettingsRepository), new(*database.ProviderSettingsRepository)),

	AppSet,
)

// SetupDatabaseApp constructs a database backed variant of the ChainUP App.
func SetupDatabaseApp() (*App, func(), error) {
	panic(wire.Build(
		dbAppSet,
	))
}

// SetupInMemoryServer constructs a database backed variant of the ChainUP Server.
func SetupDatabaseServer() (*server.Server, func(), error) {
	panic(wire.Build(
		dbAppSet,

		httpserver.Set,
	))
}
