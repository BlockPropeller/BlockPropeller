package routes

import "github.com/google/wire"

// Set is the Wire provider set for the routes package
// that does not depend on any underlying dependencies.
var Set = wire.NewSet(
	NewAuthenticationRoutes,
	NewAccountRoutes,
	NewProviderSettingsRoutes,
)
