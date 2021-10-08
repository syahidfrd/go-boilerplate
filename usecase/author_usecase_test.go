package usecase_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/syahidfrd/go-boilerplate/domain"
	"github.com/syahidfrd/go-boilerplate/domain/mocks"
	"github.com/syahidfrd/go-boilerplate/transport/request"
	"github.com/syahidfrd/go-boilerplate/usecase"
)

func TestCreate(t *testing.T) {
	mockRedisRepository := new(mocks.RedisRepository)
	mockAuthorRepository := new(mocks.AuthorRepository)
	mockCreateAuthorReq := request.CreateAuthorReq{
		Name: "name",
	}

	t.Run("success", func(t *testing.T) {
		mockAuthorRepository.On("Create", mock.Anything, mock.AnythingOfType("*domain.Author")).Return(nil).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		err := authorUsecase.Create(context.TODO(), &mockCreateAuthorReq)

		assert.NoError(t, err)
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockAuthorRepository.On("Create", mock.Anything, mock.AnythingOfType("*domain.Author")).Return(errors.New("Unexpected Error")).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		err := authorUsecase.Create(context.TODO(), &mockCreateAuthorReq)

		assert.NotNil(t, err)
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	mockRedisRepository := new(mocks.RedisRepository)
	mockAuthorRepository := new(mocks.AuthorRepository)
	mockAuthor := domain.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAuthorRepository.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		author, err := authorUsecase.GetByID(context.TODO(), mockAuthor.ID)

		assert.NoError(t, err)
		assert.NotNil(t, author)
		assert.Equal(t, author.ID, mockAuthor.ID)
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})

	t.Run("author-not-exist", func(t *testing.T) {
		mockAuthorRepository.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Author{}, sql.ErrNoRows).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		author, err := authorUsecase.GetByID(context.TODO(), mockAuthor.ID)

		assert.NotNil(t, err)
		assert.Equal(t, author, domain.Author{})
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockAuthorRepository.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Author{}, errors.New("Unexpected Error")).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		author, err := authorUsecase.GetByID(context.TODO(), mockAuthor.ID)

		assert.NotNil(t, err)
		assert.Equal(t, author, domain.Author{})
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})
}

func TestFetch(t *testing.T) {
	mockRedisRepository := new(mocks.RedisRepository)
	mockAuthorRepository := new(mocks.AuthorRepository)
	mockAuthor := domain.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockListAuthor := make([]domain.Author, 0)
	mockListAuthor = append(mockListAuthor, mockAuthor)

	t.Run("success", func(t *testing.T) {
		mockRedisRepository.On("Get", mock.AnythingOfType("string")).Return("", errors.New("Unexpected Error")).Once()
		mockAuthorRepository.On("Fetch", mock.Anything).Return(mockListAuthor, nil).Once()
		mockRedisRepository.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		authors, err := authorUsecase.Fetch(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, authors, len(mockListAuthor))
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})

	t.Run("success-get-from-cache", func(t *testing.T) {
		mockListAuthorByte, _ := json.Marshal(mockListAuthor)
		mockRedisRepository.On("Get", mock.AnythingOfType("string")).Return(string(mockListAuthorByte), nil).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		authors, err := authorUsecase.Fetch(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, authors, len(mockListAuthor))
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)

	})

	t.Run("error-db", func(t *testing.T) {
		mockRedisRepository.On("Get", mock.AnythingOfType("string")).Return("", errors.New("Unexpected Error")).Once()
		mockAuthorRepository.On("Fetch", mock.Anything).Return([]domain.Author{}, errors.New("Unexpected Error")).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		authors, err := authorUsecase.Fetch(context.TODO())

		assert.NotNil(t, err)
		assert.Len(t, authors, 0)
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockRedisRepository := new(mocks.RedisRepository)
	mockAuthorRepository := new(mocks.AuthorRepository)
	mockAuthor := domain.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUpdateAuthorReq := request.UpdateAuthorReq{
		Name: "name 2",
	}

	t.Run("success", func(t *testing.T) {
		mockAuthorRepository.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil).Once()
		mockAuthorRepository.On("Update", mock.Anything, mock.AnythingOfType("*domain.Author")).Return(nil).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		err := authorUsecase.Update(context.TODO(), mockAuthor.ID, &mockUpdateAuthorReq)

		assert.NoError(t, err)
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})

	t.Run("author-not-exist", func(t *testing.T) {
		mockAuthorRepository.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Author{}, sql.ErrNoRows).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		err := authorUsecase.Update(context.TODO(), mockAuthor.ID, &mockUpdateAuthorReq)

		assert.NotNil(t, err)
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockAuthorRepository.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil).Once()
		mockAuthorRepository.On("Update", mock.Anything, mock.AnythingOfType("*domain.Author")).Return(errors.New("Unexpected Error")).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		err := authorUsecase.Update(context.TODO(), mockAuthor.ID, &mockUpdateAuthorReq)

		assert.NotNil(t, err)
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockRedisRepository := new(mocks.RedisRepository)
	mockAuthorRepository := new(mocks.AuthorRepository)
	mockAuthor := domain.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAuthorRepository.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil).Once()
		mockAuthorRepository.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		authorRepository := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		err := authorRepository.Delete(context.TODO(), mockAuthor.ID)

		assert.NoError(t, err)
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})

	t.Run("author-not-exist", func(t *testing.T) {
		mockAuthorRepository.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Author{}, sql.ErrNoRows).Once()

		authorRepository := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		err := authorRepository.Delete(context.TODO(), mockAuthor.ID)

		assert.NotNil(t, err)
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockAuthorRepository.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil).Once()
		mockAuthorRepository.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(errors.New("Unexpected Error")).Once()

		authorRepository := usecase.NewAuthorUsecase(mockAuthorRepository, mockRedisRepository)
		err := authorRepository.Delete(context.TODO(), mockAuthor.ID)

		assert.NotNil(t, err)
		mockRedisRepository.AssertExpectations(t)
		mockAuthorRepository.AssertExpectations(t)
	})
}
