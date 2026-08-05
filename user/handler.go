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

func (h *Handler) RegisterRoutes(router *web.Router) {
	router.Group("/users", func(r *web.Router) {
		r.AddRoute(web.GET, "/{id}", h.handleGet)
		r.AddRoute(web.GET, "/", h.handleList)
		r.AddRoute(web.POST, "/", h.handlePost)
		r.AddRoute(web.PUT, "/{id}", h.handleUpdate)
		r.AddRoute(web.DELETE, "/{id}", h.handleDelete)
	})
}

func newErrUserNotFound(id int) error {
	return web.NewHTTPError(
		http.StatusNotFound,
		fmt.Sprintf(errFmtUserIDNotFound, id),
	)
}

func newErrUserNameConflict(name string) error {
	return web.NewHTTPError(
		http.StatusConflict,
		fmt.Sprintf(errFmtUserNameConflict, name),
	)
}

func (h *Handler) handleGet(c *web.Context) error {
	id, err := c.PathInt("id")
	if err != nil {
		return err
	}

	user, err := h.service.Get(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return newErrUserNotFound(id)
		}

		return err
	}

	return c.Response.WriteJSON(http.StatusOK, toResponse(user))
}

func (h *Handler) handleList(c *web.Context) error {
	users := h.service.List(c.Context())

	responses := toResponses(users)

	return c.Response.WriteJSON(http.StatusOK, responses)
}

func (h *Handler) handlePost(c *web.Context) error {
	var request Request
	if err := c.Bind(&request); err != nil {
		return err
	}

	user := toUser(0, &request)

	createdUser, err := h.service.Create(c.Context(), user)
	if err != nil {
		if errors.Is(err, ErrNameConflict) {
			return newErrUserNameConflict(request.Name)
		}

		return err
	}

	return c.Response.WriteJSON(http.StatusCreated, toIDResponse(createdUser.ID))
}

func (h *Handler) handleUpdate(c *web.Context) error {
	id, err := c.PathInt("id")
	if err != nil {
		return err
	}

	var request Request
	if err := c.Bind(&request); err != nil {
		return err
	}

	user := toUser(id, &request)

	_, err = h.service.Update(c.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return newErrUserNotFound(id)
		case errors.Is(err, ErrNameConflict):
			return newErrUserNameConflict(request.Name)
		default:
			return err
		}
	}

	c.Response.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) handleDelete(c *web.Context) error {
	id, err := c.PathInt("id")
	if err != nil {
		return err
	}

	_, err = h.service.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return newErrUserNotFound(id)
		}
		return err
	}

	c.Response.WriteHeader(http.StatusNoContent)

	return nil
}
