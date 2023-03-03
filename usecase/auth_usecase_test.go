package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/syahidfrd/go-boilerplate/domain"
	"github.com/syahidfrd/go-boilerplate/mocks"
	"github.com/syahidfrd/go-boilerplate/transport/request"
	"github.com/syahidfrd/go-boilerplate/usecase"
)

func TestAuthUC_SignUp(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockCryptoSvc := new(mocks.CryptoService)
	mockJWTSvc := new(mocks.JWTService)
	signUpReq := request.SignUpReq{
		Email:    "sample@mail.com",
		Password: "12345678",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()
		mockCryptoSvc.On("CreatePasswordHash", mock.Anything, mock.AnythingOfType("string")).Return("passwordHash", nil).Once()
		mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		authUC := usecase.NewAuthUsecase(mockUserRepo, mockCryptoSvc, mockJWTSvc, 60*time.Second)
		err := authUC.SignUp(context.TODO(), &signUpReq)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
		mockCryptoSvc.AssertExpectations(t)
		mockJWTSvc.AssertExpectations(t)
	})

	t.Run("email-already-registered", func(t *testing.T) {
		mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{ID: 1}, nil).Once()

		authUC := usecase.NewAuthUsecase(mockUserRepo, mockCryptoSvc, mockJWTSvc, 60*time.Second)
		err := authUC.SignUp(context.TODO(), &signUpReq)

		assert.NotNil(t, err)
		mockUserRepo.AssertExpectations(t)
		mockCryptoSvc.AssertExpectations(t)
		mockJWTSvc.AssertExpectations(t)
	})

	t.Run("error-password-hash", func(t *testing.T) {
		mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()
		mockCryptoSvc.On("CreatePasswordHash", mock.Anything, mock.AnythingOfType("string")).Return("", errors.New("unexpected error")).Once()

		authUC := usecase.NewAuthUsecase(mockUserRepo, mockCryptoSvc, mockJWTSvc, 60*time.Second)
		err := authUC.SignUp(context.TODO(), &signUpReq)

		assert.NotNil(t, err)
		mockUserRepo.AssertExpectations(t)
		mockCryptoSvc.AssertExpectations(t)
		mockJWTSvc.AssertExpectations(t)
	})

	t.Run("error-create-new-user", func(t *testing.T) {
		mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()
		mockCryptoSvc.On("CreatePasswordHash", mock.Anything, mock.AnythingOfType("string")).Return("passwordHash", nil).Once()
		mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("unexpected error")).Once()

		authUC := usecase.NewAuthUsecase(mockUserRepo, mockCryptoSvc, mockJWTSvc, 60*time.Second)
		err := authUC.SignUp(context.TODO(), &signUpReq)

		assert.NotNil(t, err)
		mockUserRepo.AssertExpectations(t)
		mockCryptoSvc.AssertExpectations(t)
		mockJWTSvc.AssertExpectations(t)
	})
}

func TestAuthUC_SignIn(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockCryptoSvc := new(mocks.CryptoService)
	mockJWTSvc := new(mocks.JWTService)
	mockUser := domain.User{
		ID:        1,
		Email:     "sample@mail.com",
		Password:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	signInReq := request.SignInReq{
		Email:    "sample@mail.com",
		Password: "12345678",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(mockUser, nil).Once()
		mockCryptoSvc.On("ValidatePassword", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(true).Once()
		mockJWTSvc.On("GenerateToken", mock.Anything, mock.AnythingOfType("int64")).Return("accessToken", nil).Once()

		authUC := usecase.NewAuthUsecase(mockUserRepo, mockCryptoSvc, mockJWTSvc, 60*time.Second)
		_, err := authUC.SignIn(context.TODO(), &signInReq)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
		mockCryptoSvc.AssertExpectations(t)
		mockJWTSvc.AssertExpectations(t)
	})

	t.Run("invalid-password", func(t *testing.T) {
		mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(mockUser, nil).Once()
		mockCryptoSvc.On("ValidatePassword", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(false).Once()

		authUC := usecase.NewAuthUsecase(mockUserRepo, mockCryptoSvc, mockJWTSvc, 60*time.Second)
		_, err := authUC.SignIn(context.TODO(), &signInReq)

		assert.NotNil(t, err)
		mockUserRepo.AssertExpectations(t)
		mockCryptoSvc.AssertExpectations(t)
		mockJWTSvc.AssertExpectations(t)
	})

}
