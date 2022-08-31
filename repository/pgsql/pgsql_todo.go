package pgsql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/syahidfrd/go-boilerplate/entity"
)

// TodoRepository represent the todos repository contract
type TodoRepository interface {
	Create(ctx context.Context, todo *entity.Todo) error
	GetByID(ctx context.Context, id int64) (entity.Todo, error)
	Fetch(ctx context.Context) ([]entity.Todo, error)
	Update(ctx context.Context, todo *entity.Todo) error
	Delete(ctx context.Context, id int64) error
}

type pgsqlTodoRepository struct {
	db *sql.DB
}

// NewPgsqlTodoRepository will create new an todoRepository object representation of TodoRepository interface
func NewPgsqlTodoRepository(db *sql.DB) TodoRepository {
	return &pgsqlTodoRepository{
		db: db,
	}
}

func (r *pgsqlTodoRepository) Create(ctx context.Context, todo *entity.Todo) (err error) {
	query := "INSERT INTO todos (name, created_at, updated_at) VALUES ($1, $2, $3)"
	_, err = r.db.ExecContext(ctx, query, todo.Name, todo.CreatedAt, todo.UpdatedAt)
	return
}

func (r *pgsqlTodoRepository) GetByID(ctx context.Context, id int64) (todo entity.Todo, err error) {
	query := "SELECT id, name, created_at, updated_at FROM todos WHERE id = $1"
	err = r.db.QueryRowContext(ctx, query, id).Scan(&todo.ID, &todo.Name, &todo.CreatedAt, &todo.UpdatedAt)
	return
}

func (r *pgsqlTodoRepository) Fetch(ctx context.Context) (todos []entity.Todo, err error) {
	query := "SELECT id, name, created_at, updated_at FROM todos"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return todos, err
	}

	defer rows.Close()

	for rows.Next() {
		var todo entity.Todo
		err := rows.Scan(&todo.ID, &todo.Name, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return todos, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *pgsqlTodoRepository) Update(ctx context.Context, todo *entity.Todo) (err error) {
	query := "UPDATE todos SET name = $1, updated_at = $2 WHERE id = $3"
	res, err := r.db.ExecContext(ctx, query, todo.Name, todo.UpdatedAt, todo.ID)
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

func (r *pgsqlTodoRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM todos WHERE id = $1"
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
