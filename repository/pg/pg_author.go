package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/syahidfrd/go-boilerplate/domain"
)

type pgAuthorRepository struct {
	db *sql.DB
}

// NewAuthorRepository will create new an authorRepository object representation of domain.AuthorRepository interface
func NewPostgresqlAuthorRepository(db *sql.DB) domain.AuthorRepository {
	return &pgAuthorRepository{
		db: db,
	}
}

func (r *pgAuthorRepository) Create(ctx context.Context, author *domain.Author) (err error) {
	query := "INSERT INTO authors (name, created_at, updated_at) VALUES ($1, $2, $3)"
	_, err = r.db.ExecContext(ctx, query, author.Name, author.CreatedAt, author.UpdatedAt)
	return
}

func (r *pgAuthorRepository) GetByID(ctx context.Context, id int64) (author domain.Author, err error) {
	query := "SELECT id, name, created_at, updated_at FROM authors WHERE id = $1"
	err = r.db.QueryRowContext(ctx, query, id).Scan(&author.ID, &author.Name, &author.CreatedAt, &author.UpdatedAt)
	return
}

func (r *pgAuthorRepository) Fetch(ctx context.Context) (authors []domain.Author, err error) {
	query := "SELECT id, name, created_at, updated_at FROM authors"
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

func (r *pgAuthorRepository) Update(ctx context.Context, author *domain.Author) (err error) {
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

func (r *pgAuthorRepository) Delete(ctx context.Context, id int64) (err error) {
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