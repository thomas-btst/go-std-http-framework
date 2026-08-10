package web

import (
	"fmt"
)

type HTTPError struct {
	Code    int    `json:"status"`
	Message string `json:"message"`
	Err     error  `json:"-"`
	Details any    `json:"details,omitempty"`
}

func NewHTTPErrorWithDetails(code int, message string, err error, details any) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
		Err:     err,
		Details: details,
	}
}

func NewHTTPErrorWithErr(code int, message string, err error) *HTTPError {
	return NewHTTPErrorWithDetails(code, message, err, nil)
}

func NewHTTPError(code int, message string) *HTTPError {
	return NewHTTPErrorWithErr(code, message, nil)
}

func (e *HTTPError) Error() string {
	str := fmt.Sprintf("HTTP error %d: %s", e.Code, e.Message)

	if e.Err != nil {
		return fmt.Sprintf("%s (err: %v)", str, e.Err)
	}

	return str
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}
