package web

import "net/http"

type HandlerFunc func(*Context)

func adapter(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := &Response{ResponseWriter: w}
		context := &Context{
			Request:  r,
			Response: response,
		}
		handler(context)
	}
}
