package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// HTTPRequest represents an HTTP request for testing
type HTTPRequest struct {
	Method  string
	URL     string
	Body    any
	Headers map[string]string
}

// HTTPResponse represents an HTTP response for testing
type HTTPResponse struct {
	StatusCode int
	Body       map[string]any
	RawBody    []byte
}

// MakeJSONRequest creates and executes an HTTP request with JSON body
func MakeJSONRequest(t *testing.T, handler http.HandlerFunc, req HTTPRequest) *HTTPResponse {
	t.Helper()

	var body *bytes.Buffer
	if req.Body != nil {
		jsonBody, err := json.Marshal(req.Body)
		require.NoError(t, err)
		body = bytes.NewBuffer(jsonBody)
	} else {
		body = bytes.NewBuffer(nil)
	}

	httpReq := httptest.NewRequest(req.Method, req.URL, body)

	// Set default content type for JSON requests
	if req.Body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	// Set additional headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	handler(w, httpReq)

	response := &HTTPResponse{
		StatusCode: w.Code,
		RawBody:    w.Body.Bytes(),
	}

	// Try to parse JSON response body
	if len(response.RawBody) > 0 {
		var jsonBody map[string]any
		if err := json.Unmarshal(response.RawBody, &jsonBody); err == nil {
			response.Body = jsonBody
		}
	}

	return response
}

// MakeAuthenticatedRequest creates an HTTP request with Bearer token
func MakeAuthenticatedRequest(t *testing.T, handler http.HandlerFunc, req HTTPRequest, token string) *HTTPResponse {
	t.Helper()

	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}
	req.Headers["Authorization"] = "Bearer " + token

	return MakeJSONRequest(t, handler, req)
}

// AssertJSONResponse validates HTTP response status and JSON body
func AssertJSONResponse(t *testing.T, resp *HTTPResponse, expectedStatus int, expectedBody map[string]any) {
	t.Helper()

	require.Equal(t, expectedStatus, resp.StatusCode)

	if expectedBody != nil {
		require.NotNil(t, resp.Body, "expected JSON response body")
		for key, expectedValue := range expectedBody {
			require.Equal(t, expectedValue, resp.Body[key], "mismatch for key: %s", key)
		}
	}
}

// AssertErrorResponse validates error response with message
func AssertErrorResponse(t *testing.T, resp *HTTPResponse, expectedStatus int, expectedMessage string) {
	t.Helper()

	require.Equal(t, expectedStatus, resp.StatusCode)
	require.NotNil(t, resp.Body)
	require.Equal(t, expectedMessage, resp.Body["message"])
}
