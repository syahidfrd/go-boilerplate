package pgsql_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/syahidfrd/go-boilerplate/entity"
	"github.com/syahidfrd/go-boilerplate/repository/pgsql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"testing"
	"time"
)

func TestUserRepo_Create(t *testing.T) {
	user := &entity.User{
		Email:     "sample@mail.com",
		Password:  "randomPasswordHash",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	query := "INSERT INTO users"
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(user.Email, user.Password, user.CreatedAt, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := pgsql.NewPgsqlUserRepository(db)
	err = userRepo.Create(context.TODO(), user)
	assert.NoError(t, err)
}

func TestUserRepo_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userMock := entity.User{
		ID:        1,
		Email:     "sample@mail.com",
		Password:  "randomPasswordHash",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
		AddRow(userMock.ID, userMock.Email, userMock.Password, userMock.CreatedAt, userMock.UpdatedAt)

	query := "SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userMock.Email).
		WillReturnRows(rows)

	userRepo := pgsql.NewPgsqlUserRepository(db)
	user, err := userRepo.GetByEmail(context.TODO(), userMock.Email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userMock.ID, user.ID)
}
