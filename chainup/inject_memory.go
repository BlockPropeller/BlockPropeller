//+build wireinject

package chainup

import (
	"chainup.dev/chainup/provision"
	"github.com/google/wire"
)

// SetupInMemoryProvisioner constructs an in-memory variant of the Provisioner.
func SetupInMemoryProvisioner() *provision.Provisioner {
	panic(wire.Build(
		provision.Set,

		provision.NewInMemoryJobRepository,
		wire.Bind(new(provision.JobRepository), new(provision.InMemoryJobRepository)),
	))
}
