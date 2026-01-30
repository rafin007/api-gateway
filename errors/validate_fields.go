package errors

import "github.com/go-playground/validator/v10"

func ValidateFields(err error) map[string]string {
	errorMessages := make(map[string]string)

	if validatorErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validatorErrors {
			field := fieldError.Field()
			switch fieldError.Tag() {
			case "required":
				errorMessages[field] = field + " is required"
			case "email":
				errorMessages[field] = field + " must be a valid email"
			case "min":
				errorMessages[field] = field + " is too short"
			case "max":
				errorMessages[field] = field + " is too long"
			default:
				errorMessages[field] = field + " is invalid"
			}
		}
	}

	return errorMessages
}
