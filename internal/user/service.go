package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/syahidfrd/go-boilerplate/internal/pkg/db"
	"gorm.io/gorm"
)

var (
	// ErrUserNotFound is returned when a requested user cannot be found
	ErrUserNotFound = errors.New("user not found")
)

// Service provides user business logic operations
type Service struct {
	store *store
}

// NewService creates a new user service with the provided store
func NewService(store *store) *Service {
	return &Service{
		store: store,
	}
}

// Create creates a new user with the given email and hashed password, along with default preferences
func (s *Service) Create(ctx context.Context, email, hashedPassword string) (*User, error) {
	user := NewUser(email, hashedPassword)

	// Start database transaction
	tx := s.store.dbConn.Begin()

	// Create user within transaction
	err := s.store.Save(ctx, user, db.WithTx(tx))
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Create default preferences for the user
	preference := NewPreference(user.ID)
	err = s.store.SavePreference(ctx, preference, db.WithTx(tx))
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to save user preferences: %w", err)
	}

	// Commit transaction if all operations succeed
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit db transaction: %w", err)
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
