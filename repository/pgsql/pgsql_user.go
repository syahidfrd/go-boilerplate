package pgsql

import (
	"context"

	"github.com/syahidfrd/go-boilerplate/domain"
)

type pgsqlUserRepository struct {
}

func NewPgsqlUserRepository() *pgsqlUserRepository {
	return &pgsqlUserRepository{}
}

func (r *pgsqlUserRepository) Create(ctx context.Context, tx domain.Transaction, user *domain.User) (err error) {
	query := "INSERT INTO users (email, password, created_at, updated_at) VALUES ($1, $2, $3, $4)"
	_, err = tx.ExecContext(ctx, query, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	return
}

func (r *pgsqlUserRepository) GetByEmail(ctx context.Context, tx domain.Transaction, email string) (user domain.User, err error) {
	query := "SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1"
	err = tx.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	return
}
