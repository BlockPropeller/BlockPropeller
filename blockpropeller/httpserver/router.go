package httpserver

import (
	"blockpropeller.dev/blockpropeller/httpserver/middleware"
	"blockpropeller.dev/blockpropeller/httpserver/routes"
	"github.com/labstack/echo"
)

// Router registers all the routes with the HTTP server.
type Router struct {
	AuthenticatedMiddleware *middleware.AuthenticationMiddleware

	AuthRoutes             *routes.Authentication
	AccountRoutes          *routes.Account
	ProviderSettingsRoutes *routes.ProviderSettings
	ProvisionRoutes        *routes.Provision
	ServerRoutes           *routes.Server
}

// RegisterRoutes satisfies the server.Router interface.
func (r *Router) RegisterRoutes(e *echo.Echo) error {
	e.POST("/register", r.AuthRoutes.Register)
	e.POST("/login", r.AuthRoutes.Login)

	protectedAPI := e.Group("/api/v1",
		r.AuthenticatedMiddleware.Middleware)

	protectedAPI.GET("/account/:account_id", r.AccountRoutes.Get,
		r.AccountRoutes.LoadAccount)

	protectedAPI.GET("/provider/types", r.ProviderSettingsRoutes.GetProviderTypes)
	protectedAPI.GET("/provider/settings", r.ProviderSettingsRoutes.List)
	protectedAPI.GET("/provider/settings/:settings_id", r.ProviderSettingsRoutes.Get,
		r.ProviderSettingsRoutes.LoadProviderSettings)
	protectedAPI.POST("/provider/settings", r.ProviderSettingsRoutes.Create)

	protectedAPI.POST("/provision/job", r.ProvisionRoutes.CreateJob)
	protectedAPI.GET("/provision/job/:job_id", r.ProvisionRoutes.GetJob,
		r.ProvisionRoutes.LoadJob)

	protectedAPI.GET("/server", r.ServerRoutes.List)
	protectedAPI.GET("/server/:server_id", r.ServerRoutes.Get,
		r.ServerRoutes.LoadServer)
	protectedAPI.DELETE("/server/:server_id", r.ServerRoutes.Delete,
		r.ServerRoutes.LoadServer)

	return nil
}
