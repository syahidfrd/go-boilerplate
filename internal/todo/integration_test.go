//go:build integration

package todo

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syahidfrd/go-boilerplate/internal/auth"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/cache"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/test"
)

var sharedContainer *test.Container

func TestMain(m *testing.M) {
	var cleanup func() int
	sharedContainer, cleanup = test.SetupTestMain()

	// Run standard migrations + Todo model
	sharedContainer.RunStandardMigrations(&testing.T{})
	err := sharedContainer.DB.AutoMigrate(&Todo{})
	if err != nil {
		panic("failed to migrate Todo model: " + err.Error())
	}

	code := m.Run()
	os.Exit(cleanup() + code)
}

func setupTestServices(t *testing.T) (*Service, *handler, *test.Container) {
	t.Helper()

	// Clean all data before each test
	sharedContainer.CleanupAll(t)

	store := NewStore(sharedContainer.DB)
	redisCache := cache.NewRedis(sharedContainer.Redis)
	service := NewService(store, redisCache)
	handler := NewHandler(service)

	return service, handler, sharedContainer
}

func createAuthenticatedContext(userID int64) context.Context {
	return context.WithValue(context.Background(), auth.UserIDKey, userID)
}

func TestTodoCreateIntegration(t *testing.T) {
	_, handler, _ := setupTestServices(t)

	tests := []struct {
		name           string
		userID         int64
		requestBody    CreateTodoRequest
		expectedStatus int
		checkResponse  func(t *testing.T, resp *test.HTTPResponse)
	}{
		{
			name:   "successful todo creation",
			userID: 1,
			requestBody: CreateTodoRequest{
				Title:       "Test Todo",
				Description: "Test Description",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp *test.HTTPResponse) {
				require.NotNil(t, resp.Body)
				assert.Equal(t, "Test Todo", resp.Body["title"])
				assert.Equal(t, "Test Description", resp.Body["description"])
				assert.Equal(t, false, resp.Body["completed"])
				assert.NotNil(t, resp.Body["id"])
				assert.Equal(t, float64(1), resp.Body["user_id"])
			},
		},
		{
			name:   "missing title",
			userID: 1,
			requestBody: CreateTodoRequest{
				Description: "Test Description",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "empty title",
			userID: 1,
			requestBody: CreateTodoRequest{
				Title:       "",
				Description: "Test Description",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createAuthenticatedContext(tt.userID)

			resp := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
				handler.Create(w, r.WithContext(ctx))
			}, test.HTTPRequest{
				Method: http.MethodPost,
				URL:    "/todos",
				Body:   tt.requestBody,
			})

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}
		})
	}
}

func TestTodoCreateUnauthorizedIntegration(t *testing.T) {
	_, handler, _ := setupTestServices(t)

	resp := test.MakeJSONRequest(t, handler.Create, test.HTTPRequest{
		Method: http.MethodPost,
		URL:    "/todos",
		Body: CreateTodoRequest{
			Title:       "Test Todo",
			Description: "Test Description",
		},
	})

	test.AssertErrorResponse(t, resp, http.StatusUnauthorized, "unauthorized")
}

func TestTodoGetByUserIDIntegration(t *testing.T) {
	service, handler, _ := setupTestServices(t)

	userID := int64(1)
	ctx := createAuthenticatedContext(userID)

	// Create some test todos
	_, err := service.Create(context.Background(), userID, &CreateTodoRequest{
		Title:       "Todo 1",
		Description: "Description 1",
	})
	require.NoError(t, err)

	_, err = service.Create(context.Background(), userID, &CreateTodoRequest{
		Title:       "Todo 2",
		Description: "Description 2",
	})
	require.NoError(t, err)

	resp := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		handler.GetByUserID(w, r.WithContext(ctx))
	}, test.HTTPRequest{
		Method: http.MethodGet,
		URL:    "/todos",
	})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, resp.Body)

	data, ok := resp.Body["data"].([]any)
	require.True(t, ok)
	assert.Len(t, data, 2)

	// Check if both todos are present
	foundTodos := make(map[string]bool)
	for _, item := range data {
		todoMap := item.(map[string]any)
		foundTodos[todoMap["title"].(string)] = true
	}
	assert.True(t, foundTodos["Todo 1"])
	assert.True(t, foundTodos["Todo 2"])

	// Verify todos have correct userID
	for _, item := range data {
		todoMap := item.(map[string]any)
		assert.Equal(t, float64(userID), todoMap["user_id"])
	}
}

func TestTodoGetByUserIDCachingIntegration(t *testing.T) {
	service, handler, _ := setupTestServices(t)

	userID := int64(1)
	ctx := createAuthenticatedContext(userID)

	// Create a todo
	_, err := service.Create(context.Background(), userID, &CreateTodoRequest{
		Title:       "Cached Todo",
		Description: "Description",
	})
	require.NoError(t, err)

	// First request should populate cache
	resp1 := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		handler.GetByUserID(w, r.WithContext(ctx))
	}, test.HTTPRequest{
		Method: http.MethodGet,
		URL:    "/todos",
	})

	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	require.NotNil(t, resp1.Body)

	data1, ok := resp1.Body["data"].([]any)
	require.True(t, ok)
	assert.Len(t, data1, 1)

	// Second request should hit cache
	resp2 := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		handler.GetByUserID(w, r.WithContext(ctx))
	}, test.HTTPRequest{
		Method: http.MethodGet,
		URL:    "/todos",
	})

	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	require.NotNil(t, resp2.Body)

	data2, ok := resp2.Body["data"].([]any)
	require.True(t, ok)
	assert.Len(t, data2, 1)

	// Verify core fields match (timestamps might differ due to caching serialization)
	todo1 := data1[0].(map[string]any)
	todo2 := data2[0].(map[string]any)

	assert.Equal(t, todo1["title"], todo2["title"])
	assert.Equal(t, todo1["description"], todo2["description"])
	assert.Equal(t, todo1["completed"], todo2["completed"])
	assert.Equal(t, todo1["id"], todo2["id"])
}

func TestTodoGetByIDIntegration(t *testing.T) {
	service, handler, _ := setupTestServices(t)

	userID := int64(1)

	// Create a test todo
	todo, err := service.Create(context.Background(), userID, &CreateTodoRequest{
		Title:       "Get By ID Todo",
		Description: "Get By ID Description",
	})
	require.NoError(t, err)

	resp := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		// Simulate path value
		r.SetPathValue("id", strconv.FormatInt(todo.ID, 10))
		handler.GetByID(w, r)
	}, test.HTTPRequest{
		Method: http.MethodGet,
		URL:    "/todos/" + strconv.FormatInt(todo.ID, 10),
	})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, resp.Body)
	assert.Equal(t, "Get By ID Todo", resp.Body["title"])
	assert.Equal(t, "Get By ID Description", resp.Body["description"])
	assert.Equal(t, float64(todo.ID), resp.Body["id"])
}

func TestTodoGetByIDNotFoundIntegration(t *testing.T) {
	_, handler, _ := setupTestServices(t)

	resp := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		r.SetPathValue("id", "999")
		handler.GetByID(w, r)
	}, test.HTTPRequest{
		Method: http.MethodGet,
		URL:    "/todos/999",
	})

	test.AssertErrorResponse(t, resp, http.StatusNotFound, "todo not found")
}

func TestTodoUpdateIntegration(t *testing.T) {
	service, handler, _ := setupTestServices(t)

	userID := int64(1)

	// Create a test todo
	todo, err := service.Create(context.Background(), userID, &CreateTodoRequest{
		Title:       "Original Title",
		Description: "Original Description",
	})
	require.NoError(t, err)

	updateReq := UpdateTodoRequest{
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	resp := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		r.SetPathValue("id", strconv.FormatInt(todo.ID, 10))
		handler.Update(w, r)
	}, test.HTTPRequest{
		Method: http.MethodPut,
		URL:    "/todos/" + strconv.FormatInt(todo.ID, 10),
		Body:   updateReq,
	})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, resp.Body)
	assert.Equal(t, "Updated Title", resp.Body["title"])
	assert.Equal(t, "Updated Description", resp.Body["description"])
	assert.Equal(t, float64(todo.ID), resp.Body["id"])
}

func TestTodoUpdateNotFoundIntegration(t *testing.T) {
	_, handler, _ := setupTestServices(t)

	updateReq := UpdateTodoRequest{
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	resp := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		r.SetPathValue("id", "999")
		handler.Update(w, r)
	}, test.HTTPRequest{
		Method: http.MethodPut,
		URL:    "/todos/999",
		Body:   updateReq,
	})

	test.AssertErrorResponse(t, resp, http.StatusNotFound, "todo not found")
}

func TestTodoToggleCompleteIntegration(t *testing.T) {
	service, handler, _ := setupTestServices(t)

	userID := int64(1)

	// Create a test todo
	todo, err := service.Create(context.Background(), userID, &CreateTodoRequest{
		Title:       "Toggle Todo",
		Description: "Toggle Description",
	})
	require.NoError(t, err)
	assert.False(t, todo.Completed)

	// Toggle to completed
	resp1 := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		r.SetPathValue("id", strconv.FormatInt(todo.ID, 10))
		handler.ToggleComplete(w, r)
	}, test.HTTPRequest{
		Method: http.MethodPatch,
		URL:    "/todos/" + strconv.FormatInt(todo.ID, 10) + "/toggle",
	})

	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	require.NotNil(t, resp1.Body)
	assert.Equal(t, true, resp1.Body["completed"])

	// Toggle back to incomplete
	resp2 := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		r.SetPathValue("id", strconv.FormatInt(todo.ID, 10))
		handler.ToggleComplete(w, r)
	}, test.HTTPRequest{
		Method: http.MethodPatch,
		URL:    "/todos/" + strconv.FormatInt(todo.ID, 10) + "/toggle",
	})

	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	require.NotNil(t, resp2.Body)
	assert.Equal(t, false, resp2.Body["completed"])
}

func TestTodoToggleCompleteNotFoundIntegration(t *testing.T) {
	_, handler, _ := setupTestServices(t)

	resp := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		r.SetPathValue("id", "999")
		handler.ToggleComplete(w, r)
	}, test.HTTPRequest{
		Method: http.MethodPatch,
		URL:    "/todos/999/toggle",
	})

	test.AssertErrorResponse(t, resp, http.StatusNotFound, "todo not found")
}

func TestTodoDeleteIntegration(t *testing.T) {
	service, handler, _ := setupTestServices(t)

	userID := int64(1)

	// Create a test todo
	todo, err := service.Create(context.Background(), userID, &CreateTodoRequest{
		Title:       "Delete Todo",
		Description: "Delete Description",
	})
	require.NoError(t, err)

	// Delete the todo
	resp := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		r.SetPathValue("id", strconv.FormatInt(todo.ID, 10))
		handler.Delete(w, r)
	}, test.HTTPRequest{
		Method: http.MethodDelete,
		URL:    "/todos/" + strconv.FormatInt(todo.ID, 10),
	})

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify todo is deleted
	getResp := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		r.SetPathValue("id", strconv.FormatInt(todo.ID, 10))
		handler.GetByID(w, r)
	}, test.HTTPRequest{
		Method: http.MethodGet,
		URL:    "/todos/" + strconv.FormatInt(todo.ID, 10),
	})

	test.AssertErrorResponse(t, getResp, http.StatusNotFound, "todo not found")
}

func TestTodoDeleteNotFoundIntegration(t *testing.T) {
	_, handler, _ := setupTestServices(t)

	resp := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		r.SetPathValue("id", "999")
		handler.Delete(w, r)
	}, test.HTTPRequest{
		Method: http.MethodDelete,
		URL:    "/todos/999",
	})

	test.AssertErrorResponse(t, resp, http.StatusNotFound, "todo not found")
}

func TestTodoCacheInvalidationIntegration(t *testing.T) {
	service, handler, _ := setupTestServices(t)

	userID := int64(1)
	ctx := createAuthenticatedContext(userID)

	// Create initial todo and populate cache
	_, err := service.Create(context.Background(), userID, &CreateTodoRequest{
		Title:       "Cache Test Todo",
		Description: "Cache Test Description",
	})
	require.NoError(t, err)

	// Get todos to populate cache
	resp1 := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		handler.GetByUserID(w, r.WithContext(ctx))
	}, test.HTTPRequest{
		Method: http.MethodGet,
		URL:    "/todos",
	})

	require.Equal(t, http.StatusOK, resp1.StatusCode)
	data1 := resp1.Body["data"].([]any)
	assert.Len(t, data1, 1)

	// Create another todo (should invalidate cache)
	_, err = service.Create(context.Background(), userID, &CreateTodoRequest{
		Title:       "Second Todo",
		Description: "Second Description",
	})
	require.NoError(t, err)

	// Get todos again - should show updated data
	resp2 := test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
		handler.GetByUserID(w, r.WithContext(ctx))
	}, test.HTTPRequest{
		Method: http.MethodGet,
		URL:    "/todos",
	})

	require.Equal(t, http.StatusOK, resp2.StatusCode)
	data2 := resp2.Body["data"].([]any)
	assert.Len(t, data2, 2)
}

func TestTodoServiceDirectCallsIntegration(t *testing.T) {
	service, _, _ := setupTestServices(t)

	userID := int64(1)
	ctx := context.Background()

	// Test Create
	createReq := &CreateTodoRequest{
		Title:       "Service Test Todo",
		Description: "Service Test Description",
	}

	todo, err := service.Create(ctx, userID, createReq)
	require.NoError(t, err)
	assert.NotZero(t, todo.ID)
	assert.Equal(t, userID, todo.UserID)
	assert.Equal(t, createReq.Title, todo.Title)
	assert.Equal(t, createReq.Description, todo.Description)
	assert.False(t, todo.Completed)

	// Test GetByID
	foundTodo, err := service.GetByID(ctx, todo.ID)
	require.NoError(t, err)
	assert.Equal(t, todo.ID, foundTodo.ID)
	assert.Equal(t, todo.Title, foundTodo.Title)

	// Test GetByUserID
	todos, err := service.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, todos, 1)
	assert.Equal(t, todo.ID, todos[0].ID)

	// Test Update
	updateReq := &UpdateTodoRequest{
		Title:       "Updated Service Todo",
		Description: "Updated Service Description",
	}

	updatedTodo, err := service.Update(ctx, todo.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, updateReq.Title, updatedTodo.Title)
	assert.Equal(t, updateReq.Description, updatedTodo.Description)

	// Test ToggleComplete
	toggledTodo, err := service.ToggleComplete(ctx, todo.ID)
	require.NoError(t, err)
	assert.True(t, toggledTodo.Completed)

	// Toggle again
	toggledTodo, err = service.ToggleComplete(ctx, todo.ID)
	require.NoError(t, err)
	assert.False(t, toggledTodo.Completed)

	// Test Delete
	err = service.Delete(ctx, todo.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = service.GetByID(ctx, todo.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrTodoNotFound)
}

func TestTodoUserIsolationIntegration(t *testing.T) {
	service, _, _ := setupTestServices(t)

	user1ID := int64(1)
	user2ID := int64(2)
	ctx := context.Background()

	// Create todos for different users
	todo1, err := service.Create(ctx, user1ID, &CreateTodoRequest{
		Title:       "User 1 Todo",
		Description: "User 1 Description",
	})
	require.NoError(t, err)

	todo2, err := service.Create(ctx, user2ID, &CreateTodoRequest{
		Title:       "User 2 Todo",
		Description: "User 2 Description",
	})
	require.NoError(t, err)

	// Verify user 1 only sees their todo
	user1Todos, err := service.GetByUserID(ctx, user1ID)
	require.NoError(t, err)
	assert.Len(t, user1Todos, 1)
	assert.Equal(t, todo1.ID, user1Todos[0].ID)
	assert.Equal(t, user1ID, user1Todos[0].UserID)

	// Verify user 2 only sees their todo
	user2Todos, err := service.GetByUserID(ctx, user2ID)
	require.NoError(t, err)
	assert.Len(t, user2Todos, 1)
	assert.Equal(t, todo2.ID, user2Todos[0].ID)
	assert.Equal(t, user2ID, user2Todos[0].UserID)
}

func TestTodoJSONMarshalingIntegration(t *testing.T) {
	service, _, _ := setupTestServices(t)

	userID := int64(1)
	ctx := context.Background()

	todo, err := service.Create(ctx, userID, &CreateTodoRequest{
		Title:       "JSON Test Todo",
		Description: "JSON Test Description",
	})
	require.NoError(t, err)

	// Test JSON marshaling
	jsonData, err := todo.MarshalJSON()
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "JSON Test Todo")
	assert.Contains(t, string(jsonData), "JSON Test Description")
	assert.Contains(t, string(jsonData), "\"completed\":false")
}
