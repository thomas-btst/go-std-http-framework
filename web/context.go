package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	errMsgInvalidBody = "Request body is not valid"

	errFmtInvalidPathInt = "Path parameter '%s' must be a valid integer"
)

type Context struct {
	*http.Request
	Response  *Response
	validator Validator
}

func newContext(r *http.Request, w http.ResponseWriter, validator Validator) *Context {
	return &Context{
		Request:   r,
		Response:  newResponse(w),
		validator: validator,
	}
}

func (c *Context) Validate(request any) error {
	if c.validator == nil {
		return nil
	}
	return c.validator.Validate(request)
}

func (c *Context) Body(request any) error {
	err := json.NewDecoder(c.Request.Body).Decode(request)
	if err != nil {
		return NewHTTPErrorWithErr(
			http.StatusBadRequest,
			errMsgInvalidBody,
			err,
		)
	}

	return nil
}

func (c *Context) Bind(request any) error {
	if err := c.Body(request); err != nil {
		return err
	}

	return c.Validate(request)
}

func (c *Context) PathString(name string) string {
	return c.PathValue(name)
}

func (c *Context) PathInt(name string) (int, error) {
	strVal := c.PathString(name)

	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		return 0, NewHTTPErrorWithErr(
			http.StatusBadRequest,
			fmt.Sprintf(errFmtInvalidPathInt, name),
			err,
		)
	}

	return intVal, nil
}
