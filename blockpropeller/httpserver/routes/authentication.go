package routes

import (
	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/httpserver/request"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// RegisterRequest holds the request payload for the registration endpoint.
type RegisterRequest struct {
	Email    account.Email         `json:"email" form:"email" validate:"required,email"`
	Password account.ClearPassword `json:"password" form:"password" validate:"required,min=6"`
}

// RegisterResponse holds the response returned to the caller upon successful registration.
type RegisterResponse struct {
	Account *account.Account `json:"account"`
	Token   account.Token    `json:"token"`
}

// LoginRequest holds the request payload for the login endpoint.
type LoginRequest struct {
	Email    account.Email         `json:"email" form:"email" validate:"required,email"`
	Password account.ClearPassword `json:"password" form:"password" validate:"required"`
}

// LoginResponse holds the response payload for the login endpoint.
type LoginResponse struct {
	Token account.Token `json:"token"`
}

// Authentication routes define how a user authenticates with the system.
type Authentication struct {
	accSvc *account.Service
}

// NewAuthenticationRoutes returns a new Authentication instance.
func NewAuthenticationRoutes(accSvc *account.Service) *Authentication {
	return &Authentication{accSvc: accSvc}
}

// Register an account with BlockPropeller.
func (a *Authentication) Register(c echo.Context) error {
	var req RegisterRequest
	if err := request.Parse(c, &req); err != nil {
		return err
	}

	acc, token, err := a.accSvc.Register(req.Email, req.Password)
	if errors.Cause(err) == account.ErrEmailAlreadyExists {
		return echo.ErrBadRequest.SetInternal(err)
	}
	if err != nil {
		return errors.Wrap(err, "register account")
	}

	return c.JSON(201, &RegisterResponse{
		Account: acc,
		Token:   token,
	})
}

// Login to an account with BlockPropeller.
func (a *Authentication) Login(c echo.Context) error {
	var req LoginRequest
	if err := request.Parse(c, &req); err != nil {
		return err
	}

	token, err := a.accSvc.Login(req.Email, req.Password)
	if err != nil {
		return echo.ErrForbidden.SetInternal(errors.Wrap(err, "login account"))
	}

	return c.JSON(200, &LoginResponse{Token: token})
}
