package utils

import (
	"sort"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo"
)

// Error is the response that represents an error
type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

// NewError creates a new error response
func NewError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})

	switch v := err.(type) {
	case *echo.HTTPError:
		e.Errors["message"] = v.Message
	default:
		e.Errors["message"] = v.Error()
	}

	return e
}

type invalidField struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// NewValidationError creates a new error response representing a data validation error (HTTP 400)
func NewValidationError(errs validation.Errors) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})

	var details []invalidField
	var fields []string
	for field := range errs {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		details = append(details, invalidField{
			Field: field,
			Error: errs[field].Error(),
		})
	}

	e.Errors["message"] = "there is some problem with the data you submitted"
	e.Errors["details"] = details

	return e
}

// NewAccessForbiddenError creates a new error response representing an authorization failure (HTTP 403)
func NewAccessForbiddenError() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["message"] = "access forbidden"
	return e
}

// NewNotFoundError creates a new error response representing a resource-not-found error (HTTP 404)
func NewNotFoundError() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["message"] = "resource not found"
	return e
}
