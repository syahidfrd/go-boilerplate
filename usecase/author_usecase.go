package usecase

import (
	"context"

	"github.com/syahidfrd/go-boilerplate/domain"
	"github.com/syahidfrd/go-boilerplate/transport/request"
)

type authorUsecase struct {
	authorRepository domain.AuthorRepository
}

func NewAuthorUsecase(authorRepository domain.AuthorRepository) domain.AuthorUsecase {
	return &authorUsecase{
		authorRepository: authorRepository,
	}
}

func (u *authorUsecase) Create(ctx context.Context, request *request.CreateAuthorReq) (err error) {
	err = u.authorRepository.Create(ctx, &domain.Author{
		Name: request.Name,
	})
	return
}

func (u *authorUsecase) GetByID(ctx context.Context, id uint64) (author domain.Author, err error) {
	author, err = u.authorRepository.GetByID(ctx, id)
	return
}

func (u *authorUsecase) Fetch(ctx context.Context) (authors []domain.Author, err error) {
	authors, err = u.authorRepository.Fetch(ctx)
	return
}

func (u *authorUsecase) Update(ctx context.Context, id uint64, request *request.UpdateAuthorReq) (err error) {
	author, err := u.authorRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	author.Name = request.Name

	err = u.authorRepository.Update(ctx, &author)
	return
}

func (u *authorUsecase) Delete(ctx context.Context, id uint64) (err error) {
	_, err = u.authorRepository.GetByID(ctx, id)
	if err != nil {
		return
	}

	err = u.authorRepository.Delete(ctx, id)
	return
}
