package http_test

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	httpDelivery "github.com/syahidfrd/go-boilerplate/delivery/http"
	"github.com/syahidfrd/go-boilerplate/mocks"
	"github.com/syahidfrd/go-boilerplate/transport/request"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthHandler_SignUp(t *testing.T) {
	mockAuthUC := new(mocks.AuthUsecase)
	signUpReq := request.SignUpReq{
		Email:    "sample@mail.com",
		Password: "12345678",
	}

	t.Run("success", func(t *testing.T) {
		jsonReq, err := json.Marshal(signUpReq)
		assert.NoError(t, err)

		mockAuthUC.On("SignUp", mock.Anything, mock.AnythingOfType("*request.SignUpReq")).
			Return(nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/auth/signup", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/signup")

		handler := httpDelivery.AuthHandler{
			AuthUC: mockAuthUC,
		}
		err = handler.SignUp(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockAuthUC.AssertExpectations(t)
	})

	t.Run("error-validation", func(t *testing.T) {
		invalidSignUpReq := request.SignUpReq{
			Email: "sample@mail.com",
		}
		jsonReq, err := json.Marshal(invalidSignUpReq)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/auth/signup", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/signup")

		handler := httpDelivery.AuthHandler{
			AuthUC: mockAuthUC,
		}
		err = handler.SignUp(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockAuthUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		jsonReq, err := json.Marshal(signUpReq)
		assert.NoError(t, err)

		mockAuthUC.On("SignUp", mock.Anything, mock.AnythingOfType("*request.SignUpReq")).
			Return(errors.New("unexpected error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/auth/signup", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/signup")

		handler := httpDelivery.AuthHandler{
			AuthUC: mockAuthUC,
		}
		err = handler.SignUp(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockAuthUC.AssertExpectations(t)
	})
}

func TestAuthHandler_SignIn(t *testing.T) {
	mockAuthUC := new(mocks.AuthUsecase)
	signInReq := request.SignInReq{
		Email:    "sample@mail.com",
		Password: "12345678",
	}

	t.Run("success", func(t *testing.T) {
		jsonReq, err := json.Marshal(signInReq)
		assert.NoError(t, err)

		mockAuthUC.On("SignIn", mock.Anything, mock.AnythingOfType("*request.SignInReq")).
			Return("accessToken", nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/auth/signin", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/signin")

		handler := httpDelivery.AuthHandler{
			AuthUC: mockAuthUC,
		}
		err = handler.SignIn(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockAuthUC.AssertExpectations(t)
	})

	t.Run("error-validation", func(t *testing.T) {
		invalidSignInReq := request.SignInReq{
			Email: "sample@mail.com",
		}
		jsonReq, err := json.Marshal(invalidSignInReq)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/auth/signin", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/signin")

		handler := httpDelivery.AuthHandler{
			AuthUC: mockAuthUC,
		}
		err = handler.SignIn(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockAuthUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		jsonReq, err := json.Marshal(signInReq)
		assert.NoError(t, err)

		mockAuthUC.On("SignIn", mock.Anything, mock.AnythingOfType("*request.SignInReq")).
			Return("", errors.New("unexpected error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/auth/signin", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/signin")

		handler := httpDelivery.AuthHandler{
			AuthUC: mockAuthUC,
		}
		err = handler.SignIn(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockAuthUC.AssertExpectations(t)
	})
}
