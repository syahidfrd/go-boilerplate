package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/syahidfrd/go-boilerplate/internal/jwt"
	"github.com/syahidfrd/go-boilerplate/internal/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidCredentials is returned when email/password combination is invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserAlreadyExists is returned when attempting to create a user with existing email
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrUserNotFound is returned when a requested user cannot be found
	ErrUserNotFound = errors.New("user not found")
)

// Service provides authentication business logic operations
type Service struct {
	userService UserService
	jwtService  JWTService
}

// UserService defines the interface for user data operations
type UserService interface {
	Create(ctx context.Context, email, hashedPassword string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
}

// JWTService defines the interface for JWT token operations
type JWTService interface {
	GenerateToken(userID int64) (string, error)
	ValidateToken(tokenString string) (*jwt.Claims, error)
}

// SignUpRequest represents the request payload for user registration
type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// SignUpResponse represents the response payload for successful user registration
type SignUpResponse struct {
	Message string `json:"message"`
}

// SignInRequest represents the request payload for user authentication
type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// SignInResponse represents the response payload for successful authentication
type SignInResponse struct {
	AccessToken string `json:"access_token"`
}

// NewService creates a new auth service with the provided dependencies
func NewService(userService UserService, jwtService JWTService) *Service {
	return &Service{
		userService: userService,
		jwtService:  jwtService,
	}
}

// hashPassword hashes the given password using bcrypt
func (s *Service) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// validatePassword validates the given password against the hashed password
func (s *Service) validatePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// SignUp handles user registration by validating input and creating a new user account
func (s *Service) SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error) {
	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user through user service
	_, err = s.userService.Create(ctx, req.Email, hashedPassword)
	if err != nil {
		// For now, assume any user creation error is due to duplicate email
		// You can add more specific error checking here based on user service errors
		return nil, ErrUserAlreadyExists
	}

	return &SignUpResponse{
		Message: "signup successfully",
	}, nil
}

// SignIn handles user authentication by validating credentials and generating JWT token
func (s *Service) SignIn(ctx context.Context, req *SignInRequest) (*SignInResponse, error) {
	// Get user by email
	u, err := s.userService.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// Validate password
	if err := s.validatePassword(u.Password, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := s.jwtService.GenerateToken(u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &SignInResponse{
		AccessToken: token,
	}, nil
}
