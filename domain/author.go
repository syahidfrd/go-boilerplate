package domain

import (
	"context"
	"time"

	"github.com/syahidfrd/go-boilerplate/transport/request"
)

type Author struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthorUsecase interface {
	Create(ctx context.Context, request *request.CreateAuthorReq) error
	GetByID(ctx context.Context, id int64) (Author, error)
	Fetch(ctx context.Context) ([]Author, error)
	Update(ctx context.Context, id int64, request *request.UpdateAuthorReq) error
	Delete(ctx context.Context, id int64) error
}

type AuthorRepository interface {
	Create(ctx context.Context, author *Author) error
	GetByID(ctx context.Context, id int64) (Author, error)
	Fetch(ctx context.Context) ([]Author, error)
	Update(ctx context.Context, author *Author) error
	Delete(ctx context.Context, id int64) error
}
