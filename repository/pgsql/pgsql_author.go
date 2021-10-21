package pgsql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/syahidfrd/go-boilerplate/entity"
)

// AuthorRepository represent the author's repository contract
type AuthorRepository interface {
	Create(ctx context.Context, author *entity.Author) error
	GetByID(ctx context.Context, id int64) (entity.Author, error)
	Fetch(ctx context.Context) ([]entity.Author, error)
	Update(ctx context.Context, author *entity.Author) error
	Delete(ctx context.Context, id int64) error
}

type pgsqlAuthorRepository struct {
	db *sql.DB
}

// NewAuthorRepository will create new an authorRepository object representation of AuthorRepository interface
func NewPgsqlAuthorRepository(db *sql.DB) AuthorRepository {
	return &pgsqlAuthorRepository{
		db: db,
	}
}

func (r *pgsqlAuthorRepository) Create(ctx context.Context, author *entity.Author) (err error) {
	query := "INSERT INTO authors (name, created_at, updated_at) VALUES ($1, $2, $3)"
	_, err = r.db.ExecContext(ctx, query, author.Name, author.CreatedAt, author.UpdatedAt)
	return
}

func (r *pgsqlAuthorRepository) GetByID(ctx context.Context, id int64) (author entity.Author, err error) {
	query := "SELECT id, name, created_at, updated_at FROM authors WHERE id = $1"
	err = r.db.QueryRowContext(ctx, query, id).Scan(&author.ID, &author.Name, &author.CreatedAt, &author.UpdatedAt)
	return
}

func (r *pgsqlAuthorRepository) Fetch(ctx context.Context) (authors []entity.Author, err error) {
	query := "SELECT id, name, created_at, updated_at FROM authors"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return authors, err
	}

	defer rows.Close()

	for rows.Next() {
		var author entity.Author
		err := rows.Scan(&author.ID, &author.Name, &author.CreatedAt, &author.UpdatedAt)
		if err != nil {
			return authors, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func (r *pgsqlAuthorRepository) Update(ctx context.Context, author *entity.Author) (err error) {
	query := "UPDATE authors SET name = $1, updated_at = $2 WHERE id = $3"
	res, err := r.db.ExecContext(ctx, query, author.Name, author.UpdatedAt, author.ID)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("weird behavior, total affected: %d", affect)
	}

	return
}

func (r *pgsqlAuthorRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM authors WHERE id = $1"
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("weird behavior, total affected: %d", affect)
	}

	return
}
