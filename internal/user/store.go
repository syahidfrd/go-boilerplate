package user

import (
	"context"

	"github.com/syahidfrd/go-boilerplate/internal/pkg/db"
	"gorm.io/gorm"
)

// store implements user data persistence using GORM
type store struct {
	dbConn *gorm.DB
}

// NewStore creates a new user store with the provided database connection
func NewStore(dbConn *gorm.DB) *store {
	return &store{dbConn: dbConn}
}

// Save persists a user to the database (create or update)
func (s *store) Save(ctx context.Context, user *User, options ...db.Option) error {
	dbConn := s.dbConn

	opts := &db.Options{}
	for _, opt := range options {
		opt(opts)
	}

	if opts.Tx != nil {
		dbConn = opts.Tx
	}

	return dbConn.WithContext(ctx).Save(user).Error
}

// FindByEmail retrieves a user by their email address from the database
func (s *store) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := s.dbConn.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// SavePreference persists a user preference to the database (create or update)
func (s *store) SavePreference(ctx context.Context, preference *Preference, options ...db.Option) error {
	dbConn := s.dbConn

	opts := &db.Options{}
	for _, opt := range options {
		opt(opts)
	}

	if opts.Tx != nil {
		dbConn = opts.Tx
	}

	return dbConn.WithContext(ctx).Save(preference).Error
}
