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

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	httpDelivery "github.com/syahidfrd/go-boilerplate/delivery/http"
	"github.com/syahidfrd/go-boilerplate/domain"
	"github.com/syahidfrd/go-boilerplate/mocks"
	"github.com/syahidfrd/go-boilerplate/transport/request"
	"github.com/syahidfrd/go-boilerplate/utils"
)

func TestTodoHandler_Create(t *testing.T) {
	mockTodoUC := new(mocks.TodoUsecase)
	createTodoReq := request.CreateTodoReq{
		Name: "name",
	}

	t.Run("success", func(t *testing.T) {
		jsonReq, err := json.Marshal(createTodoReq)
		assert.NoError(t, err)

		mockTodoUC.On("Create", mock.Anything, mock.AnythingOfType("*request.CreateTodoReq")).
			Return(nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/todos", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos")

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.Create(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})

	t.Run("error-validation", func(t *testing.T) {
		invalidCreateTodoReq := request.CreateTodoReq{
			Name: "",
		}
		jsonReq, err := json.Marshal(invalidCreateTodoReq)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/todos", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos")

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.Create(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		jsonReq, err := json.Marshal(createTodoReq)
		assert.NoError(t, err)

		mockTodoUC.On("Create", mock.Anything, mock.AnythingOfType("*request.CreateTodoReq")).
			Return(errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/todos", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos")

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.Create(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})

}

func TestTodoHandler_GetByID(t *testing.T) {
	mockTodoUC := new(mocks.TodoUsecase)
	mockTodo := domain.Todo{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockTodoUC.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(mockTodo, nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/todos/"+strconv.Itoa(int(mockTodo.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockTodo.ID)))

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})

	t.Run("data-not-exist", func(t *testing.T) {
		mockTodoUC.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(domain.Todo{}, utils.NewNotFoundError("todo not found")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/todos/"+strconv.Itoa(int(mockTodo.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockTodo.ID)))

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		mockTodoUC.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(domain.Todo{}, errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/todos/"+strconv.Itoa(int(mockTodo.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockTodo.ID)))

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})
}

func TestTodoHandler_Fetch(t *testing.T) {
	mockTodoUC := new(mocks.TodoUsecase)
	mockTodo := domain.Todo{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockListTodo := make([]domain.Todo, 0)
	mockListTodo = append(mockListTodo, mockTodo)

	t.Run("success", func(t *testing.T) {
		mockTodoUC.On("Fetch", mock.Anything).Return(mockListTodo, nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/todos/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos/")

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.Fetch(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		mockTodoUC.On("Fetch", mock.Anything).Return([]domain.Todo{}, errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/todos/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos/")

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.Fetch(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})
}

func TestTodoHandler_Update(t *testing.T) {
	mockTodoUC := new(mocks.TodoUsecase)
	mockTodo := domain.Todo{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	updateTodoReq := request.UpdateTodoReq{
		Name: "name",
	}

	t.Run("success", func(t *testing.T) {
		jsonReq, err := json.Marshal(updateTodoReq)
		assert.NoError(t, err)

		mockTodoUC.On("Update", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*request.UpdateTodoReq")).
			Return(nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/todos/"+strconv.Itoa(int(mockTodo.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockTodo.ID)))

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})

	t.Run("error-validation", func(t *testing.T) {
		invalidUpdateTodoReq := request.UpdateTodoReq{
			Name: "",
		}
		jsonReq, err := json.Marshal(invalidUpdateTodoReq)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/todos/"+strconv.Itoa(int(mockTodo.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockTodo.ID)))

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})

	t.Run("data-not-exist", func(t *testing.T) {
		jsonReq, err := json.Marshal(updateTodoReq)
		assert.NoError(t, err)

		mockTodoUC.On("Update", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*request.UpdateTodoReq")).
			Return(utils.NewNotFoundError("todo not found")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/todos/"+strconv.Itoa(int(mockTodo.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockTodo.ID)))

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		jsonReq, err := json.Marshal(updateTodoReq)
		assert.NoError(t, err)

		mockTodoUC.On("Update", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*request.UpdateTodoReq")).
			Return(errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/todos/"+strconv.Itoa(int(mockTodo.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockTodo.ID)))

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})

}

func TestTodoHandler_Delete(t *testing.T) {
	mockTodoUC := new(mocks.TodoUsecase)
	mockTodo := domain.Todo{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockTodoUC.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.DELETE, "/api/v1/todos/"+strconv.Itoa(int(mockTodo.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockTodo.ID)))

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})

	t.Run("data-not-exist", func(t *testing.T) {
		mockTodoUC.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(utils.NewNotFoundError("todo not found")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.DELETE, "/api/v1/todos/"+strconv.Itoa(int(mockTodo.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockTodo.ID)))

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		mockTodoUC.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.DELETE, "/api/v1/todos/"+strconv.Itoa(int(mockTodo.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/todos/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockTodo.ID)))

		handler := httpDelivery.TodoHandler{
			TodoUC: mockTodoUC,
		}
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockTodoUC.AssertExpectations(t)
	})
}
