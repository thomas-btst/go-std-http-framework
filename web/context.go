package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

const (
	errMsgInvalidBody = "Request body is not valid"

	errFmtInvalidPathInt = "Path parameter '%s' must be a valid integer"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Context struct {
	*http.Request
	Response *Response
}

func newContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		Request:  r,
		Response: newResponse(w),
	}
}

func (c *Context) Validate(request any) error {
	err := validate.Struct(request)
	if err == nil {
		return nil
	}

	if _, ok := errors.AsType[validator.ValidationErrors](err); ok {
		return NewHTTPErrorWithErr(http.StatusBadRequest,
			errMsgInvalidBody,
			err,
		) // TODO: display a real error
	}

	return err
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
