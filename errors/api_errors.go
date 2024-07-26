package api_error

import (
	"errors"
	"fmt"
)

type APIError struct {
	Title   string `json:"title"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Title, e.Message)
}

func NewAPIError(title string, code int, message string) *APIError {
	return &APIError{
		Title:   title,
		Code:    code,
		Message: message,
	}
}

func UnexpectedError(err error) *APIError {
	if err == nil {
		err = errors.New("err: Unexpected Error Occurred")
	}
	return NewAPIError("Unexpected Error", 500, err.Error())
}
