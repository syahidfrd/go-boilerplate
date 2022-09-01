package http

import (
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"github.com/syahidfrd/go-boilerplate/delivery/middleware"
	"github.com/syahidfrd/go-boilerplate/transport/request"
	"github.com/syahidfrd/go-boilerplate/usecase"
	"github.com/syahidfrd/go-boilerplate/utils"
)

type TodoHandler struct {
	TodoUC usecase.TodoUsecase
}

// NewTodoHandler will initialize the todo resources endpoint
func NewTodoHandler(e *echo.Echo, middleware *middleware.Middleware, todoUC usecase.TodoUsecase) {
	handler := &TodoHandler{
		TodoUC: todoUC,
	}

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/todos", handler.Create, middleware.JWTAuth())
	apiV1.GET("/todos/:id", handler.GetByID, middleware.JWTAuth())
	apiV1.GET("/todos", handler.Fetch, middleware.JWTAuth())
	apiV1.PUT("/todos/:id", handler.Update, middleware.JWTAuth())
	apiV1.DELETE("/todos/:id", handler.Delete, middleware.JWTAuth())
}

// Create godoc
// @Summary Create Todo
// @Description Create Todo
// @Tags Todos
// @Accept json
// @Produce json
// @Param todo body request.CreateTodoReq true "Todo to create"
// @Success 200
// @Router /api/v1/todos [post]
func (h *TodoHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.CreateTodoReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := h.TodoUC.Create(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "todo created",
	})

}

// GetByID godoc
// @Summary Get Todo
// @Description Get Todo
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path string true "todo id"
// @Success 200
// @Router /api/v1/todos/{id} [get]
func (h *TodoHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("todo not found"))
	}

	todo, err := h.TodoUC.GetByID(ctx, int64(id))
	if err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": todo})
}

// Fetch godoc
// @Summary Fetch Todo
// @Description Fetch Todo
// @Tags Todos
// @Accept json
// @Produce json
// @Success 200
// @Router /api/v1/todos [get]
func (h *TodoHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	todos, err := h.TodoUC.Fetch(ctx)
	if err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": todos})
}

// Update godoc
// @Summary Update Todo
// @Description Update Todo
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path string true "todo id"
// @Param todo body request.UpdateTodoReq true "Todo to update"
// @Success 200
// @Router /api/v1/todos/{id} [put]
func (h *TodoHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("todo not found"))
	}

	var req request.UpdateTodoReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := h.TodoUC.Update(ctx, int64(id), &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "todo updated",
	})
}

// Delete godoc
// @Summary Delete Todo
// @Description Delete Todo
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path string true "todo id"
// @Success 200
// @Router /api/v1/todos/{id} [delete]
func (h *TodoHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("todo not found"))
	}

	if err := h.TodoUC.Delete(ctx, int64(id)); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "todo deleted",
	})
}
