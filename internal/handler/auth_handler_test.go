package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"mariadiezmaback/internal/handler"
	"mariadiezmaback/internal/repository/memory"
	"mariadiezmaback/internal/service"
	"mariadiezmaback/pkg/response"
)

func TestAuthHandler_Login(t *testing.T) {
	repo := memory.NewUserRepository()
	authSvc := service.NewAuthService(repo, "test-secret-jwt-key", 24)
	ctx := context.Background()

	adminEmail := "admin@mariadiezma.com"
	adminPassword := "AdminPass123!"

	if err := authSvc.EnsureAdminUser(ctx, adminEmail, adminPassword); err != nil {
		t.Fatalf("failed to seed admin user: %v", err)
	}

	authHandler := handler.NewAuthHandler(authSvc)

	r := chi.NewRouter()
	r.Post("/api/v1/auth/login", authHandler.Login)

	// 1. Success case with username = full email
	bodySuccessUserEmail, _ := json.Marshal(map[string]string{
		"username": adminEmail,
		"password": adminPassword,
	})
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bodySuccessUserEmail))
	req1.Header.Set("Content-Type", "application/json")
	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec1.Code, rec1.Body.String())
	}

	var resp1 response.APIResponse
	if err := json.Unmarshal(rec1.Body.Bytes(), &resp1); err != nil {
		t.Fatalf("failed to parse success response: %v", err)
	}
	if !resp1.Success {
		t.Errorf("expected success true, got false")
	}

	// 2. Success case with username = "admin" (short username)
	bodySuccessShortUser, _ := json.Marshal(map[string]string{
		"username": "admin",
		"password": adminPassword,
	})
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bodySuccessShortUser))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected status 200 with short username, got %d: %s", rec2.Code, rec2.Body.String())
	}

	// 3. Backward compatibility: success with email field
	bodySuccessEmail, _ := json.Marshal(map[string]string{
		"email":    adminEmail,
		"password": adminPassword,
	})
	req3 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bodySuccessEmail))
	req3.Header.Set("Content-Type", "application/json")
	rec3 := httptest.NewRecorder()
	r.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusOK {
		t.Fatalf("expected status 200 with email fallback, got %d: %s", rec3.Code, rec3.Body.String())
	}

	// 4. Invalid credentials case
	bodyInvalid, _ := json.Marshal(map[string]string{
		"username": "admin",
		"password": "WrongPassword!",
	})
	reqInvalid := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bodyInvalid))
	reqInvalid.Header.Set("Content-Type", "application/json")
	recInvalid := httptest.NewRecorder()
	r.ServeHTTP(recInvalid, reqInvalid)

	if recInvalid.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", recInvalid.Code, recInvalid.Body.String())
	}

	var respInvalid response.APIResponse
	_ = json.Unmarshal(recInvalid.Body.Bytes(), &respInvalid)
	if respInvalid.Success {
		t.Errorf("expected success false")
	}
	if respInvalid.Message != "invalid username or password" {
		t.Errorf("expected message 'invalid username or password', got '%s'", respInvalid.Message)
	}

	// 5. Missing credentials
	bodyEmpty, _ := json.Marshal(map[string]string{
		"username": "",
	})
	reqEmpty := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bodyEmpty))
	reqEmpty.Header.Set("Content-Type", "application/json")
	recEmpty := httptest.NewRecorder()
	r.ServeHTTP(recEmpty, reqEmpty)

	if recEmpty.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for empty body, got %d", recEmpty.Code)
	}
}
