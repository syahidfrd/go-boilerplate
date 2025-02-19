package domain

import (
	"context"
	"time"

	"github.com/syahidfrd/go-boilerplate/transport/request"
)

// Todo ...
type Todo struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TodoRepository represent the todos repository contract
type TodoRepository interface {
	Create(ctx context.Context, tx Transaction, todo *Todo) error
	GetByID(ctx context.Context, tx Transaction, id int64) (Todo, error)
	Fetch(ctx context.Context, tx Transaction) ([]Todo, error)
	Update(ctx context.Context, tx Transaction, todo *Todo) error
	Delete(ctx context.Context, tx Transaction, id int64) error
}

// TodoUsecase represent the todos usecase contract
type TodoUsecase interface {
	Create(ctx context.Context, request *request.CreateTodoReq) error
	GetByID(ctx context.Context, id int64) (Todo, error)
	Fetch(ctx context.Context) ([]Todo, error)
	Update(ctx context.Context, id int64, request *request.UpdateTodoReq) error
	Delete(ctx context.Context, id int64) error
}
