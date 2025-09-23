package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/syahidfrd/go-boilerplate/internal/jwt"
	"github.com/syahidfrd/go-boilerplate/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, email, hashedPassword string) (*user.User, error) {
	args := m.Called(ctx, email, hashedPassword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID int64) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Claims), args.Error(1)
}

func TestAuthService_SignUp_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	mockJWTService := new(MockJWTService)
	authService := NewService(mockUserService, mockJWTService)

	req := &SignUpRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedUser := &user.User{
		ID:    1,
		Email: "test@example.com",
	}

	mockUserService.On("Create", mock.Anything, req.Email, mock.AnythingOfType("string")).Return(expectedUser, nil)

	response, err := authService.SignUp(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "signup successfully", response.Message)
	mockUserService.AssertExpectations(t)
}

func TestAuthService_SignUp_UserServiceError(t *testing.T) {
	mockUserService := new(MockUserService)
	mockJWTService := new(MockJWTService)
	authService := NewService(mockUserService, mockJWTService)

	req := &SignUpRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedError := errors.New("user already exists")
	mockUserService.On("Create", mock.Anything, req.Email, mock.AnythingOfType("string")).Return(nil, expectedError)

	response, err := authService.SignUp(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, ErrUserAlreadyExists, err)
	mockUserService.AssertExpectations(t)
}

func TestAuthService_SignUp_HashesPassword(t *testing.T) {
	mockUserService := new(MockUserService)
	mockJWTService := new(MockJWTService)
	authService := NewService(mockUserService, mockJWTService)

	req := &SignUpRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	var capturedHashedPassword string
	mockUserService.On("Create", mock.Anything, req.Email, mock.AnythingOfType("string")).Run(func(args mock.Arguments) {
		capturedHashedPassword = args.String(2)
	}).Return(&user.User{ID: 1}, nil)

	_, err := authService.SignUp(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEqual(t, req.Password, capturedHashedPassword)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(capturedHashedPassword), []byte(req.Password)))
	mockUserService.AssertExpectations(t)
}

func TestAuthService_SignIn_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	mockJWTService := new(MockJWTService)
	authService := NewService(mockUserService, mockJWTService)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	req := &SignInRequest{
		Email:    "test@example.com",
		Password: password,
	}

	expectedUser := &user.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	expectedToken := "jwt-token-123"

	mockUserService.On("GetByEmail", mock.Anything, req.Email).Return(expectedUser, nil)
	mockJWTService.On("GenerateToken", expectedUser.ID).Return(expectedToken, nil)

	response, err := authService.SignIn(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedToken, response.AccessToken)
	mockUserService.AssertExpectations(t)
	mockJWTService.AssertExpectations(t)
}

func TestAuthService_SignIn_UserNotFound(t *testing.T) {
	mockUserService := new(MockUserService)
	mockJWTService := new(MockJWTService)
	authService := NewService(mockUserService, mockJWTService)

	req := &SignInRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockUserService.On("GetByEmail", mock.Anything, req.Email).Return(nil, user.ErrUserNotFound)

	response, err := authService.SignIn(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, ErrInvalidCredentials, err)
	mockUserService.AssertExpectations(t)
}

func TestAuthService_SignIn_InvalidPassword(t *testing.T) {
	mockUserService := new(MockUserService)
	mockJWTService := new(MockJWTService)
	authService := NewService(mockUserService, mockJWTService)

	correctPassword := "password123"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	req := &SignInRequest{
		Email:    "test@example.com",
		Password: wrongPassword,
	}

	expectedUser := &user.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	mockUserService.On("GetByEmail", mock.Anything, req.Email).Return(expectedUser, nil)

	response, err := authService.SignIn(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, ErrInvalidCredentials, err)
	mockUserService.AssertExpectations(t)
}

func TestAuthService_SignIn_JWTGenerationError(t *testing.T) {
	mockUserService := new(MockUserService)
	mockJWTService := new(MockJWTService)
	authService := NewService(mockUserService, mockJWTService)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	req := &SignInRequest{
		Email:    "test@example.com",
		Password: password,
	}

	expectedUser := &user.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	expectedError := errors.New("jwt generation failed")

	mockUserService.On("GetByEmail", mock.Anything, req.Email).Return(expectedUser, nil)
	mockJWTService.On("GenerateToken", expectedUser.ID).Return("", expectedError)

	response, err := authService.SignIn(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to generate token")
	assert.Contains(t, err.Error(), expectedError.Error())
	mockUserService.AssertExpectations(t)
	mockJWTService.AssertExpectations(t)
}
