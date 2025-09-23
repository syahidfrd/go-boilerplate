package user

import (
	"context"

	"gorm.io/gorm"
)

// store implements user data persistence using GORM
type store struct {
	db *gorm.DB
}

// NewStore creates a new user store with the provided database connection
func NewStore(db *gorm.DB) *store {
	return &store{db: db}
}

// Save persists a user to the database (create or update)
func (s *store) Save(ctx context.Context, user *User) error {
	return s.db.WithContext(ctx).Save(user).Error
}

// FindByEmail retrieves a user by their email address from the database
func (s *store) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := s.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
