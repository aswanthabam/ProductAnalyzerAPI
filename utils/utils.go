package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) interface{} {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return gin.H{"error": "Invalid input."}
	}
	errorMap := map[string][]string{
		"required": {},
		"email":    {},
		"invalid":  {},
	}
	for _, fieldError := range validationErrors {
		fieldName := fieldError.Field()
		switch fieldError.Tag() {
		case "required":
			errorMap["required"] = append(errorMap["required"], fieldName)
		case "email":
			errorMap["email"] = append(errorMap["email"], fieldName)
		default:
			errorMap["invalid"] = append(errorMap["invalid"], fieldName)
		}
	}

	return errorMap
}
