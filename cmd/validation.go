package main

import (
	"github.com/go-playground/validator/v10"
)

func validateStruct(data interface{}) *[]ValidationError {
	validate := validator.New()
	var errors []ValidationError

	err := validate.Struct(data)
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, err := range errs {
				errors = append(errors, ValidationError{
					Field:     err.Field(),
					Tag:       err.Tag(),
					TagTarget: err.Param(),
					Value:     err.Value(),
				})
			}
		}

		return &errors
	}

	return nil
}
