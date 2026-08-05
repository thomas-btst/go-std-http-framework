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
	errMsgInvalidBody   = "Request body is not valid"
	errMsgInternalError = "Internal server error"

	errFmtInvalidPathInt = "Path parameter '%s' must be a valid integer"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Context struct {
	*http.Request
	Response *Response
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		Request:  r,
		Response: NewResponse(w),
	}
}

func (c *Context) Validate(request any) error {
	return validate.Struct(request)
}

func (c *Context) MustValidate(request any) bool {
	err := c.Validate(request)
	if err != nil {
		if _, ok := errors.AsType[validator.ValidationErrors](err); ok {
			c.Response.WriteError(http.StatusBadRequest, errMsgInvalidBody) // TODO: display a real error
		} else {
			c.Response.WriteInternalError()
		}
		return false
	}

	return true
}

func (c *Context) Body(request any) error {
	return json.NewDecoder(c.Request.Body).Decode(request)
}

func (c *Context) MustBody(request any) bool {
	err := c.Body(request)
	if err != nil {
		c.Response.WriteError(http.StatusBadRequest, errMsgInvalidBody)
		return false
	}

	return true
}

func (c *Context) Bind(request any) error {
	if err := c.Body(request); err != nil {
		return err
	}

	return c.Validate(request)
}

func (c *Context) MustBind(request any) bool {
	return c.MustBody(request) && c.MustValidate(request)
}

func (c *Context) PathString(name string) string {
	return c.PathValue(name)
}

func (c *Context) PathInt(name string) (int, error) {
	strValue := c.PathString(name)
	return strconv.Atoi(strValue)
}

func (c *Context) MustPathInt(name string) (int, bool) {
	value, err := c.PathInt(name)
	if err != nil {
		c.Response.WriteError(http.StatusBadRequest,
			fmt.Sprintf(errFmtInvalidPathInt, name))
		return value, false
	}

	return value, true
}
