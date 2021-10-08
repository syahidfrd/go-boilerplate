package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/syahidfrd/go-boilerplate/domain"
	"github.com/syahidfrd/go-boilerplate/repository/redis"
	"github.com/syahidfrd/go-boilerplate/transport/request"
)

type authorUsecase struct {
	authorRepository domain.AuthorRepository
	redisRepository  redis.RedisRepository
}

// NewAuthorUsecase will create new an authorUsecase object representation of domain.AuthorUsecase interface
func NewAuthorUsecase(authorRepository domain.AuthorRepository, redisRepository redis.RedisRepository) domain.AuthorUsecase {
	return &authorUsecase{
		authorRepository: authorRepository,
		redisRepository:  redisRepository,
	}
}

func (u *authorUsecase) Create(ctx context.Context, request *request.CreateAuthorReq) (err error) {
	err = u.authorRepository.Create(ctx, &domain.Author{
		Name:      request.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	return
}

func (u *authorUsecase) GetByID(ctx context.Context, id int64) (author domain.Author, err error) {
	author, err = u.authorRepository.GetByID(ctx, id)
	return
}

func (u *authorUsecase) Fetch(ctx context.Context) (authors []domain.Author, err error) {
	authorsCached, _ := u.redisRepository.Get("authors")
	if err = json.Unmarshal([]byte(authorsCached), &authors); err == nil {
		return
	}

	authors, err = u.authorRepository.Fetch(ctx)
	if err != nil {
		return
	}

	authorsString, _ := json.Marshal(&authors)
	u.redisRepository.Set("authors", authorsString, 60*time.Second)
	return
}

func (u *authorUsecase) Update(ctx context.Context, id int64, request *request.UpdateAuthorReq) (err error) {
	author, err := u.authorRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	author.Name = request.Name
	author.UpdatedAt = time.Now()

	err = u.authorRepository.Update(ctx, &author)
	return
}

func (u *authorUsecase) Delete(ctx context.Context, id int64) (err error) {
	_, err = u.authorRepository.GetByID(ctx, id)
	if err != nil {
		return
	}

	err = u.authorRepository.Delete(ctx, id)
	return
}
