package pgsql

import (
	"context"
	"fmt"

	"github.com/syahidfrd/go-boilerplate/domain"
)

type pgsqlTodoRepository struct {
}

// NewPgsqlTodoRepository will create new an todoRepository object representation of TodoRepository interface
func NewPgsqlTodoRepository() *pgsqlTodoRepository {
	return &pgsqlTodoRepository{}
}

func (r *pgsqlTodoRepository) Create(ctx context.Context, tx domain.Transaction, todo *domain.Todo) (err error) {
	query := "INSERT INTO todos (name, created_at, updated_at) VALUES ($1, $2, $3)"
	_, err = tx.ExecContext(ctx, query, todo.Name, todo.CreatedAt, todo.UpdatedAt)
	return
}

func (r *pgsqlTodoRepository) GetByID(ctx context.Context, tx domain.Transaction, id int64) (todo domain.Todo, err error) {
	query := "SELECT id, name, created_at, updated_at FROM todos WHERE id = $1"
	err = tx.QueryRowContext(ctx, query, id).Scan(&todo.ID, &todo.Name, &todo.CreatedAt, &todo.UpdatedAt)
	return
}

func (r *pgsqlTodoRepository) Fetch(ctx context.Context, tx domain.Transaction) (todos []domain.Todo, err error) {
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

	return todos, nil
}

func (r *pgsqlTodoRepository) Update(ctx context.Context, tx domain.Transaction, todo *domain.Todo) (err error) {
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

	return
}

func (r *pgsqlTodoRepository) Delete(ctx context.Context, tx domain.Transaction, id int64) (err error) {
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

	return
}
