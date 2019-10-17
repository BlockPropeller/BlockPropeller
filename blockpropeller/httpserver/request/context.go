package request

import (
	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/provision"
	"github.com/labstack/echo"
)

// WithAuth adds an authenticated Account to echo.Context.
func WithAuth(c echo.Context, acc *account.Account) {
	c.Set("_auth", acc)
}

// AuthFromContext returns an Account from echo.Context.
func AuthFromContext(c echo.Context) *account.Account {
	acc := c.Get("_auth")
	if acc == nil {
		return nil
	}

	return acc.(*account.Account)
}

// WithAccount adds an Account resource to echo.Context.
func WithAccount(c echo.Context, acc *account.Account) {
	c.Set("_account", acc)
}

// AccountFromContext returns an Account from echo.Context.
func AccountFromContext(c echo.Context) *account.Account {
	acc := c.Get("_account")
	if acc == nil {
		return nil
	}

	return acc.(*account.Account)
}

// WithProviderSettings adds an ProviderSettings resource to echo.Context.
func WithProviderSettings(c echo.Context, settings *infrastructure.ProviderSettings) {
	c.Set("_provider_settings", settings)
}

// ProviderSettingsFromContext returns an ProviderSettings from echo.Context.
func ProviderSettingsFromContext(c echo.Context) *infrastructure.ProviderSettings {
	settings := c.Get("_provider_settings")
	if settings == nil {
		return nil
	}

	return settings.(*infrastructure.ProviderSettings)
}

// WithJob adds an Job resource to echo.Context.
func WithJob(c echo.Context, job *provision.Job) {
	c.Set("_job", job)
}

// JobFromContext returns an Job from echo.Context.
func JobFromContext(c echo.Context) *provision.Job {
	job := c.Get("_job")
	if job == nil {
		return nil
	}

	return job.(*provision.Job)
}

// WithServer adds an Server resource to echo.Context.
func WithServer(c echo.Context, srv *infrastructure.Server) {
	c.Set("_server", srv)
}

// ServerFromContext returns an Server from echo.Context.
func ServerFromContext(c echo.Context) *infrastructure.Server {
	srv := c.Get("_server")
	if srv == nil {
		return nil
	}

	return srv.(*infrastructure.Server)
}
