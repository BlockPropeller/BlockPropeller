//+build wireinject

package chainup

import (
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"chainup.dev/lib/log"
	"github.com/google/wire"
)

// SetupInMemoryApp constructs an in-memory variant of the StateMachine handling Server state transitions.
func SetupInMemoryApp() *App {
	panic(wire.Build(
		ProvideFileConfigProvider,

		log.NewConsoleLogger,
		wire.Bind(new(log.Logger), new(*log.ConsoleLogger)),

		infrastructure.NewInMemoryServerRepository,
		wire.Bind(new(infrastructure.ServerRepository), new(*infrastructure.InMemoryServerRepository)),

		provision.NewInMemoryJobRepository,
		wire.Bind(new(provision.JobRepository), new(*provision.InMemoryJobRepository)),

		AppSet,
	))
}
