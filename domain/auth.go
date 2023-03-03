package domain

import (
	"context"

	"github.com/syahidfrd/go-boilerplate/transport/request"
)

// AuthUsecase represent the auth usecase contract
type AuthUsecase interface {
	SignUp(ctx context.Context, request *request.SignUpReq) error
	SignIn(ctx context.Context, request *request.SignInReq) (string, error)
}
