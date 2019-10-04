package server

import "github.com/google/wire"

// Set is the Wire provider set for the server package
// that does not depend on any underlying dependencies.
var Set = wire.NewSet(
	ProvideServer,
)
