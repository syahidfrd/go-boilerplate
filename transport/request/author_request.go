package request

import validation "github.com/go-ozzo/ozzo-validation"

type CreateAuthorReq struct {
	Name string `json:"name"`
}

func (request CreateAuthorReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
	)
}

type UpdateAuthorReq struct {
	Name string `json:"name"`
}

func (request UpdateAuthorReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
	)
}
