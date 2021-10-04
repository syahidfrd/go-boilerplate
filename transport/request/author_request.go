package request

import validation "github.com/go-ozzo/ozzo-validation"

// CreateAuthorReq represent create author request body
type CreateAuthorReq struct {
	Name string `json:"name"`
}

func (request CreateAuthorReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
	)
}

// UpdateAuthorReq represent update author request body
type UpdateAuthorReq struct {
	Name string `json:"name"`
}

func (request UpdateAuthorReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
	)
}
