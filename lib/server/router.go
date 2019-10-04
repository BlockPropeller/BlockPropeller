package server

import "github.com/labstack/echo"

// Router is the interface applications are supposed to provide in order
// to register all the required routes to the echo.Echo HTTP server.
type Router interface {
	RegisterRoutes(*echo.Echo) error
}
