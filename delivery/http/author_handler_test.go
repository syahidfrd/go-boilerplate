package http_test

import (
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
	"github.com/syahidfrd/go-boilerplate/entity"
	"github.com/syahidfrd/go-boilerplate/mocks"
	"github.com/syahidfrd/go-boilerplate/transport/request"
	"github.com/syahidfrd/go-boilerplate/utils"
)

func TestAuthorHandler_Create(t *testing.T) {
	mockAuthorUC := new(mocks.AuthorUsecase)
	createAuthorReq := request.CreateAuthorReq{
		Name: "name",
	}

	t.Run("success", func(t *testing.T) {
		jsonReq, err := json.Marshal(createAuthorReq)
		assert.NoError(t, err)

		mockAuthorUC.On("Create", mock.Anything, mock.AnythingOfType("*request.CreateAuthorReq")).
			Return(nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/authors", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors")

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.Create(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})

	t.Run("error-validation", func(t *testing.T) {
		invalidCreateAuthorReq := request.CreateAuthorReq{
			Name: "",
		}
		jsonReq, err := json.Marshal(invalidCreateAuthorReq)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/authors", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors")

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.Create(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		jsonReq, err := json.Marshal(createAuthorReq)
		assert.NoError(t, err)

		mockAuthorUC.On("Create", mock.Anything, mock.AnythingOfType("*request.CreateAuthorReq")).
			Return(errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/authors", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors")

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.Create(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})

}

func TestAuthorHandler_GetByID(t *testing.T) {
	mockAuthorUC := new(mocks.AuthorUsecase)
	mockAuthor := entity.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAuthorUC.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(mockAuthor, nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})

	t.Run("data-not-exist", func(t *testing.T) {
		mockAuthorUC.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(entity.Author{}, utils.NewNotFoundError("author not found")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		mockAuthorUC.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(entity.Author{}, errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})
}

func TestAuthorHandler_Fetch(t *testing.T) {
	mockAuthorUC := new(mocks.AuthorUsecase)
	mockAuthor := entity.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockListAuthor := make([]entity.Author, 0)
	mockListAuthor = append(mockListAuthor, mockAuthor)

	t.Run("success", func(t *testing.T) {
		mockAuthorUC.On("Fetch", mock.Anything).Return(mockListAuthor, nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/authors/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors/")

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.Fetch(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		mockAuthorUC.On("Fetch", mock.Anything).Return([]entity.Author{}, errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/authors/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors/")

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.Fetch(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})
}

func TestAuthorHandler_Update(t *testing.T) {
	mockAuthorUC := new(mocks.AuthorUsecase)
	mockAuthor := entity.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	updateAuthorReq := request.UpdateAuthorReq{
		Name: "name",
	}

	t.Run("success", func(t *testing.T) {
		jsonReq, err := json.Marshal(updateAuthorReq)
		assert.NoError(t, err)

		mockAuthorUC.On("Update", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*request.UpdateAuthorReq")).
			Return(nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})

	t.Run("error-validation", func(t *testing.T) {
		invalidUpdateAuthorReq := request.UpdateAuthorReq{
			Name: "",
		}
		jsonReq, err := json.Marshal(invalidUpdateAuthorReq)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})

	t.Run("data-not-exist", func(t *testing.T) {
		jsonReq, err := json.Marshal(updateAuthorReq)
		assert.NoError(t, err)

		mockAuthorUC.On("Update", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*request.UpdateAuthorReq")).
			Return(utils.NewNotFoundError("author not found")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		jsonReq, err := json.Marshal(updateAuthorReq)
		assert.NoError(t, err)

		mockAuthorUC.On("Update", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*request.UpdateAuthorReq")).
			Return(errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})

}

func TestAuthorHandler_Delete(t *testing.T) {
	mockAuthorUC := new(mocks.AuthorUsecase)
	mockAuthor := entity.Author{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAuthorUC.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.DELETE, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})

	t.Run("data-not-exist", func(t *testing.T) {
		mockAuthorUC.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(utils.NewNotFoundError("author not found")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.DELETE, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		mockAuthorUC.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.DELETE, "/api/v1/authors/"+strconv.Itoa(int(mockAuthor.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/authors/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockAuthor.ID)))

		handler := httpDelivery.AuthorHandler{
			AuthorUC: mockAuthorUC,
		}
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockAuthorUC.AssertExpectations(t)
	})
}
