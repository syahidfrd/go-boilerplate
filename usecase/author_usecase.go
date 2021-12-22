package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/syahidfrd/go-boilerplate/entity"
	"github.com/syahidfrd/go-boilerplate/repository/pgsql"
	"github.com/syahidfrd/go-boilerplate/repository/redis"
	"github.com/syahidfrd/go-boilerplate/transport/request"
	"github.com/syahidfrd/go-boilerplate/utils"
)

// AuthorUsecase represent the author's usecase contract
type AuthorUsecase interface {
	Create(ctx context.Context, request *request.CreateAuthorReq) error
	GetByID(ctx context.Context, id int64) (entity.Author, error)
	Fetch(ctx context.Context) ([]entity.Author, error)
	Update(ctx context.Context, id int64, request *request.UpdateAuthorReq) error
	Delete(ctx context.Context, id int64) error
}

type authorUsecase struct {
	authorRepo pgsql.AuthorRepository
	redisRepo  redis.RedisRepository
}

// NewAuthorUsecase will create new an authorUsecase object representation of AuthorUsecase interface
func NewAuthorUsecase(authorRepo pgsql.AuthorRepository, redisRepo redis.RedisRepository) AuthorUsecase {
	return &authorUsecase{
		authorRepo: authorRepo,
		redisRepo:  redisRepo,
	}
}

func (u *authorUsecase) Create(ctx context.Context, request *request.CreateAuthorReq) (err error) {
	err = u.authorRepo.Create(ctx, &entity.Author{
		Name:      request.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	return
}

func (u *authorUsecase) GetByID(ctx context.Context, id int64) (author entity.Author, err error) {
	author, err = u.authorRepo.GetByID(ctx, id)
	if err != nil && err == sql.ErrNoRows {
		err = utils.NewNotFoundError("author not found")
		return
	}
	return
}

func (u *authorUsecase) Fetch(ctx context.Context) (authors []entity.Author, err error) {
	authorsCached, _ := u.redisRepo.Get("authors")
	if err = json.Unmarshal([]byte(authorsCached), &authors); err == nil {
		return
	}

	authors, err = u.authorRepo.Fetch(ctx)
	if err != nil {
		return
	}

	authorsString, _ := json.Marshal(&authors)
	u.redisRepo.Set("authors", authorsString, 60*time.Second)
	return
}

func (u *authorUsecase) Update(ctx context.Context, id int64, request *request.UpdateAuthorReq) (err error) {
	author, err := u.authorRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewNotFoundError("author not found")
			return
		}
		return
	}

	author.Name = request.Name
	author.UpdatedAt = time.Now()

	err = u.authorRepo.Update(ctx, &author)
	return
}

func (u *authorUsecase) Delete(ctx context.Context, id int64) (err error) {
	_, err = u.authorRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewNotFoundError("author not found")
			return
		}
		return
	}

	err = u.authorRepo.Delete(ctx, id)
	return
}
