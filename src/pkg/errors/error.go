package errors

import (
	"errors"
	"net/http"
)

type PageOutOfBoundError struct {
	RequestedPage int
	TotalPages    int
}

func (e *PageOutOfBoundError) Error() string {
	return "Requested page is out of bounds"
}

var ErrFileNotFound = NewAppError("File not found", http.StatusNotFound)
var ErrRecordNotFound = NewAppError("Record not found", http.StatusNotFound)
var ErrUnauthorized = NewAppError("Unauthorized", http.StatusUnauthorized)
var ErrForbidden = NewAppError("Forbidden", http.StatusForbidden)
var ErrBadRequest = NewAppError("Bad request", http.StatusBadRequest)
var ErrInternalServerError = NewAppError("Internal server error", http.StatusInternalServerError)

type AppError struct {
	Message       string      `json:"message"`
	StatusCode    int         `json:"statusCode"`
	StatusMessage string      `json:"statusMessage"`
	Details       interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// wrapper for errors.Is function (makes it easier to use)
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func NewAppError(message string, statusCode int) *AppError {
	return &AppError{
		Message:       message,
		StatusCode:    statusCode,
		StatusMessage: http.StatusText(statusCode),
	}
}

func (e *AppError) SetMessage(message string) *AppError {
	e.Message = message
	return e
}

func (e *AppError) SetStatusCode(statusCode int) *AppError {
	e.StatusCode = statusCode
	e.StatusMessage = http.StatusText(statusCode)
	return e
}

func (e *AppError) SetStatusMessage(statusMessage string) *AppError {
	e.StatusMessage = statusMessage
	return e
}

func (e *AppError) SetError(err error) *AppError {
	e.Message = err.Error()
	return e
}

func (e *AppError) SetDetails(details interface{}) *AppError {
	e.Details = details
	return e
}

func NewAppErrorFromError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewAppError(err.Error(), http.StatusInternalServerError)
}

func NewAppErrorFromErrorWithCode(err error, code int) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewAppError(err.Error(), code)
}
