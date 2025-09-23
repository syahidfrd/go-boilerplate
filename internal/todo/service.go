package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/syahidfrd/go-boilerplate/internal/pkg/cache"
)

// Service provides todo business logic operations with caching support
type Service struct {
	store *store
	cache *cache.RedisCache
}

// CreateTodoRequest represents the request payload for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

// UpdateTodoRequest represents the request payload for updating a todo
type UpdateTodoRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

// NewService creates a new todo service with the provided dependencies
func NewService(store *store, cache *cache.RedisCache) *Service {
	return &Service{
		store: store,
		cache: cache,
	}
}

// Create creates a new todo item for the specified user
func (s *Service) Create(ctx context.Context, userID int64, req *CreateTodoRequest) (*Todo, error) {
	todo := NewTodo(userID, req.Title, req.Description)

	if err := s.store.Save(ctx, todo); err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	// Invalidate user's todo cache
	cacheKey := fmt.Sprintf("todos:user:%d", userID)
	s.cache.Delete(ctx, cacheKey)

	return todo, nil
}

// GetByID retrieves a todo by its ID
func (s *Service) GetByID(ctx context.Context, id int64) (*Todo, error) {
	todo, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo by id: %w", err)
	}
	return todo, nil
}

// GetByUserID retrieves all todos for a specific user with caching support
func (s *Service) GetByUserID(ctx context.Context, userID int64) ([]Todo, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("todos:user:%d", userID)

	if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
		var todos []Todo
		if json.Unmarshal([]byte(cached), &todos) == nil {
			return todos, nil
		}
	}

	// Cache miss, get from database
	todos, err := s.store.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos by user id: %w", err)
	}

	// Cache the result
	if data, err := json.Marshal(todos); err == nil {
		s.cache.Set(ctx, cacheKey, string(data), 10*time.Minute)
	}

	return todos, nil
}

// Update updates an existing todo with new title and description
func (s *Service) Update(ctx context.Context, id int64, req *UpdateTodoRequest) (*Todo, error) {
	todo, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo for update: %w", err)
	}

	todo.Title = req.Title
	todo.Description = req.Description

	if err := s.store.Save(ctx, todo); err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	// Invalidate user's todo cache
	cacheKey := fmt.Sprintf("todos:user:%d", todo.UserID)
	s.cache.Delete(ctx, cacheKey)

	return todo, nil
}

// ToggleComplete toggles the completion status of a todo
func (s *Service) ToggleComplete(ctx context.Context, id int64) (*Todo, error) {
	todo, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo for toggle: %w", err)
	}

	if todo.Completed {
		todo.MarkAsIncomplete()
	} else {
		todo.MarkAsCompleted()
	}

	if err := s.store.Save(ctx, todo); err != nil {
		return nil, fmt.Errorf("failed to toggle todo completion: %w", err)
	}

	// Invalidate user's todo cache
	cacheKey := fmt.Sprintf("todos:user:%d", todo.UserID)
	s.cache.Delete(ctx, cacheKey)

	return todo, nil
}

// Delete removes a todo by its ID
func (s *Service) Delete(ctx context.Context, id int64) error {
	// Get todo first to get UserID for cache invalidation
	todo, err := s.store.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get todo for delete: %w", err)
	}

	if err := s.store.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	// Invalidate user's todo cache
	cacheKey := fmt.Sprintf("todos:user:%d", todo.UserID)
	s.cache.Delete(ctx, cacheKey)

	return nil
}
