package server

import (
	"fmt"
	"net/http"

	"blockpropeller.dev/lib/log"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func httpErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if he, ok := errors.Cause(err).(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
		if he.Internal != nil {
			c.Set("_internal_error", he.Internal)
			err = fmt.Errorf("%v, %v", err, he.Internal)
		}
	} else if ve, ok := errors.Cause(err).(validator.ValidationErrors); ok {
		code = http.StatusBadRequest
		msg = http.StatusText(http.StatusBadRequest)
		if len(ve) > 0 {
			msg = fmt.Sprintf("Validation failed for '%s' on tag '%s'",
				ve[0].Field(), ve[0].Tag())
		}
	} else {
		msg = http.StatusText(code)
	}
	if _, ok := msg.(string); ok {
		msg = echo.Map{"message": msg}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			log.ErrorErr(err, "send error response")
		}
	}
}
