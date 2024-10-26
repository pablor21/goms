package errors

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ValidationError struct {
	*AppError
	ValidationErrors error `json:"validationErrors"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

type FieldErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Params  map[string]interface{} `json:"params"`
}

func (e *FieldErrorDetail) Error() string {
	return e.Message
}

type FieldError struct {
	Key   string `json:"index,omitempty"`
	Field string `json:"field"`
	Error error  `json:"error"`
}

type FieldErrors []FieldError

func (e FieldErrors) Error() string {
	return "Validation error"
}

func ParseFieldErrors(err error, index string) error {
	// ferrors := []FieldError{}

	if verrs, ok := err.(*FieldErrorDetail); ok {
		return verrs
	}

	if verrs, ok := err.(validation.Errors); ok {
		for key, verr := range verrs {

			return ParseValidationErrors(verr, key)

		}
	}

	if verrs, ok := err.(validation.ErrorObject); ok {
		return &FieldErrorDetail{
			Code:    verrs.Code(),
			Message: verrs.Error(),
			Params:  verrs.Params(),
		}
	}
	return &FieldErrorDetail{
		Code:    "unknown",
		Message: err.Error(),
		Params:  map[string]interface{}{},
	}
}

func ParseValidationErrors(err error, key string) error {
	verrors := FieldErrors{}

	if verrs, ok := err.(validation.Errors); ok {
		for field, verr := range verrs {
			verrors = append(verrors, FieldError{
				Key:   key,
				Field: field,
				Error: ParseFieldErrors(verr, key),
			})
		}
	}
	return verrors
}

func NewValidationErrorWithMessage(message string, err error) *ValidationError {
	return &ValidationError{
		AppError:         NewAppError(message, http.StatusUnprocessableEntity),
		ValidationErrors: ParseValidationErrors(err, ""),
	}
}

func NewValidationError(err error) *ValidationError {
	return NewValidationErrorWithMessage("Validation error", err)
}

func NewUniqueValidationErrorDetail(params map[string]interface{}) error {
	return &FieldErrorDetail{
		Code:    "validation_unique",
		Message: "value already exists",
		Params:  params,
	}
}
