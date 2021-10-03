package repository

import (
	"context"
	"database/sql"

	"github.com/syahidfrd/go-boilerplate/domain"
)

type authorRepository struct {
	db *sql.DB
}

func NewAuthorRepository(db *sql.DB) domain.AuthorRepository {
	return &authorRepository{
		db: db,
	}
}

func (r *authorRepository) Create(ctx context.Context, author *domain.Author) (err error) {
	query := "INSERT INTO authors (name, created_at, updated_at) VALUES ($1, now(), now())"
	_, err = r.db.ExecContext(ctx, query, author.Name)
	return
}

func (r *authorRepository) GetByID(ctx context.Context, id uint64) (author domain.Author, err error) {
	query := "SELECT * FROM authors WHERE id = $1"
	err = r.db.QueryRowContext(ctx, query, id).Scan(&author.ID, &author.Name, &author.CreatedAt, &author.UpdatedAt)
	return
}

func (r *authorRepository) Fetch(ctx context.Context) (authors []domain.Author, err error) {
	query := "SELECT * FROM authors"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return authors, err
	}

	defer rows.Close()

	for rows.Next() {
		var author domain.Author
		err := rows.Scan(&author.ID, &author.Name, &author.CreatedAt, &author.UpdatedAt)
		if err != nil {
			return authors, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func (r *authorRepository) Update(ctx context.Context, author *domain.Author) (err error) {
	query := "UPDATE authors SET name = $1, updated_at = now() WHERE id = $2"
	_, err = r.db.ExecContext(ctx, query, author.Name, author.ID)
	return
}

func (r *authorRepository) Delete(ctx context.Context, id uint64) (err error) {
	query := "DELETE FROM authors WHERE id = $1"
	_, err = r.db.ExecContext(ctx, query, id)
	return
}
