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

func TestTodoUC_Create(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockTodoRepo := new(mocks.TodoRepository)
	createTodoReq := request.CreateTodoReq{
		Name: "name",
	}

	t.Run("success", func(t *testing.T) {
		mockTodoRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Todo")).Return(nil).Once()

		todoUsecase := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		err := todoUsecase.Create(context.TODO(), &createTodoReq)

		assert.NoError(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockTodoRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Todo")).Return(errors.New("Unexpected Error")).Once()

		todoUsecase := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		err := todoUsecase.Create(context.TODO(), &createTodoReq)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)
	})
}

func TestTodoUC_GetByID(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockTodoRepo := new(mocks.TodoRepository)
	mockTodo := entity.Todo{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockTodoRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockTodo, nil).Once()

		todoUsecase := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		todo, err := todoUsecase.GetByID(context.TODO(), mockTodo.ID)

		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.Equal(t, todo.ID, mockTodo.ID)
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)
	})

	t.Run("todo-not-exist", func(t *testing.T) {
		mockTodoRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entity.Todo{}, sql.ErrNoRows).Once()

		todoUsecase := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		todo, err := todoUsecase.GetByID(context.TODO(), mockTodo.ID)

		assert.NotNil(t, err)
		assert.Equal(t, todo, entity.Todo{})
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockTodoRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entity.Todo{}, errors.New("Unexpected Error")).Once()

		todoUsecase := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		todo, err := todoUsecase.GetByID(context.TODO(), mockTodo.ID)

		assert.NotNil(t, err)
		assert.Equal(t, todo, entity.Todo{})
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)
	})
}

func TestTooUC_Fetch(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockTodoRepo := new(mocks.TodoRepository)
	mockTodo := entity.Todo{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockListTodo := make([]entity.Todo, 0)
	mockListTodo = append(mockListTodo, mockTodo)

	t.Run("success", func(t *testing.T) {
		mockRedisRepo.On("Get", mock.AnythingOfType("string")).Return("", errors.New("Unexpected Error")).Once()
		mockTodoRepo.On("Fetch", mock.Anything).Return(mockListTodo, nil).Once()
		mockRedisRepo.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil).Once()

		todoUsecase := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		todos, err := todoUsecase.Fetch(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, todos, len(mockListTodo))
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)
	})

	t.Run("success-get-from-cache", func(t *testing.T) {
		mockListTodoByte, _ := json.Marshal(mockListTodo)
		mockRedisRepo.On("Get", mock.AnythingOfType("string")).Return(string(mockListTodoByte), nil).Once()

		todoUsecase := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		todos, err := todoUsecase.Fetch(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, todos, len(mockListTodo))
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)

	})

	t.Run("error-db", func(t *testing.T) {
		mockRedisRepo.On("Get", mock.AnythingOfType("string")).Return("", errors.New("Unexpected Error")).Once()
		mockTodoRepo.On("Fetch", mock.Anything).Return([]entity.Todo{}, errors.New("Unexpected Error")).Once()

		todoUsecase := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		todos, err := todoUsecase.Fetch(context.TODO())

		assert.NotNil(t, err)
		assert.Len(t, todos, 0)
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)
	})
}

func TestTodoUC_Update(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockTodoRepo := new(mocks.TodoRepository)
	mockTodo := entity.Todo{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	updateTodoReq := request.UpdateTodoReq{
		Name: "name 2",
	}

	t.Run("success", func(t *testing.T) {
		mockTodoRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockTodo, nil).Once()
		mockTodoRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Todo")).Return(nil).Once()

		todoUsecase := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		err := todoUsecase.Update(context.TODO(), mockTodo.ID, &updateTodoReq)

		assert.NoError(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)
	})

	t.Run("todo-not-exist", func(t *testing.T) {
		mockTodoRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entity.Todo{}, sql.ErrNoRows).Once()

		todoUsecase := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		err := todoUsecase.Update(context.TODO(), mockTodo.ID, &updateTodoReq)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockTodoRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockTodo, nil).Once()
		mockTodoRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Todo")).Return(errors.New("Unexpected Error")).Once()

		todoUsecase := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		err := todoUsecase.Update(context.TODO(), mockTodo.ID, &updateTodoReq)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)
	})
}

func TestTodoUC_Delete(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockTodoRepo := new(mocks.TodoRepository)
	mockTodo := entity.Todo{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockTodoRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockTodo, nil).Once()
		mockTodoRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		todoRepository := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		err := todoRepository.Delete(context.TODO(), mockTodo.ID)

		assert.NoError(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)
	})

	t.Run("todo-not-exist", func(t *testing.T) {
		mockTodoRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entity.Todo{}, sql.ErrNoRows).Once()

		todoRepository := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		err := todoRepository.Delete(context.TODO(), mockTodo.ID)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockTodoRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockTodo, nil).Once()
		mockTodoRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(errors.New("Unexpected Error")).Once()

		todoRepository := usecase.NewTodoUsecase(mockTodoRepo, mockRedisRepo, 60*time.Second)
		err := todoRepository.Delete(context.TODO(), mockTodo.ID)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockTodoRepo.AssertExpectations(t)
	})
}
