package provision

import "github.com/google/wire"

// Set is the Wire provider set for the Provisioner
// that does not depend on any underlying dependencies.
var Set = wire.NewSet(
	NewProvisioner,
	ConfigureStateMachine,
)
