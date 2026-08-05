package web

import (
	"fmt"
)

type HTTPError struct {
	Code    int
	Message string
	Err     error
}

func NewHTTPErrorWithErr(code int, message string, err error) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
		Err:     err,
	}
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
