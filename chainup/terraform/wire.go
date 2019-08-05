package terraform

import "github.com/google/wire"

// Set is the Wire provider set for the provisioning package
// that does not depend on any underlying dependencies.
var Set = wire.NewSet(
	ConfigureTerraform,
)
