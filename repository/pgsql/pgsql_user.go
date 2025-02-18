package pgsql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/syahidfrd/go-boilerplate/domain"
)

type pgsqlUserRepository struct {
	db *sql.DB
}

func NewPgsqlUserRepository(db *sql.DB) *pgsqlUserRepository {
	return &pgsqlUserRepository{
		db: db,
	}
}

func (r *pgsqlUserRepository) Create(ctx context.Context, user *domain.User) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()

	query := "INSERT INTO users (email, password, created_at, updated_at) VALUES ($1, $2, $3, $4)"
	_, err = tx.ExecContext(ctx, query, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("failed to commit transaction: %v", err)
		return
	}

	return
}

func (r *pgsqlUserRepository) GetByEmail(ctx context.Context, email string) (user domain.User, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()

	query := "SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1"
	err = tx.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("failed to commit transaction: %v", err)
		return
	}

	return
}
