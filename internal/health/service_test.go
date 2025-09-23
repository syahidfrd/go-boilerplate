package health

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) PingDatabase(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockStore) PingCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestNewService(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	assert.NotNil(t, service)
	assert.Equal(t, mockStore, service.store)
}

func TestService_Check_AllHealthy(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	// Mock successful ping responses
	mockStore.On("PingDatabase", mock.Anything).Return(nil)
	mockStore.On("PingCache", mock.Anything).Return(nil)

	response, err := service.Check(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, StatusHealthy, response.Database)
	assert.Equal(t, StatusHealthy, response.Cache)

	mockStore.AssertExpectations(t)
}

func TestService_Check_DatabaseUnhealthy(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	// Mock database failure
	mockStore.On("PingDatabase", mock.Anything).Return(errors.New("database connection failed"))
	mockStore.On("PingCache", mock.Anything).Return(nil)

	response, err := service.Check(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, StatusUnhealthy, response.Database)
	assert.Equal(t, StatusHealthy, response.Cache)

	mockStore.AssertExpectations(t)
}

func TestService_Check_CacheUnhealthy(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	// Mock cache failure
	mockStore.On("PingDatabase", mock.Anything).Return(nil)
	mockStore.On("PingCache", mock.Anything).Return(errors.New("redis connection failed"))

	response, err := service.Check(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, StatusHealthy, response.Database)
	assert.Equal(t, StatusUnhealthy, response.Cache)

	mockStore.AssertExpectations(t)
}
