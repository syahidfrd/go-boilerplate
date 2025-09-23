package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	email := "test@example.com"
	hashedPassword := "hashed_password_123"

	user := NewUser(email, hashedPassword)

	assert.Equal(t, email, user.Email)
	assert.Equal(t, hashedPassword, user.Password)
	assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
	assert.Equal(t, user.CreatedAt, user.UpdatedAt)
}

func TestNewUser_EmptyEmail(t *testing.T) {
	email := ""
	hashedPassword := "hashed_password_123"

	user := NewUser(email, hashedPassword)

	assert.Equal(t, email, user.Email)
	assert.Equal(t, hashedPassword, user.Password)
}

func TestNewUser_EmptyPassword(t *testing.T) {
	email := "test@example.com"
	hashedPassword := ""

	user := NewUser(email, hashedPassword)

	assert.Equal(t, email, user.Email)
	assert.Equal(t, hashedPassword, user.Password)
}

func TestNewUser_DifferentInstances(t *testing.T) {
	email1 := "user1@example.com"
	email2 := "user2@example.com"
	password := "same_password"

	user1 := NewUser(email1, password)
	user2 := NewUser(email2, password)

	assert.NotEqual(t, user1, user2)
	assert.Equal(t, email1, user1.Email)
	assert.Equal(t, email2, user2.Email)
	assert.Equal(t, password, user1.Password)
	assert.Equal(t, password, user2.Password)
}

func TestNewUser_TimestampConsistency(t *testing.T) {
	before := time.Now()
	user := NewUser("test@example.com", "password")
	after := time.Now()

	// CreatedAt and UpdatedAt should be the same
	assert.Equal(t, user.CreatedAt, user.UpdatedAt)

	// Timestamps should be within reasonable range
	assert.True(t, user.CreatedAt.After(before) || user.CreatedAt.Equal(before))
	assert.True(t, user.CreatedAt.Before(after) || user.CreatedAt.Equal(after))
}

func TestUser_Fields(t *testing.T) {
	user := NewUser("test@example.com", "hashed_password")

	// Test that all expected fields exist and are properly typed
	assert.IsType(t, int64(0), user.ID)
	assert.IsType(t, "", user.Email)
	assert.IsType(t, "", user.Password)
	assert.IsType(t, time.Time{}, user.CreatedAt)
	assert.IsType(t, time.Time{}, user.UpdatedAt)

	// ID should be zero value (will be set by database)
	assert.Equal(t, int64(0), user.ID)
}

func TestUser_SpecialCharacters(t *testing.T) {
	// Test with special characters in email and password
	email := "test+user@example-domain.com"
	password := "password!@#$%^&*()_+"

	user := NewUser(email, password)

	assert.Equal(t, email, user.Email)
	assert.Equal(t, password, user.Password)
}

func TestUser_LongStrings(t *testing.T) {
	// Test with long email and password
	longEmail := "very.long.email.address.for.testing.purposes@very-long-domain-name-for-testing.com"
	longPassword := "very_long_hashed_password_string_that_might_be_generated_by_bcrypt_or_similar_hashing_function_123456789"

	user := NewUser(longEmail, longPassword)

	assert.Equal(t, longEmail, user.Email)
	assert.Equal(t, longPassword, user.Password)
}
