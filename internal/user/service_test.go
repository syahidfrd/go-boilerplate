package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) Save(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockStore) FindByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func TestService_Create_Success(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	email := "test@example.com"
	hashedPassword := "hashed_password_123"

	mockStore.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

	user, err := service.Create(context.Background(), email, hashedPassword)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, hashedPassword, user.Password)
	mockStore.AssertExpectations(t)
}

func TestService_Create_SaveError(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	email := "test@example.com"
	hashedPassword := "hashed_password_123"
	expectedError := errors.New("database error")

	mockStore.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(expectedError)

	user, err := service.Create(context.Background(), email, hashedPassword)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to save user")
	mockStore.AssertExpectations(t)
}

func TestService_Create_EmptyEmail(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	email := ""
	hashedPassword := "hashed_password_123"

	mockStore.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

	user, err := service.Create(context.Background(), email, hashedPassword)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, hashedPassword, user.Password)
	mockStore.AssertExpectations(t)
}

func TestService_Create_EmptyPassword(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	email := "test@example.com"
	hashedPassword := ""

	mockStore.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

	user, err := service.Create(context.Background(), email, hashedPassword)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, hashedPassword, user.Password)
	mockStore.AssertExpectations(t)
}

func TestService_GetByEmail_Success(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	email := "test@example.com"
	expectedUser := &User{
		ID:       1,
		Email:    email,
		Password: "hashed_password_123",
	}

	mockStore.On("FindByEmail", mock.Anything, email).Return(expectedUser, nil)

	user, err := service.GetByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockStore.AssertExpectations(t)
}

func TestService_GetByEmail_UserNotFound(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	email := "nonexistent@example.com"
	mockStore.On("FindByEmail", mock.Anything, email).Return(nil, gorm.ErrRecordNotFound)

	user, err := service.GetByEmail(context.Background(), email)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrUserNotFound, err)
	mockStore.AssertExpectations(t)
}

func TestService_GetByEmail_DatabaseError(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	email := "test@example.com"
	expectedError := errors.New("database connection error")
	mockStore.On("FindByEmail", mock.Anything, email).Return(nil, expectedError)

	user, err := service.GetByEmail(context.Background(), email)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to find user by email")
	assert.NotEqual(t, ErrUserNotFound, err)
	mockStore.AssertExpectations(t)
}

func TestService_GetByEmail_EmptyEmail(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	email := ""
	mockStore.On("FindByEmail", mock.Anything, email).Return(nil, gorm.ErrRecordNotFound)

	user, err := service.GetByEmail(context.Background(), email)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrUserNotFound, err)
	mockStore.AssertExpectations(t)
}

func TestService_GetByEmail_SpecialCharacterEmail(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	email := "test+user@example-domain.com"
	expectedUser := &User{
		ID:       1,
		Email:    email,
		Password: "hashed_password_123",
	}

	mockStore.On("FindByEmail", mock.Anything, email).Return(expectedUser, nil)

	user, err := service.GetByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.Equal(t, email, user.Email)
	mockStore.AssertExpectations(t)
}

func TestNewService(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	assert.NotNil(t, service)
	assert.Equal(t, mockStore, service.store)
}

func TestService_Create_VerifyUserConstruction(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	email := "verify@example.com"
	hashedPassword := "verify_password_123"

	// Capture the user that's passed to Save
	var capturedUser *User
	mockStore.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Run(func(args mock.Arguments) {
		capturedUser = args.Get(1).(*User)
	}).Return(nil)

	user, err := service.Create(context.Background(), email, hashedPassword)

	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Verify the user passed to Save is the same as returned
	assert.Equal(t, capturedUser, user)

	// Verify the NewUser constructor was used properly
	assert.Equal(t, email, capturedUser.Email)
	assert.Equal(t, hashedPassword, capturedUser.Password)
	assert.Equal(t, int64(0), capturedUser.ID) // Should be zero before database saves it

	mockStore.AssertExpectations(t)
}

func TestService_ErrorWrapping(t *testing.T) {
	mockStore := new(MockStore)
	service := NewService(mockStore)

	t.Run("Create error wrapping", func(t *testing.T) {
		originalError := errors.New("constraint violation")
		mockStore.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(originalError)

		_, err := service.Create(context.Background(), "test@example.com", "password")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to save user")
		assert.True(t, errors.Is(err, originalError))
	})

	t.Run("GetByEmail error wrapping", func(t *testing.T) {
		originalError := errors.New("connection timeout")
		mockStore.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, originalError)

		_, err := service.GetByEmail(context.Background(), "test@example.com")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find user by email")
		assert.True(t, errors.Is(err, originalError))
	})
}
