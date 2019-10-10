package httpserver

import (
	"chainup.dev/chainup/httpserver/middleware"
	"chainup.dev/chainup/httpserver/routes"
	"github.com/labstack/echo"
)

// Router registers all the routes with the HTTP server.
type Router struct {
	AuthenticatedMiddleware *middleware.AuthenticationMiddleware

	AuthRoutes       *routes.Authentication
	AccountRoutes    *routes.Account
	ProviderSettings *routes.ProviderSettings
}

// RegisterRoutes satisfies the server.Router interface.
func (r *Router) RegisterRoutes(e *echo.Echo) error {
	e.POST("/register", r.AuthRoutes.Register)
	e.POST("/login", r.AuthRoutes.Login)

	protectedAPI := e.Group("/api/v1",
		r.AuthenticatedMiddleware.Middleware)

	protectedAPI.GET("/account/:account_id", r.AccountRoutes.Get,
		r.AccountRoutes.LoadAccount)
	protectedAPI.GET("/provider/types", r.ProviderSettings.GetProviderTypes)
	protectedAPI.GET("/provider/settings", r.ProviderSettings.List)
	protectedAPI.GET("/provider/settings/:settings_id", r.ProviderSettings.Get,
		r.ProviderSettings.LoadProviderSettings)
	protectedAPI.POST("/provider/settings", r.ProviderSettings.Create,
		r.ProviderSettings.LoadProviderSettings)

	return nil
}
