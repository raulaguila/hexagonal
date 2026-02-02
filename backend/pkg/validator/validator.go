package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// ValidateError represents a validation error
type ValidateError struct {
	Field string
	Tag   string
	Param string
	Value any
}

func (e *ValidateError) Error() string {
	return e.Field + " failed on " + e.Tag + " validation"
}

// validatorStruct wraps the validator
type validatorStruct struct {
	validator *validator.Validate
}

// StructValidator is the global validator instance
var StructValidator *validatorStruct

func init() {
	StructValidator = &validatorStruct{
		validator: validator.New(),
	}
}

// Validate validates a struct
func (v validatorStruct) Validate(data any) error {
	if result := v.validator.Struct(data); result != nil {
		var errs validator.ValidationErrors
		if errors.As(result, &errs) {
			for _, err := range errs {
				return &ValidateError{err.Field(), err.Tag(), err.Param(), err.Value()}
			}
		}
	}

	return nil
}
