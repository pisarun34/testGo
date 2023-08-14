package handlers

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func handleValidationErrors(err error, translator *ut.Translator) map[string]string {
	validationErrors := err.(validator.ValidationErrors)
	errorMessages := make(map[string]string)
	for _, e := range validationErrors {
		fmt.Println(e.Tag())
		errorMessages[e.Field()] = e.Translate(*translator)
	}
	return errorMessages
}
