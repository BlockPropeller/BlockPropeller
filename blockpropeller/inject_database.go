//+build wireinject

package blockpropeller

import (
	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/database"
	"blockpropeller.dev/blockpropeller/httpserver"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/provision"
	"blockpropeller.dev/lib/log"
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

// SetupDatabaseApp constructs a database backed variant of the BlockPropeller App.
func SetupDatabaseApp() (*App, func(), error) {
	panic(wire.Build(
		dbAppSet,
	))
}

// SetupInMemoryServer constructs a database backed variant of the BlockPropeller Server.
func SetupDatabaseServer() (*AppServer, func(), error) {
	panic(wire.Build(
		dbAppSet,

		httpserver.Set,
		NewAppServer,
	))
}
