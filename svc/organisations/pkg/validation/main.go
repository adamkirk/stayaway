package validation

import (
	"errors"

	"gopkg.in/validator.v2"
)

type FieldError struct {
	Key string
	Errors []string
}

type ValidationError struct {
	Errs []FieldError
}

func (err ValidationError) Error() string {
	return "invalid data"
}

type Validator struct {}

func (v *Validator) Validate(in any) error {
	err := validator.Validate(in)

	if err == nil {
		return nil
	}

	errMap, ok := err.(validator.ErrorMap)

	if ! ok {
		return errors.New("unable to handle validation")
	}

	fieldErrors := []FieldError{}

	for key, errArray := range(errMap) {
		messages := []string{}
		for _, fieldError := range(errArray) {
			messages = append(messages, fieldError.Error())
		}

		fieldErrors = append(fieldErrors, FieldError{
			Key: key,
			Errors: messages,
		})
	}

	return ValidationError{
		Errs: fieldErrors,
	}
}

func NewValidator() *Validator {
	return &Validator{}
}

