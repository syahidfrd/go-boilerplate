package render

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	tests := []struct {
		name         string
		code         int
		data         any
		expectedCode int
		expectedBody string
	}{
		{
			name:         "success with map",
			code:         http.StatusOK,
			data:         map[string]string{"message": "success"},
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"success"}`,
		},
		{
			name:         "success with struct",
			code:         http.StatusCreated,
			data:         struct{ ID int }{ID: 1},
			expectedCode: http.StatusCreated,
			expectedBody: `{"ID":1}`,
		},
		{
			name:         "success with nil",
			code:         http.StatusNoContent,
			data:         nil,
			expectedCode: http.StatusNoContent,
			expectedBody: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			JSON(w, tt.code, tt.data)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}

type mockHTTPError struct {
	message string
	code    int
}

func (e mockHTTPError) Error() string {
	return e.message
}

func (e mockHTTPError) HTTPStatusCode() int {
	return e.code
}

func TestJSONFromError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedMsg  string
		requestID    string
	}{
		{
			name:         "custom HTTP error",
			err:          mockHTTPError{message: "user not found", code: http.StatusNotFound},
			expectedCode: http.StatusNotFound,
			expectedMsg:  "user not found",
			requestID:    "req-123",
		},
		{
			name:         "json unmarshal error",
			err:          &json.UnmarshalTypeError{Value: "string", Type: reflect.TypeOf(0)},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "json: cannot unmarshal string into Go value of type int",
			requestID:    "",
		},
		{
			name:         "syntax error",
			err:          func() error { var v any; return json.Unmarshal([]byte("{invalid"), &v) }(),
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "invalid character 'i' looking for beginning of object key string",
			requestID:    "",
		},
		{
			name:         "strconv error",
			err:          strconv.ErrSyntax,
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "invalid syntax",
			requestID:    "",
		},
		{
			name:         "generic error",
			err:          errors.New("database connection failed"),
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Something went wrong",
			requestID:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if tt.requestID != "" {
				w.Header().Set("X-Request-Id", tt.requestID)
			}

			JSONFromError(w, tt.err)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var response HTTPError
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedMsg, response.Message)
			assert.Equal(t, tt.requestID, response.RequestID)
		})
	}
}

func TestJSONFromError_ValidationErrors(t *testing.T) {
	// Create a validator instance
	validate := validator.New()

	// Define a struct with validation tags
	type TestStruct struct {
		Email string `validate:"required,email"`
		Age   int    `validate:"gte=0,lte=130"`
	}

	// Create invalid data
	testData := TestStruct{
		Email: "invalid-email",
		Age:   -1,
	}

	// Validate and get validation errors
	err := validate.Struct(testData)
	assert.Error(t, err)

	w := httptest.NewRecorder()
	JSONFromError(w, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response HTTPError
	jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, jsonErr)

	// Validation error message should contain field validation details
	assert.Contains(t, response.Message, "Email")
	assert.Contains(t, response.Message, "Age")
}

func TestHTTPError(t *testing.T) {
	httpErr := HTTPError{
		Message:   "test error",
		RequestID: "req-456",
	}

	jsonData, err := json.Marshal(httpErr)
	assert.NoError(t, err)

	expected := `{"message":"test error","request_id":"req-456"}`
	assert.Equal(t, expected, string(jsonData))
}

func TestEmpty(t *testing.T) {
	empty := Empty{}

	jsonData, err := json.Marshal(empty)
	assert.NoError(t, err)

	expected := `{}`
	assert.Equal(t, expected, string(jsonData))
}
