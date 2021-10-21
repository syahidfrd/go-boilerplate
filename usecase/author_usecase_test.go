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
	"github.com/syahidfrd/go-boilerplate/entity"
	"github.com/syahidfrd/go-boilerplate/mocks"
	"github.com/syahidfrd/go-boilerplate/transport/request"
	"github.com/syahidfrd/go-boilerplate/usecase"
)

func TestAuthorUC_Create(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockAuthorRepo := new(mocks.AuthorRepository)
	createAuthorReq := request.CreateAuthorReq{
		Name: "name",
	}

	t.Run("success", func(t *testing.T) {
		mockAuthorRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Author")).Return(nil).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		err := authorUsecase.Create(context.TODO(), &createAuthorReq)

		assert.NoError(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockAuthorRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Author")).Return(errors.New("Unexpected Error")).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		err := authorUsecase.Create(context.TODO(), &createAuthorReq)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestAuthorUC_GetByID(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockAuthorRepo := new(mocks.AuthorRepository)
	mockAuthor := entity.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		author, err := authorUsecase.GetByID(context.TODO(), mockAuthor.ID)

		assert.NoError(t, err)
		assert.NotNil(t, author)
		assert.Equal(t, author.ID, mockAuthor.ID)
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("author-not-exist", func(t *testing.T) {
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entity.Author{}, sql.ErrNoRows).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		author, err := authorUsecase.GetByID(context.TODO(), mockAuthor.ID)

		assert.NotNil(t, err)
		assert.Equal(t, author, entity.Author{})
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entity.Author{}, errors.New("Unexpected Error")).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		author, err := authorUsecase.GetByID(context.TODO(), mockAuthor.ID)

		assert.NotNil(t, err)
		assert.Equal(t, author, entity.Author{})
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestAuthorUC_Fetch(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockAuthorRepo := new(mocks.AuthorRepository)
	mockAuthor := entity.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockListAuthor := make([]entity.Author, 0)
	mockListAuthor = append(mockListAuthor, mockAuthor)

	t.Run("success", func(t *testing.T) {
		mockRedisRepo.On("Get", mock.AnythingOfType("string")).Return("", errors.New("Unexpected Error")).Once()
		mockAuthorRepo.On("Fetch", mock.Anything).Return(mockListAuthor, nil).Once()
		mockRedisRepo.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		authors, err := authorUsecase.Fetch(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, authors, len(mockListAuthor))
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("success-get-from-cache", func(t *testing.T) {
		mockListAuthorByte, _ := json.Marshal(mockListAuthor)
		mockRedisRepo.On("Get", mock.AnythingOfType("string")).Return(string(mockListAuthorByte), nil).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		authors, err := authorUsecase.Fetch(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, authors, len(mockListAuthor))
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)

	})

	t.Run("error-db", func(t *testing.T) {
		mockRedisRepo.On("Get", mock.AnythingOfType("string")).Return("", errors.New("Unexpected Error")).Once()
		mockAuthorRepo.On("Fetch", mock.Anything).Return([]entity.Author{}, errors.New("Unexpected Error")).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		authors, err := authorUsecase.Fetch(context.TODO())

		assert.NotNil(t, err)
		assert.Len(t, authors, 0)
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestAuthorUC_Update(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockAuthorRepo := new(mocks.AuthorRepository)
	mockAuthor := entity.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	updateAuthorReq := request.UpdateAuthorReq{
		Name: "name 2",
	}

	t.Run("success", func(t *testing.T) {
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil).Once()
		mockAuthorRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Author")).Return(nil).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		err := authorUsecase.Update(context.TODO(), mockAuthor.ID, &updateAuthorReq)

		assert.NoError(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("author-not-exist", func(t *testing.T) {
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entity.Author{}, sql.ErrNoRows).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		err := authorUsecase.Update(context.TODO(), mockAuthor.ID, &updateAuthorReq)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil).Once()
		mockAuthorRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Author")).Return(errors.New("Unexpected Error")).Once()

		authorUsecase := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		err := authorUsecase.Update(context.TODO(), mockAuthor.ID, &updateAuthorReq)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestAuthorUC_Delete(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockAuthorRepo := new(mocks.AuthorRepository)
	mockAuthor := entity.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil).Once()
		mockAuthorRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		authorRepository := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		err := authorRepository.Delete(context.TODO(), mockAuthor.ID)

		assert.NoError(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("author-not-exist", func(t *testing.T) {
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entity.Author{}, sql.ErrNoRows).Once()

		authorRepository := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		err := authorRepository.Delete(context.TODO(), mockAuthor.ID)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil).Once()
		mockAuthorRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(errors.New("Unexpected Error")).Once()

		authorRepository := usecase.NewAuthorUsecase(mockAuthorRepo, mockRedisRepo)
		err := authorRepository.Delete(context.TODO(), mockAuthor.ID)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}
