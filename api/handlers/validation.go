package handlers

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func handleValidationErrors(err error, trans *ut.Translator) []string {
	var errorMessages []string
	for _, e := range err.(validator.ValidationErrors) {
		// Translate each error one at a time
		errorMessages = append(errorMessages, e.Translate(*trans))
	}
	return errorMessages
}
