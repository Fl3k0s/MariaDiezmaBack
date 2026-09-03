package service_test

import (
	"context"
	"testing"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/repository/memory"
	"mariadiezmaback/internal/service"
)

func TestAuthService_EnsureAdminAndLogin(t *testing.T) {
	repo := memory.NewUserRepository()
	authSvc := service.NewAuthService(repo, "test-secret-key-12345", 2)
	ctx := context.Background()

	adminEmail := "admin@test.com"
	adminPassword := "SecretPassword123"

	// 1. Ensure Admin user
	err := authSvc.EnsureAdminUser(ctx, adminEmail, adminPassword)
	if err != nil {
		t.Fatalf("unexpected error creating admin: %v", err)
	}

	// 2. Ensure idempotent
	err = authSvc.EnsureAdminUser(ctx, adminEmail, adminPassword)
	if err != nil {
		t.Fatalf("unexpected error on second call to EnsureAdminUser: %v", err)
	}

	// 3. Login with correct credentials
	resp, err := authSvc.Login(ctx, domain.LoginInput{
		Email:    adminEmail,
		Password: adminPassword,
	})
	if err != nil {
		t.Fatalf("expected successful login, got: %v", err)
	}
	if resp.Token == "" {
		t.Errorf("expected non-empty token")
	}
	if resp.User.Email != adminEmail {
		t.Errorf("expected email %s, got %s", adminEmail, resp.User.Email)
	}

	// 4. Validate token
	claims, err := authSvc.ValidateToken(resp.Token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}
	if claims.Email != adminEmail {
		t.Errorf("expected claim email %s, got %s", adminEmail, claims.Email)
	}
	if claims.Role != domain.RoleAdmin {
		t.Errorf("expected admin role, got %s", claims.Role)
	}

	// 5. Login with wrong password
	_, err = authSvc.Login(ctx, domain.LoginInput{
		Email:    adminEmail,
		Password: "wrongpassword",
	})
	if err == nil {
		t.Fatalf("expected error on wrong password, got nil")
	}
}
