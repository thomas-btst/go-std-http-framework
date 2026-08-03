package web

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

func (r *Response) WriteJSON(status int, data any) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		slog.Error("Error during JSON response serialization", slog.Any("err", err))
		return
	}

	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(status)
	_, err := buf.WriteTo(r)
	if err != nil {
		slog.Error("Error writing JSON response", slog.Any("err", err))
	}
}

func (r *Response) WriteError(status int, message string) {
	r.WriteJSON(status, map[string]string{"error": message})
}

func (r *Response) WriteInternalError() {
	r.WriteError(http.StatusInternalServerError, errMsgInternalError)
}
