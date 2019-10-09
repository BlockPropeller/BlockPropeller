package routes

import (
	"context"

	"chainup.dev/chainup/account"
	"chainup.dev/chainup/httpserver/request"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// GetAccountResponse is the response for a get account request.
type GetAccountResponse struct {
	Account *account.Account `json:"account"`
}

// Account REST Resource for accessing account information.
type Account struct {
	accRepo account.Repository
}

// NewAccountRoutes returns a new Account REST resource.
func NewAccountRoutes(accRepo account.Repository) *Account {
	return &Account{accRepo: accRepo}
}

// LoadAccount is a middleware for loading Accounts into request context
// as well as checking for correct permissions of an authenticated user.
func (a *Account) LoadAccount(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authAcc := request.AuthFromContext(c)
		if authAcc == nil {
			return echo.ErrForbidden.SetInternal(errors.New("missing authenticated account"))
		}

		accID := account.IDFromString(c.Param("account_id"))
		if accID == "me" {
			accID = authAcc.ID
		}

		if authAcc.ID != accID {
			return echo.ErrForbidden.
				SetInternal(errors.Errorf("unauthorized account access: authenticated %s, account %s",
					authAcc.ID, accID))
		}

		acc, err := a.accRepo.FindByID(context.Background(), accID)
		if err != nil {
			return echo.ErrInternalServerError.SetInternal(err)
		}

		request.WithAccount(c, acc)

		return next(c)
	}
}

// Get an Account.
func (a *Account) Get(c echo.Context) error {
	acc := request.AccountFromContext(c)
	if acc == nil {
		return echo.ErrNotFound
	}

	return c.JSON(200, &GetAccountResponse{Account: acc})
}
