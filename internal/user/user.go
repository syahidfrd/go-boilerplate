package user

import "time"

// User represents a user account with authentication credentials
type User struct {
	ID        int64
	Email     string `gorm:"uniqueIndex"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user with the given email and hashed password
func NewUser(email, hashedPassword string) *User {
	now := time.Now()
	return &User{
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
