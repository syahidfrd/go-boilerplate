package usecase

import (
	"context"
	"database/sql"
	"github.com/syahidfrd/go-boilerplate/entity"
	"github.com/syahidfrd/go-boilerplate/repository/pgsql"
	"github.com/syahidfrd/go-boilerplate/transport/request"
	"github.com/syahidfrd/go-boilerplate/utils"
	"github.com/syahidfrd/go-boilerplate/utils/crypto"
	"github.com/syahidfrd/go-boilerplate/utils/jwt"
	"time"
)

// AuthUsecase represent the todos usecase contract
type AuthUsecase interface {
	SignUp(ctx context.Context, request *request.SignUpReq) error
	SignIn(ctx context.Context, request *request.SignInReq) (string, error)
}

type authUsecase struct {
	userRepo       pgsql.UserRepository
	cryptoSvc      crypto.CryptoService
	jwtSvc         jwt.JWTService
	contextTimeout time.Duration
}

// NewAuthUsecase will create new an authUsecase object representation of AuthUsecase interface
func NewAuthUsecase(userRepo pgsql.UserRepository, cryptoSvc crypto.CryptoService, jwtSvc jwt.JWTService, contextTimeout time.Duration) AuthUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		cryptoSvc:      cryptoSvc,
		jwtSvc:         jwtSvc,
		contextTimeout: contextTimeout,
	}
}

func (u *authUsecase) SignUp(c context.Context, request *request.SignUpReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, request.Email)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if user.ID != 0 {
		err = utils.NewBadRequestError("email already registered")
		return
	}

	passwordHash, err := u.cryptoSvc.CreatePasswordHash(ctx, request.Password)
	if err != nil {
		return
	}

	err = u.userRepo.Create(ctx, &entity.User{
		Email:     request.Email,
		Password:  passwordHash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	return
}

func (u *authUsecase) SignIn(c context.Context, request *request.SignInReq) (accessToken string, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, request.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewBadRequestError("email and password not match")
			return
		}
		return
	}

	if !u.cryptoSvc.ValidatePassword(ctx, user.Password, request.Password) {
		err = utils.NewBadRequestError("email and password not match")
		return
	}

	accessToken, err = u.jwtSvc.GenerateToken(ctx, user.ID)
	return
}
