package server

import (
	"github.com/go-playground/validator"
)

type formValidator struct {
	validator *validator.Validate
}

// Validate satisfies the echo.Validator interface.
func (v *formValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
