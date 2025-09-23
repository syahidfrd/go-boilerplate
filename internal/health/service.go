package health

import (
	"context"
)

// Store defines the interface for health check operations
type Store interface {
	PingDatabase(ctx context.Context) error
	PingCache(ctx context.Context) error
}

// Service provides health check business logic operations
type Service struct {
	store Store
}

// NewService creates a new health service with the provided store
func NewService(store Store) *Service {
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
