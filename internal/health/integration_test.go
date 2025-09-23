//go:build integration

package health

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/test"
)

func setupTestServices(t *testing.T) (*Service, *Handler, *test.Container) {
	t.Helper()

	tc := test.SetupFullContainer(t)
	tc.MigrateAll(t)

	store := NewStore(tc.DB, tc.Redis)
	service := NewService(store)
	handler := NewHandler(service)

	return service, handler, tc
}

func setupTestServicesWithDBOnly(t *testing.T) (*Service, *Handler, *test.Container) {
	t.Helper()

	tc := test.SetupPostgresContainer(t)
	tc.MigrateAll(t)

	// Create a nil Redis client to simulate Redis being down
	store := NewStore(tc.DB, nil)
	service := NewService(store)
	handler := NewHandler(service)

	return service, handler, tc
}

func setupTestServicesWithRedisOnly(t *testing.T) (*Service, *Handler, *test.Container) {
	t.Helper()

	tc := test.SetupRedisContainer(t)

	// Create a nil DB connection to simulate DB being down
	store := NewStore(nil, tc.Redis)
	service := NewService(store)
	handler := NewHandler(service)

	return service, handler, tc
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
	tc := test.SetupRedisContainer(t)

	// Create store with nil database to simulate database being down
	store := NewStore(nil, tc.Redis)
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
	tc := test.SetupPostgresContainer(t)
	tc.MigrateAll(t)

	// Create store with nil Redis to simulate cache being down
	store := NewStore(tc.DB, nil)
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
	tc := test.SetupRedisContainer(t)

	// Create service with nil database
	store := NewStore(nil, tc.Redis)
	service := NewService(store)

	resp, err := service.Check(tc.Redis.Context())
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, StatusUnhealthy, resp.Database)
	assert.Equal(t, StatusHealthy, resp.Cache)
}

func TestHealthServiceCheckWithCacheErrorIntegration(t *testing.T) {
	tc := test.SetupPostgresContainer(t)
	tc.MigrateAll(t)

	// Create service with nil Redis
	store := NewStore(tc.DB, nil)
	service := NewService(store)

	resp, err := service.Check(tc.DB.Statement.Context)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, StatusHealthy, resp.Database)
	assert.Equal(t, StatusUnhealthy, resp.Cache)
}

func TestHealthStorePingDatabaseIntegration(t *testing.T) {
	tc := test.SetupPostgresContainer(t)
	tc.MigrateAll(t)

	store := NewStore(tc.DB, nil)

	err := store.PingDatabase(tc.DB.Statement.Context)
	assert.NoError(t, err)
}

func TestHealthStorePingCacheIntegration(t *testing.T) {
	tc := test.SetupRedisContainer(t)

	store := NewStore(nil, tc.Redis)

	err := store.PingCache(tc.Redis.Context())
	assert.NoError(t, err)
}

func TestHealthStorePingDatabaseErrorIntegration(t *testing.T) {
	store := NewStore(nil, nil)

	err := store.PingDatabase(nil)
	assert.Error(t, err)
}

func TestHealthStorePingCacheErrorIntegration(t *testing.T) {
	store := NewStore(nil, nil)

	err := store.PingCache(nil)
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