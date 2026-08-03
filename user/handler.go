package user

import (
	"errors"
	"fmt"
	"net/http"

	"standard/web"
)

const (
	errFmtUserNameConflict = "User with name '%s' already exists"
	errFmtUserIDNotFound   = "User with ID '%d' does not exist"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r *web.Router) {
	r.AddRoute(web.GET, "/users/{id}", h.handleGet)
	r.AddRoute(web.GET, "/users", h.handleList)
	r.AddRoute(web.POST, "/users", h.handlePost)
	r.AddRoute(web.PUT, "/users/{id}", h.handleUpdate)
	r.AddRoute(web.DELETE, "/users/{id}", h.handleDelete)
}

func writeUserNotFoundError(c *web.Context, id int) {
	c.Response.WriteError(
		http.StatusNotFound,
		fmt.Sprintf(errFmtUserIDNotFound, id),
	)
}

func writeUserNameConflictError(c *web.Context, name string) {
	c.Response.WriteError(
		http.StatusConflict,
		fmt.Sprintf(errFmtUserNameConflict, name),
	)
}

func (h *Handler) handleGet(c *web.Context) {
	id, ok := c.MustPathInt("id")
	if !ok {
		return
	}

	user, err := h.service.Get(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeUserNotFoundError(c, id)
		} else {
			c.Response.WriteInternalError()
		}
		return
	}

	c.Response.WriteJSON(http.StatusOK, toResponse(user))
}

func (h *Handler) handleList(c *web.Context) {
	users := h.service.List(c.Context())

	responses := toResponses(users)

	c.Response.WriteJSON(http.StatusOK, responses)
}

func (h *Handler) handlePost(c *web.Context) {
	var request Request
	if !c.MustBind(&request) {
		return
	}

	user := toUser(0, &request)

	createdUser, err := h.service.Create(c.Context(), user)
	if err != nil {
		if errors.Is(err, ErrNameConflict) {
			writeUserNameConflictError(c, request.Name)
		} else {
			c.Response.WriteInternalError()
		}
		return
	}

	c.Response.WriteJSON(http.StatusCreated, toIDResponse(createdUser.ID))
}

func (h *Handler) handleUpdate(c *web.Context) {
	id, ok := c.MustPathInt("id")
	if !ok {
		return
	}

	var request Request
	if !c.MustBind(&request) {
		return
	}

	user := toUser(id, &request)

	_, err := h.service.Update(c.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			writeUserNotFoundError(c, id)
		case errors.Is(err, ErrNameConflict):
			writeUserNameConflictError(c, request.Name)
		default:
			c.Response.WriteInternalError()
		}
		return
	}
	c.Response.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleDelete(c *web.Context) {
	id, ok := c.MustPathInt("id")
	if !ok {
		return
	}

	_, err := h.service.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeUserNotFoundError(c, id)
		} else {
			c.Response.WriteInternalError()
		}
		return
	}

	c.Response.WriteHeader(http.StatusNoContent)
}
