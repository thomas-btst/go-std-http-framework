package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	http.ResponseWriter
	StatusCode int
}

func newResponse(w http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}
}

func (r *Response) WriteHeader(status int) {
	r.StatusCode = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *Response) Status(status int) error {
	r.WriteHeader(status)
	return nil
}

func (r *Response) NoContent() error {
	return r.Status(http.StatusNoContent)
}

func (r *Response) JSON(status int, data any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return fmt.Errorf("Error during JSON response serialization: %w", err)
	}

	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(status)
	_, err := buf.WriteTo(r)
	if err != nil {
		return fmt.Errorf("Error writing JSON response: %w", err)
	}

	return nil
}

func (r *Response) Error(status int, message string) error {
	return r.JSON(status, map[string]string{"error": message})
}
