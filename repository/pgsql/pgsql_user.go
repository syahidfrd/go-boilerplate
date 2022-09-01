package pgsql

import (
	"context"
	"database/sql"
	"github.com/syahidfrd/go-boilerplate/entity"
)

// UserRepository represent the todos repository contract
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type pgsqlUserRepository struct {
	db *sql.DB
}

func NewPgsqlUserRepository(db *sql.DB) UserRepository {
	return &pgsqlUserRepository{
		db: db,
	}
}

func (r *pgsqlUserRepository) Create(ctx context.Context, user *entity.User) (err error) {
	query := "INSERT INTO users (email, password, created_at, updated_at) VALUES ($1, $2, $3, $4)"
	_, err = r.db.ExecContext(ctx, query, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	return
}

func (r *pgsqlUserRepository) GetByEmail(ctx context.Context, email string) (user entity.User, err error) {
	query := "SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1"
	err = r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	return
}
