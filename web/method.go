package web

import "net/http"

type HTTPMethod string

const (
	GET     HTTPMethod = http.MethodGet
	POST    HTTPMethod = http.MethodPost
	PUT     HTTPMethod = http.MethodPut
	DELETE  HTTPMethod = http.MethodDelete
	HEAD    HTTPMethod = http.MethodPatch
	OPTIONS HTTPMethod = http.MethodOptions
	PATCH   HTTPMethod = http.MethodPatch
	CONNECT HTTPMethod = http.MethodConnect
	TRACE   HTTPMethod = http.MethodTrace
)
