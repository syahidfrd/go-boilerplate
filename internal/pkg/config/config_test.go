package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase_DataSourceName(t *testing.T) {
	testCases := []struct {
		name     string
		database Database
		expected string
	}{
		{
			name: "Complete database config",
			database: Database{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				Password: "testpass",
				Name:     "testdb",
			},
			expected: "user=testuser password=testpass host=localhost port=5432 dbname=testdb sslmode=disable",
		},
		{
			name: "Database config with different port",
			database: Database{
				Host:     "prod-db.example.com",
				Port:     5433,
				User:     "produser",
				Password: "prodpass",
				Name:     "proddb",
			},
			expected: "user=produser password=prodpass host=prod-db.example.com port=5433 dbname=proddb sslmode=disable",
		},
		{
			name: "Database config with empty values",
			database: Database{
				Host:     "",
				Port:     0,
				User:     "",
				Password: "",
				Name:     "",
			},
			expected: "user= password= host= port=0 dbname= sslmode=disable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.database.DataSourceName()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLoadEnv(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{
		"APP_SECRET", "CACHE_URL",
		"DATABASE_HOST", "DATABASE_PORT", "DATABASE_USER", "DATABASE_PASSWORD",
		"DATABASE_NAME", "DATABASE_MAX_IDLE_CONN", "DATABASE_MAX_OPEN_CONN",
	}

	for _, envVar := range envVars {
		originalEnv[envVar] = os.Getenv(envVar)
	}

	// Clean environment
	defer func() {
		for _, envVar := range envVars {
			if original, exists := originalEnv[envVar]; exists && original != "" {
				os.Setenv(envVar, original)
			} else {
				os.Unsetenv(envVar)
			}
		}
	}()

	// Set test environment variables
	testEnv := map[string]string{
		"APP_SECRET":             "test-jwt-secret",
		"CACHE_URL":              "localhost:6379",
		"DATABASE_HOST":          "localhost",
		"DATABASE_PORT":          "5432",
		"DATABASE_USER":          "testuser",
		"DATABASE_PASSWORD":      "testpass",
		"DATABASE_NAME":          "testdb",
		"DATABASE_MAX_IDLE_CONN": "5",
		"DATABASE_MAX_OPEN_CONN": "10",
	}

	for key, value := range testEnv {
		os.Setenv(key, value)
	}

	// Load configuration
	config := LoadEnv()

	// Verify configuration
	assert.NotNil(t, config)
	assert.Equal(t, "test-jwt-secret", config.AppSecret)
	assert.Equal(t, "localhost:6379", config.CacheURL)

	// Verify database configuration
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, 5432, config.Database.Port)
	assert.Equal(t, "testuser", config.Database.User)
	assert.Equal(t, "testpass", config.Database.Password)
	assert.Equal(t, "testdb", config.Database.Name)
	assert.Equal(t, 5, config.Database.MaxIdleConn)
	assert.Equal(t, 10, config.Database.MaxOpenConn)
}

func TestLoadEnv_WithDefaults(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{
		"APP_SECRET", "CACHE_URL",
		"DATABASE_HOST", "DATABASE_PORT", "DATABASE_USER", "DATABASE_PASSWORD",
		"DATABASE_NAME", "DATABASE_MAX_IDLE_CONN", "DATABASE_MAX_OPEN_CONN",
	}

	for _, envVar := range envVars {
		originalEnv[envVar] = os.Getenv(envVar)
	}

	// Clean environment to test defaults
	defer func() {
		for _, envVar := range envVars {
			if original, exists := originalEnv[envVar]; exists && original != "" {
				os.Setenv(envVar, original)
			} else {
				os.Unsetenv(envVar)
			}
		}
	}()

	// Clear all environment variables
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}

	// Load configuration
	config := LoadEnv()

	// Verify configuration with defaults (empty values)
	assert.NotNil(t, config)
	assert.Equal(t, "", config.AppSecret)
	assert.Equal(t, "", config.CacheURL)
	assert.Equal(t, 0, config.Database.MaxIdleConn)
	assert.Equal(t, 0, config.Database.MaxOpenConn)
}

func TestLoadEnv_IntegerParsing(t *testing.T) {
	// Save original environment
	originalEnv := map[string]string{
		"DATABASE_MAX_IDLE_CONN": os.Getenv("DATABASE_MAX_IDLE_CONN"),
		"DATABASE_MAX_OPEN_CONN": os.Getenv("DATABASE_MAX_OPEN_CONN"),
	}

	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	// Test valid integer values
	os.Setenv("DATABASE_MAX_IDLE_CONN", "15")
	os.Setenv("DATABASE_MAX_OPEN_CONN", "25")

	config := LoadEnv()

	assert.Equal(t, 15, config.Database.MaxIdleConn)
	assert.Equal(t, 25, config.Database.MaxOpenConn)
}

func TestConfig_Struct_Initialization(t *testing.T) {
	config := &Config{
		AppSecret: "my-secret",
		CacheURL:  "redis:6379",
		Database: Database{
			Host:        "db-host",
			Port:        5432,
			User:        "user",
			Password:    "pass",
			Name:        "dbname",
			MaxIdleConn: 10,
			MaxOpenConn: 20,
		},
	}

	assert.Equal(t, "my-secret", config.AppSecret)
	assert.Equal(t, "redis:6379", config.CacheURL)
	assert.Equal(t, "db-host", config.Database.Host)
	assert.Equal(t, 5432, config.Database.Port)
	assert.Equal(t, "user", config.Database.User)
	assert.Equal(t, "pass", config.Database.Password)
	assert.Equal(t, "dbname", config.Database.Name)
	assert.Equal(t, 10, config.Database.MaxIdleConn)
	assert.Equal(t, 20, config.Database.MaxOpenConn)
}

func TestDatabase_Struct_Initialization(t *testing.T) {
	db := Database{
		Host:        "localhost",
		Port:        5432,
		User:        "testuser",
		Password:    "testpass",
		Name:        "testdb",
		MaxIdleConn: 5,
		MaxOpenConn: 10,
	}

	assert.Equal(t, "localhost", db.Host)
	assert.Equal(t, 5432, db.Port)
	assert.Equal(t, "testuser", db.User)
	assert.Equal(t, "testpass", db.Password)
	assert.Equal(t, "testdb", db.Name)
	assert.Equal(t, 5, db.MaxIdleConn)
	assert.Equal(t, 10, db.MaxOpenConn)
}
