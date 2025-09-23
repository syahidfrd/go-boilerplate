package render

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// HTTPError represents an error response structure
type HTTPError struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

// Empty represents an empty struct for responses with no data
type Empty struct{}

// JSON writes JSON response with the given status code and data
func JSON(w http.ResponseWriter, code int, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
}

// JSONFromError writes JSON error response with appropriate status code based on error type
func JSONFromError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	msg := "Something went wrong"

	var httpErr interface{ HTTPStatusCode() int }
	var validationErrs validator.ValidationErrors

	switch {
	case errors.As(err, &httpErr):
		code = httpErr.HTTPStatusCode()
		msg = err.Error()
	case errors.As(err, new(*json.UnmarshalTypeError)),
		errors.As(err, new(*json.SyntaxError)),
		errors.As(err, &validationErrs),
		errors.Is(err, io.EOF),
		errors.Is(err, io.ErrUnexpectedEOF),
		errors.Is(err, strconv.ErrSyntax),
		errors.Is(err, strconv.ErrRange):
		code = http.StatusBadRequest
		msg = err.Error()
	}

	errResp := HTTPError{
		Message:   msg,
		RequestID: w.Header().Get("X-Request-Id"),
	}

	resp, _ := json.Marshal(errResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
