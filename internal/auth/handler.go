package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/render"
)

// handler handles HTTP requests for authentication endpoints
type handler struct {
	svc       *Service
	validator *validator.Validate
}

// NewHandler creates a new auth handler with the provided service
func NewHandler(svc *Service) *handler {
	return &handler{
		svc:       svc,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// SignIn handles user authentication requests
func (h *handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSONFromError(w, err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		render.JSONFromError(w, err)
		return
	}

	ctx := r.Context()
	resp, err := h.svc.SignIn(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			render.JSON(w, http.StatusUnauthorized, map[string]string{"message": "invalid credentials"})
		default:
			log.Ctx(ctx).Error().Msgf("failed to signin: %s", err.Error())
			render.JSONFromError(w, err)
		}
		return
	}

	render.JSON(w, http.StatusOK, resp)
}

// SignUp handles user registration requests
func (h *handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSONFromError(w, err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		render.JSONFromError(w, err)
		return
	}

	ctx := r.Context()
	resp, err := h.svc.SignUp(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserAlreadyExists):
			render.JSON(w, http.StatusBadRequest, map[string]string{"message": "user already exists"})
		default:
			log.Ctx(ctx).Error().Msgf("failed to signup: %s", err.Error())
			render.JSONFromError(w, err)
		}
		return
	}

	render.JSON(w, http.StatusOK, resp)
}
