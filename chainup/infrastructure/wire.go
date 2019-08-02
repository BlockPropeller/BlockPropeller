package infrastructure

import "github.com/google/wire"

// Set is the Wire provider set for infrastructure package
// that does not depend on any underlying dependencies.
var Set = wire.NewSet(
	ConfigureServerStateMachine,
)
