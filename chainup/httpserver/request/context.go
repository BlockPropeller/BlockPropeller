package request

import (
	"chainup.dev/chainup/account"
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
