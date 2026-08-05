package web

import "net/http"

type HandlerFunc func(*Context)

func adapter(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context := newContext(r, w)
		handler(context)
	}
}
