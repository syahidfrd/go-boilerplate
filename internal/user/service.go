package user

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	// ErrUserNotFound is returned when a requested user cannot be found
	ErrUserNotFound = errors.New("user not found")
)

// Store defines the interface for user data persistence operations
type Store interface {
	Save(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

// Service provides user business logic operations
type Service struct {
	store Store
}

// NewService creates a new user service with the provided store
func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}

// Create creates a new user with the given email and hashed password
func (s *Service) Create(ctx context.Context, email, hashedPassword string) (*User, error) {
	user := NewUser(email, hashedPassword)

	err := s.store.Save(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by their email address
func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.store.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return user, nil
}
