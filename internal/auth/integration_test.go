//go:build integration

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/jwt"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/test"
	"github.com/syahidfrd/go-boilerplate/internal/user"
)

func setupTestServices(t *testing.T) (*Service, *handler, *JWTMiddleware, *test.Container) {
	t.Helper()

	tc := test.SetupPostgresContainer(t)
	tc.MigrateAll(t)

	userStore := user.NewStore(tc.DB)
	userService := user.NewService(userStore)
	jwtService := jwt.NewService("test-secret-key-for-integration-tests")

	authService := NewService(userService, jwtService)
	authHandler := NewHandler(authService)
	jwtMiddleware := NewJWTMiddleware(jwtService)

	return authService, authHandler, jwtMiddleware, tc
}

func TestSignUpIntegration(t *testing.T) {
	_, handler, _, _ := setupTestServices(t)

	tests := []struct {
		name           string
		requestBody    SignUpRequest
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name: "successful signup",
			requestBody: SignUpRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "signup successfully",
			},
		},
		{
			name: "invalid email format",
			requestBody: SignUpRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "password too short",
			requestBody: SignUpRequest{
				Email:    "test2@example.com",
				Password: "123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing email",
			requestBody: SignUpRequest{
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			requestBody: SignUpRequest{
				Email: "test3@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := test.MakeJSONRequest(t, handler.SignUp, test.HTTPRequest{
				Method: http.MethodPost,
				URL:    "/signup",
				Body:   tt.requestBody,
			})

			if tt.expectedBody != nil {
				test.AssertJSONResponse(t, resp, tt.expectedStatus, tt.expectedBody)
			} else {
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestSignUpDuplicateEmailIntegration(t *testing.T) {
	_, handler, _, _ := setupTestServices(t)

	signupReq := SignUpRequest{
		Email:    "duplicate@example.com",
		Password: "password123",
	}

	// First signup should succeed
	resp1 := test.MakeJSONRequest(t, handler.SignUp, test.HTTPRequest{
		Method: http.MethodPost,
		URL:    "/signup",
		Body:   signupReq,
	})

	test.AssertJSONResponse(t, resp1, http.StatusOK, map[string]any{
		"message": "signup successfully",
	})

	// Second signup with same email should fail
	resp2 := test.MakeJSONRequest(t, handler.SignUp, test.HTTPRequest{
		Method: http.MethodPost,
		URL:    "/signup",
		Body:   signupReq,
	})

	test.AssertErrorResponse(t, resp2, http.StatusBadRequest, "user already exists")
}

func TestSignInIntegration(t *testing.T) {
	_, handler, _, _ := setupTestServices(t)

	// First, create a user
	signupReq := SignUpRequest{
		Email:    "signin@example.com",
		Password: "password123",
	}

	signupResp := test.MakeJSONRequest(t, handler.SignUp, test.HTTPRequest{
		Method: http.MethodPost,
		URL:    "/signup",
		Body:   signupReq,
	})
	require.Equal(t, http.StatusOK, signupResp.StatusCode)

	tests := []struct {
		name           string
		requestBody    SignInRequest
		expectedStatus int
		checkToken     bool
		expectedError  string
	}{
		{
			name: "successful signin",
			requestBody: SignInRequest{
				Email:    "signin@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			checkToken:     true,
		},
		{
			name: "wrong password",
			requestBody: SignInRequest{
				Email:    "signin@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid credentials",
		},
		{
			name: "non-existent user",
			requestBody: SignInRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid credentials",
		},
		{
			name: "invalid email format",
			requestBody: SignInRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			requestBody: SignInRequest{
				Email: "signin@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := test.MakeJSONRequest(t, handler.SignIn, test.HTTPRequest{
				Method: http.MethodPost,
				URL:    "/signin",
				Body:   tt.requestBody,
			})

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.checkToken {
				require.NotNil(t, resp.Body)
				assert.NotEmpty(t, resp.Body["access_token"])
			}

			if tt.expectedError != "" {
				test.AssertErrorResponse(t, resp, tt.expectedStatus, tt.expectedError)
			}
		})
	}
}

func TestJWTMiddlewareIntegration(t *testing.T) {
	authService, _, middleware, _ := setupTestServices(t)

	// Create a user and get a token
	ctx := context.Background()
	signupReq := &SignUpRequest{
		Email:    "middleware@example.com",
		Password: "password123",
	}

	_, err := authService.SignUp(ctx, signupReq)
	require.NoError(t, err)

	signinReq := &SignInRequest{
		Email:    "middleware@example.com",
		Password: "password123",
	}

	signinResp, err := authService.SignIn(ctx, signinReq)
	require.NoError(t, err)
	require.NotEmpty(t, signinResp.AccessToken)

	// Create a test handler that requires authentication
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "user id not found in context", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"user_id": userID,
			"message": "authenticated",
		})
	})

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		checkUserID    bool
		expectedError  string
	}{
		{
			name:           "valid token",
			token:          signinResp.AccessToken,
			expectedStatus: http.StatusOK,
			checkUserID:    true,
		},
		{
			name:           "missing authorization header",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "authorization header required",
		},
		{
			name:           "invalid token",
			token:          "invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *test.HTTPResponse

			if tt.token != "" {
				resp = test.MakeAuthenticatedRequest(t, func(w http.ResponseWriter, r *http.Request) {
					middleware.Authenticate(protectedHandler).ServeHTTP(w, r)
				}, test.HTTPRequest{
					Method: http.MethodGet,
					URL:    "/protected",
				}, tt.token)
			} else {
				resp = test.MakeJSONRequest(t, func(w http.ResponseWriter, r *http.Request) {
					middleware.Authenticate(protectedHandler).ServeHTTP(w, r)
				}, test.HTTPRequest{
					Method: http.MethodGet,
					URL:    "/protected",
				})
			}

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.checkUserID {
				require.NotNil(t, resp.Body)
				assert.Equal(t, "authenticated", resp.Body["message"])
				assert.NotNil(t, resp.Body["user_id"])
			}

			if tt.expectedError != "" {
				test.AssertErrorResponse(t, resp, tt.expectedStatus, tt.expectedError)
			}
		})
	}
}

func TestGetUserIDFromContextIntegration(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func() context.Context
		expectedUserID int64
		expectedOK     bool
	}{
		{
			name: "valid user ID in context",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), UserIDKey, int64(123))
			},
			expectedUserID: 123,
			expectedOK:     true,
		},
		{
			name: "no user ID in context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedUserID: 0,
			expectedOK:     false,
		},
		{
			name: "wrong type in context",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), UserIDKey, "123")
			},
			expectedUserID: 0,
			expectedOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			userID, ok := GetUserIDFromContext(ctx)

			assert.Equal(t, tt.expectedUserID, userID)
			assert.Equal(t, tt.expectedOK, ok)
		})
	}
}

func TestFullAuthFlowIntegration(t *testing.T) {
	_, handler, middleware, _ := setupTestServices(t)

	userEmail := "fullflow@example.com"
	userPassword := "password123"

	// Step 1: Sign up
	signupResp := test.MakeJSONRequest(t, handler.SignUp, test.HTTPRequest{
		Method: http.MethodPost,
		URL:    "/signup",
		Body: SignUpRequest{
			Email:    userEmail,
			Password: userPassword,
		},
	})

	test.AssertJSONResponse(t, signupResp, http.StatusOK, map[string]any{
		"message": "signup successfully",
	})

	// Step 2: Sign in
	signinResp := test.MakeJSONRequest(t, handler.SignIn, test.HTTPRequest{
		Method: http.MethodPost,
		URL:    "/signin",
		Body: SignInRequest{
			Email:    userEmail,
			Password: userPassword,
		},
	})

	require.Equal(t, http.StatusOK, signinResp.StatusCode)
	require.NotNil(t, signinResp.Body)
	token := signinResp.Body["access_token"].(string)
	require.NotEmpty(t, token)

	// Step 3: Access protected endpoint
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r.Context())
		require.True(t, ok)
		require.Greater(t, userID, int64(0))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"user_id": userID,
			"message": "access granted",
		})
	})

	protectedResp := test.MakeAuthenticatedRequest(t, func(w http.ResponseWriter, r *http.Request) {
		middleware.Authenticate(protectedHandler).ServeHTTP(w, r)
	}, test.HTTPRequest{
		Method: http.MethodGet,
		URL:    "/protected",
	}, token)

	test.AssertJSONResponse(t, protectedResp, http.StatusOK, map[string]any{
		"message": "access granted",
	})
	assert.NotNil(t, protectedResp.Body["user_id"])
}

func TestConcurrentAuthOperationsIntegration(t *testing.T) {
	_, handler, _, _ := setupTestServices(t)

	const numGoroutines = 10
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			resp := test.MakeJSONRequest(t, handler.SignUp, test.HTTPRequest{
				Method: http.MethodPost,
				URL:    "/signup",
				Body: SignUpRequest{
					Email:    fmt.Sprintf("concurrent%d@example.com", id),
					Password: "password123",
				},
			})

			if resp.StatusCode != http.StatusOK {
				results <- fmt.Errorf("expected status 200, got %d", resp.StatusCode)
				return
			}

			results <- nil
		}(i)
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		require.NoError(t, err)
	}
}
