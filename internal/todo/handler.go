package todo

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/syahidfrd/go-boilerplate/internal/auth"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/render"
)

// handler handles HTTP requests for todo endpoints
type handler struct {
	svc       *Service
	validator *validator.Validate
}

// NewHandler creates a new todo handler with the provided service
func NewHandler(svc *Service) *handler {
	return &handler{
		svc:       svc,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// Create handles todo creation requests for authenticated users
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		render.JSON(w, http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		return
	}

	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSONFromError(w, err)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		render.JSONFromError(w, err)
		return
	}

	todo, err := h.svc.Create(ctx, userID, &req)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("failed to create todo: %s", err.Error())
		render.JSONFromError(w, err)
		return
	}

	render.JSON(w, http.StatusCreated, todo)
}

// GetByUserID handles requests to retrieve all todos for the authenticated user
func (h *handler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		render.JSON(w, http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		return
	}

	todos, err := h.svc.GetByUserID(ctx, userID)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("failed to get todos: %s", err.Error())
		render.JSONFromError(w, err)
		return
	}

	render.JSON(w, http.StatusOK, map[string]any{"data": todos})
}

// GetByID handles requests to retrieve a specific todo by ID
func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		render.JSONFromError(w, err)
		return
	}

	todo, err := h.svc.GetByID(ctx, int64(id))
	if err != nil {
		switch {
		case errors.Is(err, ErrTodoNotFound):
			render.JSON(w, http.StatusNotFound, map[string]string{"message": "todo not found"})
		default:
			log.Ctx(ctx).Error().Msgf("failed to get todo: %s", err.Error())
			render.JSONFromError(w, err)
		}
		return
	}

	render.JSON(w, http.StatusOK, todo)
}

// Update handles requests to update an existing todo
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		render.JSONFromError(w, err)
		return
	}

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSONFromError(w, err)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		render.JSONFromError(w, err)
		return
	}

	todo, err := h.svc.Update(ctx, int64(id), &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrTodoNotFound):
			render.JSON(w, http.StatusNotFound, map[string]string{"message": "todo not found"})
		default:
			log.Ctx(ctx).Error().Msgf("failed to update todo: %s", err.Error())
			render.JSONFromError(w, err)
		}
		return
	}

	render.JSON(w, http.StatusOK, todo)
}

// ToggleComplete handles requests to toggle the completion status of a todo
func (h *handler) ToggleComplete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		render.JSONFromError(w, err)
		return
	}

	todo, err := h.svc.ToggleComplete(ctx, int64(id))
	if err != nil {
		switch {
		case errors.Is(err, ErrTodoNotFound):
			render.JSON(w, http.StatusNotFound, map[string]string{"message": "todo not found"})
		default:
			log.Ctx(ctx).Error().Msgf("failed to toggle todo: %s", err.Error())
			render.JSONFromError(w, err)
		}
		return
	}

	render.JSON(w, http.StatusOK, todo)
}

// Delete handles requests to delete a todo by ID
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		render.JSONFromError(w, err)
		return
	}

	if err := h.svc.Delete(ctx, int64(id)); err != nil {
		switch {
		case errors.Is(err, ErrTodoNotFound):
			render.JSON(w, http.StatusNotFound, map[string]string{"message": "todo not found"})
		default:
			log.Ctx(ctx).Error().Msgf("failed to delete todo: %s", err.Error())
			render.JSONFromError(w, err)
		}
		return
	}

	render.JSON(w, http.StatusNoContent, nil)
}
