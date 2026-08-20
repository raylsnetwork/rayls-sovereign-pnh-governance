package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
)

const (
	// TokenExpiryDays is the number of days before a JWT token expires
	TokenExpiryDays = 30
)

var _ AuthService = (*authService)(nil)

// authService implements the AuthService interface
type authService struct {
	pnRepo    PrivateNetworkRepository
	jwtSecret []byte
	log       logger.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(pnRepo PrivateNetworkRepository, jwtSecret []byte, log logger.Logger) AuthService {
	return &authService{
		pnRepo:    pnRepo,
		jwtSecret: jwtSecret,
		log:       log,
	}
}

// SignUp creates a new private network account
func (s *authService) SignUp(ctx context.Context, username, password string) error {
	s.log.Info("Attempting to create new private network", "username", username)

	// Validate inputs
	if username == "" || password == "" {
		return NewValidationError("credentials", "Username and password cannot be empty")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("Failed to hash password", "error", err, "username", username)
		return NewInternalError("password hashing", err)
	}

	// Create private network
	err = s.pnRepo.Create(ctx, username, string(hashedPassword))
	if err != nil {
		if errors.Is(err, ErrRecordConflict) {
			s.log.Warn("Signup failed - username already exists", "username", username)
			return NewConflictError("username")
		}
		s.log.Error("Failed to create private network", "error", err, "username", username)
		return NewInternalError("private network creation", err)
	}

	s.log.Info("Successfully created private network", "username", username)
	return nil
}

// Login authenticates a private network and returns a JWT token
func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
	s.log.Info("Login attempt", "username", username)

	// Validate inputs
	if username == "" || password == "" {
		return "", NewValidationError("credentials", "Username and password cannot be empty")
	}

	// Find private network
	pn, err := s.pnRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			s.log.Warn("Login failed - user not found", "username", username)
			return "", NewValidationError("credentials", "Invalid username or password")
		}
		s.log.Error("Failed to retrieve private network", "error", err, "username", username)
		return "", NewInternalError("private network retrieval", err)
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(pn.Password), []byte(password))
	if err != nil {
		s.log.Warn("Login failed - invalid password", "username", username)
		return "", NewValidationError("credentials", "Invalid username or password")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": pn.Username,
		"exp": time.Now().Add(time.Hour * 24 * TokenExpiryDays).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		s.log.Error("Failed to generate JWT token", "error", err, "username", username)
		return "", NewInternalError("JWT generation", err)
	}

	s.log.Info("Login successful", "username", username)
	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the associated private network
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*domain.PrivateNetwork, error) {
	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		s.log.Warn("Token parsing failed", "error", err.Error())
		return nil, NewValidationError("token", "invalid token")
	}

	// Extract and validate claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		s.log.Warn("Invalid token claims")
		return nil, NewValidationError("token", "invalid token")
	}

	// Check expiration
	expFloat64, ok := claims["exp"].(float64)
	if !ok {
		s.log.Warn("Token missing expiration")
		return nil, NewValidationError("token", "invalid token")
	}
	if float64(time.Now().Unix()) > expFloat64 {
		s.log.Warn("Token expired")
		return nil, NewValidationError("token", "token expired")
	}

	// Extract username
	username, ok := claims["sub"].(string)
	if !ok || username == "" {
		s.log.Warn("Token missing subject")
		return nil, NewValidationError("token", "invalid token")
	}

	// Lookup private network
	pn, err := s.pnRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			s.log.Warn("Token validation failed - user not found", "username", username)
			return nil, NewValidationError("token", "invalid token")
		}
		s.log.Error("Failed to retrieve private network during token validation", "error", err, "username", username)
		return nil, NewInternalError("private network retrieval", err)
	}

	return pn, nil
}
