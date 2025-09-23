package health

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/syahidfrd/go-boilerplate/internal/render"
)

// Handler handles HTTP requests for health check operations
type Handler struct {
	service *Service
}

// NewHandler creates a new health handler with the provided service
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Health handles health check requests
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response, err := h.service.Check(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("health check failed: %s", err.Error())
	}

	statusCode := http.StatusOK
	if response.Database == StatusUnhealthy || response.Cache == StatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	}

	render.JSON(w, statusCode, response)
}
