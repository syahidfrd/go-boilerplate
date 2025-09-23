package todo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) Save(ctx context.Context, todo *Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockStore) GetByID(ctx context.Context, id int64) (*Todo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Todo), args.Error(1)
}

func (m *MockStore) GetByUserID(ctx context.Context, userID int64) ([]Todo, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Todo), args.Error(1)
}

func (m *MockStore) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// NoOpCache for tests that don't need cache behavior
type NoOpCache struct{}

func (n *NoOpCache) Get(ctx context.Context, key string) (string, error) {
	return "", errors.New("cache miss")
}

func (n *NoOpCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return nil
}

func (n *NoOpCache) Delete(ctx context.Context, key string) error {
	return nil
}

func TestService_Create_Success(t *testing.T) {
	mockStore := new(MockStore)
	mockCache := new(MockCache)
	service := NewService(mockStore, mockCache)

	userID := int64(123)
	req := &CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
	}

	mockStore.On("Save", mock.Anything, mock.AnythingOfType("*todo.Todo")).Return(nil)
	mockCache.On("Delete", mock.Anything, "todos:user:123").Return(nil)

	todo, err := service.Create(context.Background(), userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, userID, todo.UserID)
	assert.Equal(t, req.Title, todo.Title)
	assert.Equal(t, req.Description, todo.Description)
	assert.False(t, todo.Completed)
	mockStore.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestService_Create_SaveError(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore, &NoOpCache{})

	userID := int64(123)
	req := &CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
	}

	expectedError := errors.New("database error")
	mockStore.On("Save", mock.Anything, mock.AnythingOfType("*todo.Todo")).Return(expectedError)

	todo, err := service.Create(context.Background(), userID, req)

	assert.Error(t, err)
	assert.Nil(t, todo)
	assert.Contains(t, err.Error(), "failed to create todo")
	mockStore.AssertExpectations(t)
}

func TestService_GetByID_Success(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore, &NoOpCache{})

	todoID := int64(1)
	expectedTodo := &Todo{
		ID:     todoID,
		UserID: 123,
		Title:  "Test Todo",
	}

	mockStore.On("GetByID", mock.Anything, todoID).Return(expectedTodo, nil)

	todo, err := service.GetByID(context.Background(), todoID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTodo, todo)
	mockStore.AssertExpectations(t)
}

func TestService_GetByID_NotFound(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore, &NoOpCache{})

	todoID := int64(1)
	mockStore.On("GetByID", mock.Anything, todoID).Return(nil, ErrTodoNotFound)

	todo, err := service.GetByID(context.Background(), todoID)

	assert.Error(t, err)
	assert.Nil(t, todo)
	assert.Contains(t, err.Error(), "failed to get todo by id")
	mockStore.AssertExpectations(t)
}

func TestService_GetByUserID_Success(t *testing.T) {
	mockStore := new(MockStore)
	mockCache := new(MockCache)
	service := NewService(mockStore, mockCache)

	userID := int64(123)
	expectedTodos := []Todo{
		{ID: 1, UserID: userID, Title: "Todo 1"},
		{ID: 2, UserID: userID, Title: "Todo 2"},
	}

	// Mock cache miss, then database hit
	mockCache.On("Get", mock.Anything, "todos:user:123").Return("", errors.New("cache miss"))
	mockStore.On("GetByUserID", mock.Anything, userID).Return(expectedTodos, nil)
	mockCache.On("Set", mock.Anything, "todos:user:123", mock.AnythingOfType("string"), 10*time.Minute).Return(nil)

	todos, err := service.GetByUserID(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTodos, todos)
	mockStore.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestService_GetByUserID_Error(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore, &NoOpCache{})

	userID := int64(123)
	expectedError := errors.New("database error")
	mockStore.On("GetByUserID", mock.Anything, userID).Return(nil, expectedError)

	todos, err := service.GetByUserID(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, todos)
	assert.Contains(t, err.Error(), "failed to get todos by user id")
	mockStore.AssertExpectations(t)
}

func TestService_Update_Success(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore, &NoOpCache{})

	todoID := int64(1)
	existingTodo := &Todo{
		ID:          todoID,
		UserID:      123,
		Title:       "Old Title",
		Description: "Old Description",
		Completed:   false,
	}

	req := &UpdateTodoRequest{
		Title:       "New Title",
		Description: "New Description",
	}

	mockStore.On("GetByID", mock.Anything, todoID).Return(existingTodo, nil)
	mockStore.On("Save", mock.Anything, existingTodo).Return(nil)

	todo, err := service.Update(context.Background(), todoID, req)

	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, req.Title, todo.Title)
	assert.Equal(t, req.Description, todo.Description)
	assert.False(t, todo.Completed) // Should remain unchanged
	mockStore.AssertExpectations(t)
}

func TestService_Update_TodoNotFound(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore, &NoOpCache{})

	todoID := int64(1)
	req := &UpdateTodoRequest{
		Title:       "New Title",
		Description: "New Description",
	}

	mockStore.On("GetByID", mock.Anything, todoID).Return(nil, ErrTodoNotFound)

	todo, err := service.Update(context.Background(), todoID, req)

	assert.Error(t, err)
	assert.Nil(t, todo)
	assert.Contains(t, err.Error(), "failed to get todo for update")
	mockStore.AssertExpectations(t)
}

func TestService_Update_SaveError(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore, &NoOpCache{})

	todoID := int64(1)
	existingTodo := &Todo{
		ID:     todoID,
		UserID: 123,
		Title:  "Old Title",
	}

	req := &UpdateTodoRequest{
		Title:       "New Title",
		Description: "New Description",
	}

	expectedError := errors.New("database error")
	mockStore.On("GetByID", mock.Anything, todoID).Return(existingTodo, nil)
	mockStore.On("Save", mock.Anything, existingTodo).Return(expectedError)

	todo, err := service.Update(context.Background(), todoID, req)

	assert.Error(t, err)
	assert.Nil(t, todo)
	assert.Contains(t, err.Error(), "failed to update todo")
	mockStore.AssertExpectations(t)
}

func TestService_ToggleComplete_MarkAsCompleted(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore, &NoOpCache{})

	todoID := int64(1)
	existingTodo := &Todo{
		ID:        todoID,
		UserID:    123,
		Title:     "Test Todo",
		Completed: false,
	}

	mockStore.On("GetByID", mock.Anything, todoID).Return(existingTodo, nil)
	mockStore.On("Save", mock.Anything, existingTodo).Return(nil)

	todo, err := service.ToggleComplete(context.Background(), todoID)

	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.True(t, todo.Completed)
	mockStore.AssertExpectations(t)
}

func TestService_ToggleComplete_MarkAsIncomplete(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore, &NoOpCache{})

	todoID := int64(1)
	existingTodo := &Todo{
		ID:        todoID,
		UserID:    123,
		Title:     "Test Todo",
		Completed: true,
	}

	mockStore.On("GetByID", mock.Anything, todoID).Return(existingTodo, nil)
	mockStore.On("Save", mock.Anything, existingTodo).Return(nil)

	todo, err := service.ToggleComplete(context.Background(), todoID)

	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.False(t, todo.Completed)
	mockStore.AssertExpectations(t)
}

func TestService_ToggleComplete_TodoNotFound(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore, &NoOpCache{})

	todoID := int64(1)
	mockStore.On("GetByID", mock.Anything, todoID).Return(nil, ErrTodoNotFound)

	todo, err := service.ToggleComplete(context.Background(), todoID)

	assert.Error(t, err)
	assert.Nil(t, todo)
	assert.Contains(t, err.Error(), "failed to get todo for toggle")
	mockStore.AssertExpectations(t)
}

func TestService_ToggleComplete_SaveError(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore, &NoOpCache{})

	todoID := int64(1)
	existingTodo := &Todo{
		ID:        todoID,
		UserID:    123,
		Title:     "Test Todo",
		Completed: false,
	}

	expectedError := errors.New("database error")
	mockStore.On("GetByID", mock.Anything, todoID).Return(existingTodo, nil)
	mockStore.On("Save", mock.Anything, existingTodo).Return(expectedError)

	todo, err := service.ToggleComplete(context.Background(), todoID)

	assert.Error(t, err)
	assert.Nil(t, todo)
	assert.Contains(t, err.Error(), "failed to toggle todo completion")
	mockStore.AssertExpectations(t)
}

func TestService_Delete_Success(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore, &NoOpCache{})

	todoID := int64(1)
	existingTodo := &Todo{
		ID:     todoID,
		UserID: 123,
		Title:  "Test Todo",
	}

	mockStore.On("GetByID", mock.Anything, todoID).Return(existingTodo, nil)
	mockStore.On("Delete", mock.Anything, todoID).Return(nil)

	err := service.Delete(context.Background(), todoID)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}

func TestService_Delete_Error(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore, &NoOpCache{})

	todoID := int64(1)
	existingTodo := &Todo{
		ID:     todoID,
		UserID: 123,
		Title:  "Test Todo",
	}
	expectedError := errors.New("database error")
	mockStore.On("GetByID", mock.Anything, todoID).Return(existingTodo, nil)
	mockStore.On("Delete", mock.Anything, todoID).Return(expectedError)

	err := service.Delete(context.Background(), todoID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete todo")
	mockStore.AssertExpectations(t)
}

func TestNewService(t *testing.T) {
	mockStore := new(MockStore)
	mockCache := &NoOpCache{}
	service := NewService(mockStore, mockCache)

	assert.NotNil(t, service)
	assert.Equal(t, mockStore, service.store)
	assert.Equal(t, mockCache, service.cache)
}
