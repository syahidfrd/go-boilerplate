package http

import (
	"database/sql"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo"
	"github.com/syahidfrd/go-boilerplate/transport/request"
	"github.com/syahidfrd/go-boilerplate/usecase"
	"github.com/syahidfrd/go-boilerplate/utils"
)

type AuthorHandler struct {
	AuthorUsecase usecase.AuthorUsecase
}

// NewAuthorHandler will initialize the authors/ resources endpoint
func NewAuthorHandler(e *echo.Group, authorUsecase usecase.AuthorUsecase) {
	handler := &AuthorHandler{
		AuthorUsecase: authorUsecase,
	}

	e.POST("/authors", handler.Create)
	e.GET("/authors/:id", handler.GetByID)
	e.GET("/authors", handler.Fetch)
	e.PUT("/authors/:id", handler.Update)
	e.DELETE("/authors/:id", handler.Delete)
}

// Create will store the author by given request body
func (h *AuthorHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.CreateAuthorReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err := req.Validate(); err != nil {
		errValidations := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewValidationError(errValidations))
	}

	if err := h.AuthorUsecase.Create(ctx, &req); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "author created",
	})

}

// GetByID will get author by given id
func (h *AuthorHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError())
	}

	author, err := h.AuthorUsecase.GetByID(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, utils.NewNotFoundError())
		}

		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": author})
}

// Fetch will fetch the author
func (h *AuthorHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	authors, err := h.AuthorUsecase.Fetch(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": authors})
}

// Update will get author by given request body
func (h *AuthorHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError())
	}

	var req request.UpdateAuthorReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err := req.Validate(); err != nil {
		errValidations := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewValidationError(errValidations))
	}

	if err := h.AuthorUsecase.Update(ctx, int64(id), &req); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, utils.NewNotFoundError())
		}

		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "author updated",
	})
}

// Delete will delete author by given param
func (h *AuthorHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError())
	}

	if err := h.AuthorUsecase.Delete(ctx, int64(id)); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, utils.NewNotFoundError())
		}

		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "author deleted",
	})
}
