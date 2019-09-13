//+build wireinject

package chainup

import (
	"chainup.dev/chainup/database"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"chainup.dev/lib/log"
	"github.com/google/wire"
)

// SetupDatabaseApp constructs a database backed variant of the ChainUP App.
func SetupDatabaseApp() (*App, func(), error) {
	panic(wire.Build(
		ProvideFileConfigProvider,

		log.NewConsoleLogger,
		wire.Bind(new(log.Logger), new(*log.ConsoleLogger)),

		database.Set,

		database.NewJobRepository,
		wire.Bind(new(provision.JobRepository), new(*database.JobRepository)),

		database.NewServerRepository,
		wire.Bind(new(infrastructure.ServerRepository), new(*database.ServerRepository)),

		database.NewDeploymentRepository,
		wire.Bind(new(infrastructure.DeploymentRepository), new(*database.DeploymentRepository)),

		database.NewProviderSettingsRepository,
		wire.Bind(new(infrastructure.ProviderSettingsRepository), new(*database.ProviderSettingsRepository)),

		AppSet,
	))
}
