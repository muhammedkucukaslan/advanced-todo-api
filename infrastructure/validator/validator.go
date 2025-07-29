package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type Validator struct {
	validator *validator.Validate
	logger    domain.Logger
}

func NewValidator(logger domain.Logger) *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &Validator{
		validate,
		logger,
	}
}

func (v *Validator) Validate(data any) error {
	if err := v.validator.Struct(data); err != nil {
		v.logger.Error("Validation error", "err", err.Error())
	}
	return nil
}
