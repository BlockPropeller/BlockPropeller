package middleware

import (
	"chainup.dev/chainup/account"
	"chainup.dev/chainup/httpserver/request"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// AuthenticationMiddleware ensures that a request is properly authenticated.
type AuthenticationMiddleware struct {
	accSvc *account.Service
}

// NewAuthenticationMiddleware returns a new AuthenticationMiddleware.
func NewAuthenticationMiddleware(accSvc *account.Service) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{accSvc: accSvc}
}

// Middleware satisfies the echo.MiddlewareFunc interface.
func (s *AuthenticationMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		prefixLen := len("Bearer ")
		if len(authHeader) <= prefixLen {
			return echo.ErrUnauthorized.
				SetInternal(errors.New("missing authorization token"))
		}

		token := account.NewToken(authHeader[prefixLen:])

		acc, err := s.accSvc.Authenticate(token)
		if err != nil {
			return echo.ErrUnauthorized.
				SetInternal(errors.Wrap(err, "authenticate token"))
		}

		request.WithAuth(c, acc)

		return next(c)
	}
}
