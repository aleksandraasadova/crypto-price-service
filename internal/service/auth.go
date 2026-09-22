package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aleksandraasadova/crypto-price-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Get(ctx context.Context, username string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}

type AuthService struct {
	userRepository UserRepository
	jwtSecret      string
}

func NewAuthService(userRepository UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

func (a *AuthService) Register(ctx context.Context, username, password string) (string, error) {
	if err := validateCredentials(username, password); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	if _, err := a.userRepository.Get(ctx, username); err == nil {
		return "", domain.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := a.userRepository.Create(ctx, user); err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	token, err := a.generateToken(user.Username)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (a *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	if err := validateCredentials(username, password); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	user, err := a.userRepository.Get(ctx, username)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := a.generateToken(username)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (a *AuthService) generateToken(username string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()), // Когда создан
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func validateCredentials(username, password string) error {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return domain.ErrInvalidCredentials
	}

	if strings.Contains(username, " ") {
		return domain.ErrInvalidCredentials
	}

	return nil
}
