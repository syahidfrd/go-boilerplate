//go:build integration

package health

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/test"
)

var sharedContainer *test.Container

func TestMain(m *testing.M) {
	var cleanup func() int
	sharedContainer, cleanup = test.SetupTestMain()

	// Run standard migrations
	sharedContainer.RunStandardMigrations(&testing.T{})

	code := m.Run()
	os.Exit(cleanup() + code)
}

func setupTestServices(t *testing.T) (*Service, *Handler, *test.Container) {
	t.Helper()

	// Clean all data before each test
	sharedContainer.CleanupAll(t)

	store := NewStore(sharedContainer.DB, sharedContainer.Redis)
	service := NewService(store)
	handler := NewHandler(service)

	return service, handler, sharedContainer
}

func TestHealthCheckIntegration(t *testing.T) {
	_, handler, _ := setupTestServices(t)

	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "healthy services",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"database": "healthy",
				"cache":    "healthy",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := test.MakeJSONRequest(t, handler.Health, test.HTTPRequest{
				Method: http.MethodGet,
				URL:    "/health",
			})

			test.AssertJSONResponse(t, resp, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func TestHealthCheckDatabaseDownIntegration(t *testing.T) {
	// Clean cache before test
	sharedContainer.CleanupCache(t)

	// Create store with nil database to simulate database being down
	store := NewStore(nil, sharedContainer.Redis)
	service := NewService(store)
	handler := NewHandler(service)

	resp := test.MakeJSONRequest(t, handler.Health, test.HTTPRequest{
		Method: http.MethodGet,
		URL:    "/health",
	})

	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	require.NotNil(t, resp.Body)
	assert.Equal(t, "unhealthy", resp.Body["database"])
	assert.Equal(t, "healthy", resp.Body["cache"])
}

func TestHealthCheckCacheDownIntegration(t *testing.T) {
	// Clean database before test
	sharedContainer.CleanupDatabase(t)

	// Create store with nil Redis to simulate cache being down
	store := NewStore(sharedContainer.DB, nil)
	service := NewService(store)
	handler := NewHandler(service)

	resp := test.MakeJSONRequest(t, handler.Health, test.HTTPRequest{
		Method: http.MethodGet,
		URL:    "/health",
	})

	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	require.NotNil(t, resp.Body)
	assert.Equal(t, "healthy", resp.Body["database"])
	assert.Equal(t, "unhealthy", resp.Body["cache"])
}

func TestHealthCheckBothServicesDownIntegration(t *testing.T) {
	// Create store with both nil connections to simulate both services being down
	store := NewStore(nil, nil)
	service := NewService(store)
	handler := NewHandler(service)

	resp := test.MakeJSONRequest(t, handler.Health, test.HTTPRequest{
		Method: http.MethodGet,
		URL:    "/health",
	})

	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	require.NotNil(t, resp.Body)
	assert.Equal(t, "unhealthy", resp.Body["database"])
	assert.Equal(t, "unhealthy", resp.Body["cache"])
}

func TestHealthServiceCheckIntegration(t *testing.T) {
	service, _, tc := setupTestServices(t)

	// Test service check method directly
	resp, err := service.Check(tc.DB.Statement.Context)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, StatusHealthy, resp.Database)
	assert.Equal(t, StatusHealthy, resp.Cache)
}

func TestHealthServiceCheckWithDatabaseErrorIntegration(t *testing.T) {
	// Clean cache before test
	sharedContainer.CleanupCache(t)

	// Create service with nil database
	store := NewStore(nil, sharedContainer.Redis)
	service := NewService(store)

	resp, err := service.Check(context.Background())
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, StatusUnhealthy, resp.Database)
	assert.Equal(t, StatusHealthy, resp.Cache)
}

func TestHealthServiceCheckWithCacheErrorIntegration(t *testing.T) {
	// Clean database before test
	sharedContainer.CleanupDatabase(t)

	// Create service with nil Redis
	store := NewStore(sharedContainer.DB, nil)
	service := NewService(store)

	resp, err := service.Check(context.Background())
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, StatusHealthy, resp.Database)
	assert.Equal(t, StatusUnhealthy, resp.Cache)
}

func TestHealthStorePingDatabaseIntegration(t *testing.T) {
	// Clean database before test
	sharedContainer.CleanupDatabase(t)

	store := NewStore(sharedContainer.DB, nil)

	err := store.PingDatabase(context.Background())
	assert.NoError(t, err)
}

func TestHealthStorePingCacheIntegration(t *testing.T) {
	// Clean cache before test
	sharedContainer.CleanupCache(t)

	store := NewStore(nil, sharedContainer.Redis)

	err := store.PingCache(context.Background())
	assert.NoError(t, err)
}

func TestHealthStorePingDatabaseErrorIntegration(t *testing.T) {
	store := NewStore(nil, nil)

	err := store.PingDatabase(context.Background())
	assert.Error(t, err)
}

func TestHealthStorePingCacheErrorIntegration(t *testing.T) {
	store := NewStore(nil, nil)

	err := store.PingCache(context.Background())
	assert.Error(t, err)
}

func TestHealthEndpointResponseFormatIntegration(t *testing.T) {
	_, handler, _ := setupTestServices(t)

	resp := test.MakeJSONRequest(t, handler.Health, test.HTTPRequest{
		Method: http.MethodGet,
		URL:    "/health",
	})

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, resp.Body)

	// Check that response has the expected structure
	_, hasDatabaseField := resp.Body["database"]
	_, hasCacheField := resp.Body["cache"]

	assert.True(t, hasDatabaseField, "response should have 'database' field")
	assert.True(t, hasCacheField, "response should have 'cache' field")

	// Check that values are valid status strings
	assert.Contains(t, []string{"healthy", "unhealthy"}, resp.Body["database"])
	assert.Contains(t, []string{"healthy", "unhealthy"}, resp.Body["cache"])
}
