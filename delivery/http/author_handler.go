package http

import (
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo"
	appMiddleware "github.com/syahidfrd/go-boilerplate/middleware"
	"github.com/syahidfrd/go-boilerplate/transport/request"
	"github.com/syahidfrd/go-boilerplate/usecase"
	"github.com/syahidfrd/go-boilerplate/utils"
)

type AuthorHandler struct {
	AuthorUC usecase.AuthorUsecase
}

// NewAuthorHandler will initialize the authors/ resources endpoint
func NewAuthorHandler(e *echo.Echo, middManager *appMiddleware.MiddlewareManager, authorUC usecase.AuthorUsecase) {
	handler := &AuthorHandler{
		AuthorUC: authorUC,
	}

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/authors", handler.Create)
	apiV1.GET("/authors/:id", handler.GetByID)
	apiV1.GET("/authors", handler.Fetch)
	apiV1.PUT("/authors/:id", handler.Update)
	apiV1.DELETE("/authors/:id", handler.Delete)
}

// Create will store the author by given request body
func (h *AuthorHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.CreateAuthorReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := h.AuthorUC.Create(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
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
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("author not found"))
	}

	author, err := h.AuthorUC.GetByID(ctx, int64(id))
	if err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": author})
}

// Fetch will fetch the author
func (h *AuthorHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	authors, err := h.AuthorUC.Fetch(ctx)
	if err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": authors})
}

// Update will get author by given request body
func (h *AuthorHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("author not found"))
	}

	var req request.UpdateAuthorReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := h.AuthorUC.Update(ctx, int64(id), &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
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
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("author not found"))
	}

	if err := h.AuthorUC.Delete(ctx, int64(id)); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "author deleted",
	})
}
