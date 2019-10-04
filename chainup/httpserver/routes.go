package httpserver

import (
	"net/http"

	"github.com/labstack/echo"
)

// Router registers all the routes with the HTTP server.
type Router struct {
}

// RegisterRoutes satisfies the server.Router interface.
func (r *Router) RegisterRoutes(e *echo.Echo) error {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Injections")
	})

	return nil
}
