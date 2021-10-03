package utils

import (
	"sort"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo"
)

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

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

func NewAccessForbiddenError() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["message"] = "access forbidden"
	return e
}

func NewNotFoundError() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["message"] = "resource not found"
	return e
}
