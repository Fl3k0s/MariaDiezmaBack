package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/repository"
)

type JWTClaims struct {
	UserID string          `json:"user_id"`
	Email  string          `json:"email"`
	Role   domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
	expHours  int
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, expHours int) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
		expHours:  expHours,
	}
}

func (s *AuthService) EnsureAdminUser(ctx context.Context, email, password string) error {
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		// Admin already exists. Verify if the stored hash matches the configured password.
		if err := bcrypt.CompareHashAndPassword([]byte(existing.PasswordHash), []byte(password)); err != nil {
			newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return fmt.Errorf("failed to hash password for admin update: %w", err)
			}
			if err := s.userRepo.UpdatePassword(ctx, existing.ID, string(newHash)); err != nil {
				return fmt.Errorf("failed to update admin password: %w", err)
			}
		}
		return nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now().UTC()
	admin := &domain.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: string(hash),
		Name:         "Administrator",
		Role:         domain.RoleAdmin,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return s.userRepo.Create(ctx, admin)
}

func (s *AuthService) Login(ctx context.Context, input domain.LoginInput) (*domain.AuthResponse, error) {
	identifier := input.GetIdentifier()
	user, err := s.userRepo.GetByUsernameOrEmail(ctx, identifier)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	expiresAt := time.Now().UTC().Add(time.Duration(s.expHours) * time.Hour)
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &domain.AuthResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		User:      user.ToResponse(),
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, domain.ErrUnauthorized
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, domain.ErrUnauthorized
	}

	return claims, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id string) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := user.ToResponse()
	return &resp, nil
}
