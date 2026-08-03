package web

import "net/http"

type HTTPMethod string

const (
	GET    HTTPMethod = http.MethodGet
	POST   HTTPMethod = http.MethodPost
	PUT    HTTPMethod = http.MethodPut
	DELETE HTTPMethod = http.MethodDelete
	PATCH  HTTPMethod = http.MethodPatch
)
