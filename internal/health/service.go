package health

import (
	"context"
)

// Service provides health check business logic operations
type Service struct {
	store *store
}

// NewService creates a new health service with the provided store
func NewService(store *store) *Service {
	return &Service{
		store: store,
	}
}

// Check performs health checks on all components and returns the overall status
func (s *Service) Check(ctx context.Context) (*HealthResponse, error) {
	response := &HealthResponse{
		Database: StatusHealthy,
		Cache:    StatusHealthy,
	}

	// Check database
	if err := s.store.PingDatabase(ctx); err != nil {
		response.Database = StatusUnhealthy
	}

	// Check cache
	if err := s.store.PingCache(ctx); err != nil {
		response.Cache = StatusUnhealthy
	}

	return response, nil
}
