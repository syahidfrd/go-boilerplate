package http_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	httpDelivery "github.com/syahidfrd/go-boilerplate/delivery/http"
	"github.com/syahidfrd/go-boilerplate/domain"
	"github.com/syahidfrd/go-boilerplate/domain/mocks"
	"github.com/syahidfrd/go-boilerplate/transport/request"
)

func TestCreate(t *testing.T) {
	mockAuthorUsecase := new(mocks.AuthorUsecase)
	mockCreateAuthorReq := request.CreateAuthorReq{
		Name: "name",
	}

	t.Run("success", func(t *testing.T) {
		jsonReq, err := json.Marshal(mockCreateAuthorReq)
		assert.NoError(t, err)

		mockAuthorUsecase.On("Create", mock.Anything, mock.AnythingOfType("*request.CreateAuthorReq")).
			Return(nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/authors", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.Create(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})

	t.Run("error-validation", func(t *testing.T) {
		mockInvalidCreateAuthorReq := request.CreateAuthorReq{
			Name: "",
		}
		jsonReq, err := json.Marshal(mockInvalidCreateAuthorReq)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/authors", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.Create(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		jsonReq, err := json.Marshal(mockCreateAuthorReq)
		assert.NoError(t, err)

		mockAuthorUsecase.On("Create", mock.Anything, mock.AnythingOfType("*request.CreateAuthorReq")).
			Return(errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/authors", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.Create(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})

}

func TestGetByID(t *testing.T) {
	mockAuthorUsecase := new(mocks.AuthorUsecase)
	mockAuthor := domain.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAuthorUsecase.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(mockAuthor, nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})

	t.Run("data-not-exist", func(t *testing.T) {
		mockAuthorUsecase.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(domain.Author{}, sql.ErrNoRows).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		mockAuthorUsecase.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(domain.Author{}, errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})
}

func TestFetch(t *testing.T) {
	mockAuthorUsecase := new(mocks.AuthorUsecase)
	mockAuthor := domain.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockListAuthor := make([]domain.Author, 0)
	mockListAuthor = append(mockListAuthor, mockAuthor)

	t.Run("success", func(t *testing.T) {
		mockAuthorUsecase.On("Fetch", mock.Anything).Return(mockListAuthor, nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/authors/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.Fetch(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		mockAuthorUsecase.On("Fetch", mock.Anything).Return([]domain.Author{}, errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/authors/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.Fetch(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockAuthorUsecase := new(mocks.AuthorUsecase)
	mockAuthor := domain.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUpdateAuthorReq := request.UpdateAuthorReq{
		Name: "name",
	}

	t.Run("success", func(t *testing.T) {
		jsonReq, err := json.Marshal(mockUpdateAuthorReq)
		assert.NoError(t, err)

		mockAuthorUsecase.On("Update", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*request.UpdateAuthorReq")).
			Return(nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})

	t.Run("error-validation", func(t *testing.T) {
		mockInvalidUpdateAuthorReq := request.UpdateAuthorReq{
			Name: "",
		}
		jsonReq, err := json.Marshal(mockInvalidUpdateAuthorReq)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})

	t.Run("data-not-exist", func(t *testing.T) {
		jsonReq, err := json.Marshal(mockUpdateAuthorReq)
		assert.NoError(t, err)

		mockAuthorUsecase.On("Update", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*request.UpdateAuthorReq")).
			Return(sql.ErrNoRows).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		jsonReq, err := json.Marshal(mockUpdateAuthorReq)
		assert.NoError(t, err)

		mockAuthorUsecase.On("Update", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*request.UpdateAuthorReq")).
			Return(errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})

}

func TestDelete(t *testing.T) {
	mockAuthorUsecase := new(mocks.AuthorUsecase)
	mockAuthor := domain.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAuthorUsecase.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.DELETE, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})

	t.Run("data-not-exist", func(t *testing.T) {
		mockAuthorUsecase.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(sql.ErrNoRows).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.DELETE, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		mockAuthorUsecase.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.DELETE, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUsecase: mockAuthorUsecase,
		}
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockAuthorUsecase.AssertExpectations(t)
	})
}
