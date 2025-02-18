package routes

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

const (
	parametrizedMessage   = "unexpected input value in field '%s'. expected value matching: %s(%s), but got: %v"
	unparametrizedMessage = "unexpected input value in field '%s'. expected value matching: %s, but got: %v"
)

func GenerateErrors(validationErrors validator.ValidationErrors) []error {
	errs := make([]error, len(validationErrors))
	for _, err := range validationErrors {
		if err.Param() != "" {
			errs = append(errs, fmt.Errorf(parametrizedMessage, err.Field(), err.ActualTag(), err.Param(), err.Value()))
			continue
		}
		errs = append(errs, fmt.Errorf(unparametrizedMessage, err.Field(), err.ActualTag(), err.Value()))
	}
	return errs
}
