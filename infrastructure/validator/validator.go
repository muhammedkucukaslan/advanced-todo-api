package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &Validator{
		validator: validate,
	}
}

func (v *Validator) Validate(data any) error {
	return v.validator.Struct(data)
}
