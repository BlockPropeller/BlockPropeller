package httpserver

import (
	"chainup.dev/chainup/httpserver/middleware"
	"chainup.dev/chainup/httpserver/routes"
	"github.com/labstack/echo"
)

// Router registers all the routes with the HTTP server.
type Router struct {
	AuthenticatedMiddleware *middleware.AuthenticationMiddleware

	AuthRoutes    *routes.Authentication
	AccountRoutes *routes.Account
}

// RegisterRoutes satisfies the server.Router interface.
func (r *Router) RegisterRoutes(e *echo.Echo) error {
	e.POST("/register", r.AuthRoutes.Register)
	e.POST("/login", r.AuthRoutes.Login)

	protected := e.Group("/api/v1",
		r.AuthenticatedMiddleware.Middleware)

	protected.GET("/account/:account_id", r.AccountRoutes.Get,
		r.AccountRoutes.LoadAccount)

	return nil
}
