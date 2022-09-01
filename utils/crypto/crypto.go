package crypto

import (
	"context"
	"crypto/md5"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const HashingCost int = 4 // if under 4, it will return error, see golang.org/x/crypto/bcrypt/bcrypt.go:289

type cryptoService struct {
	bcryptHashingCost int
}

func NewCryptoService() CryptoService {
	return &cryptoService{
		bcryptHashingCost: HashingCost,
	}
}

// CreatePasswordHash creates a password hash of given `plainPassword`
func (s *cryptoService) CreatePasswordHash(ctx context.Context, plainPassword string) (hashedPassword string, err error) {
	passwordHashInBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), s.bcryptHashingCost)
	if err != nil {
		return
	}
	hashedPassword = string(passwordHashInBytes)
	return
}

// ValidatePassword validates given `hashedPassword` against `plainPassword`. It returns true if given passwords are matched.
func (s *cryptoService) ValidatePassword(ctx context.Context, hashedPassword, plainPassword string) (isValid bool) {
	hashedPasswordInBytes := []byte(hashedPassword)
	plainPasswordInBytes := []byte(plainPassword)
	err := bcrypt.CompareHashAndPassword(hashedPasswordInBytes, plainPasswordInBytes)
	isValid = err == nil
	return
}

// CreateMD5Hash returns md5 hash value of `plainText`
func (s *cryptoService) CreateMD5Hash(ctx context.Context, plainText string) (hashedText string) {
	strInByte := []byte(plainText)
	resultInByte := md5.Sum(strInByte)
	hashedText = fmt.Sprintf("%x", resultInByte)
	return
}
