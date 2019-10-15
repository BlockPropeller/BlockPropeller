package server

import (
	"blockpropeller.dev/lib/log"
	"github.com/go-playground/validator"
)

type requestValidator struct {
	validator *validator.Validate
}

func newRequestValidator() *requestValidator {
	v := validator.New()
	_ = v.RegisterValidation("valid", validateIsValid)

	return &requestValidator{validator: v}
}

// Validate satisfies the echo.Validator interface.
func (v *requestValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func validateIsValid(fl validator.FieldLevel) bool {
	type validatable interface {
		IsValid() bool
	}

	prop, ok := fl.Field().Interface().(validatable)
	if !ok {
		log.Warn("field does not have an IsValid() method", log.Fields{
			"field_name": fl.FieldName(),
			"param":      fl.Param(),
		})

		return false
	}

	return prop.IsValid()
}
