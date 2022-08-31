package request

import validation "github.com/go-ozzo/ozzo-validation"

// CreateTodoReq represent create todo request body
type CreateTodoReq struct {
	Name string `json:"name"`
}

func (request CreateTodoReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
	)
}

// UpdateTodoReq represent update todo request body
type UpdateTodoReq struct {
	Name string `json:"name"`
}

func (request UpdateTodoReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
	)
}
