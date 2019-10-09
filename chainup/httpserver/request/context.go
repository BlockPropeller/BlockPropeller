package request

import (
	"chainup.dev/chainup/account"
	"chainup.dev/chainup/infrastructure"
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
