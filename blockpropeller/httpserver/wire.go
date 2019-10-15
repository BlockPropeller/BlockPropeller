package httpserver

import (
	"blockpropeller.dev/blockpropeller/httpserver/middleware"
	"blockpropeller.dev/blockpropeller/httpserver/routes"
	"blockpropeller.dev/lib/server"
	"github.com/google/wire"
)

// Set is the Wire provider set for the httpserver package
// that does not depend on any underlying dependencies.
var Set = wire.NewSet(
	wire.Struct(new(Router), "*"),
	wire.Bind(new(server.Router), new(*Router)),

	middleware.NewAuthenticationMiddleware,

	routes.Set,
	server.Set,
)
