package main

import (
	modelHTTP "github.com/Kosfedev/auth/pkg/user_v1"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func validateStruct(data interface{}) *[]modelHTTP.ValidationError {
	var errors []modelHTTP.ValidationError

	err := validate.Struct(data)
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, err := range errs {
				errors = append(errors, modelHTTP.ValidationError{
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
