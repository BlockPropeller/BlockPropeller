//+build wireinject

package chainup

import (
	"testing"

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

		infrastructure.NewInMemoryServerRepository,
		wire.Bind(new(infrastructure.ServerRepository), new(*infrastructure.InMemoryServerRepository)),

		provision.NewInMemoryJobRepository,
		wire.Bind(new(provision.JobRepository), new(*provision.InMemoryJobRepository)),

		AppSet,
	))
}
