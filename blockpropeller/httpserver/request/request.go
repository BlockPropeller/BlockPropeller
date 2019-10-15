package request

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Parse the request into a struct.
func Parse(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return errors.Wrap(err, "bind request")
	}
	if err := c.Validate(req); err != nil {
		return errors.Wrap(err, "validate request")
	}

	return nil
}
