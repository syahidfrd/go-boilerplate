package pgsql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/syahidfrd/go-boilerplate/domain"
)

type pgsqlTodoRepository struct {
	db *sql.DB
}

// NewPgsqlTodoRepository will create a new todoRepository object representation of TodoRepository interface
func NewPgsqlTodoRepository(db *sql.DB) *pgsqlTodoRepository {
	return &pgsqlTodoRepository{
		db: db,
	}
}

func (r *pgsqlTodoRepository) Create(ctx context.Context, todo *domain.Todo) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback()

	query := "INSERT INTO todos (name, created_at, updated_at) VALUES ($1, $2, $3)"
	_, err = tx.ExecContext(ctx, query, todo.Name, todo.CreatedAt, todo.UpdatedAt)

	if err != nil {
		err = fmt.Errorf("failed to insert todo: %v", err)
		return
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("failed to commit transaction: %v", err)
		return
	}

	return
}

func (r *pgsqlTodoRepository) GetByID(ctx context.Context, id int64) (todo domain.Todo, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback()

	query := "SELECT id, name, created_at, updated_at FROM todos WHERE id = $1"
	err = tx.QueryRowContext(ctx, query, id).Scan(&todo.ID, &todo.Name, &todo.CreatedAt, &todo.UpdatedAt)

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("failed to commit transaction: %v", err)
		return
	}

	return
}

func (r *pgsqlTodoRepository) Fetch(ctx context.Context) (todos []domain.Todo, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to begin transaction: %v", err)
		return
	}

	query := "SELECT id, name, created_at, updated_at FROM todos"
	rows, err := tx.QueryContext(ctx, query)

	if err != nil {
		return todos, err
	}

	defer rows.Close()

	for rows.Next() {
		var todo domain.Todo
		err := rows.Scan(&todo.ID, &todo.Name, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return todos, err
		}

		todos = append(todos, todo)
	}

	if err = tx.Commit(); err != nil {
		return todos, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return todos, nil
}

func (r *pgsqlTodoRepository) Update(ctx context.Context, todo *domain.Todo) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback()

	query := "UPDATE todos SET name = $1, updated_at = $2 WHERE id = $3"
	res, err := tx.ExecContext(ctx, query, todo.Name, todo.UpdatedAt, todo.ID)
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

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("failed to commit transaction: %v", err)
		return
	}

	return
}

func (r *pgsqlTodoRepository) Delete(ctx context.Context, id int64) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback()

	query := "DELETE FROM todos WHERE id = $1"
	res, err := tx.ExecContext(ctx, query, id)
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

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("failed to commit transaction: %v", err)
		return
	}

	return
}
