//+build wireinject

package chainup

import (
	"chainup.dev/chainup/infrastructure"
	"github.com/google/wire"
)

// SetupTestApp constructs an in-memory variant of the StateMachine handling Server state transitions.
func SetupTestApp() *App {
	panic(wire.Build(
		ProvideTestConfigProvider,

		infrastructure.NewInMemoryServerRepository,
		wire.Bind(new(infrastructure.ServerRepository), new(infrastructure.InMemoryServerRepository)),

		//provision.NewInMemoryJobRepository,
		//wire.Bind(new(provision.JobRepository), new(provision.InMemoryJobRepository)),

		AppSet,
	))
}
