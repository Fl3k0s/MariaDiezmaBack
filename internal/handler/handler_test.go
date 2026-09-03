package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/handler"
	"mariadiezmaback/internal/middleware"
	"mariadiezmaback/internal/repository/memory"
	"mariadiezmaback/internal/service"
	"mariadiezmaback/pkg/response"
)

func setupTestServer(t *testing.T) (*chi.Mux, *service.AuthService, string) {
	userRepo := memory.NewUserRepository()
	reqRepo := memory.NewRequestRepository()

	jwtSecret := "test-secret-key-9876543210"
	authSvc := service.NewAuthService(userRepo, jwtSecret, 1)
	reqSvc := service.NewRequestService(reqRepo)

	ctx := context.Background()
	_ = authSvc.EnsureAdminUser(ctx, "admin@test.com", "Password123!")

	loginResp, err := authSvc.Login(ctx, domain.LoginInput{
		Email:    "admin@test.com",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("failed to login during test setup: %v", err)
	}

	r := chi.NewRouter()

	healthH := handler.NewHealthHandler()
	authH := handler.NewAuthHandler(authSvc)
	reqH := handler.NewRequestHandler(reqSvc)

	r.Get("/health", healthH.Health)

	r.Route("/api/v1", func(api chi.Router) {
		api.Post("/auth/login", authH.Login)
		api.Post("/requests", reqH.Create)

		api.Group(func(backoffice chi.Router) {
			backoffice.Use(middleware.Auth(authSvc))

			backoffice.Get("/auth/me", authH.Me)
			backoffice.Get("/requests", reqH.List)
			backoffice.Get("/requests/{id}", reqH.GetByID)
			backoffice.Patch("/requests/{id}/status", reqH.UpdateStatus)
			backoffice.Put("/requests/{id}", reqH.Update)
			backoffice.Delete("/requests/{id}", reqH.Delete)
		})
	})

	return r, authSvc, loginResp.Token
}

func TestHealthEndpoint(t *testing.T) {
	r, _, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res response.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !res.Success {
		t.Errorf("expected success true")
	}
}

func TestAPI_RequestFlow(t *testing.T) {
	r, _, token := setupTestServer(t)

	// 1. Unauthenticated request listing should fail with 401
	unauthReq := httptest.NewRequest(http.MethodGet, "/api/v1/requests", nil)
	unauthRec := httptest.NewRecorder()
	r.ServeHTTP(unauthRec, unauthReq)
	if unauthRec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", unauthRec.Code)
	}

	// 2. Submit a request publicly
	createPayload := map[string]any{
		"type":         "contact",
		"priority":     "high",
		"sender_name":  "Carlos Gomez",
		"sender_email": "carlos@example.com",
		"subject":      "Consulta de servicios",
		"message":      "Hola, me gustaría información detallada.",
	}
	body, _ := json.Marshal(createPayload)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/requests", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	r.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", createRec.Code, createRec.Body.String())
	}

	var createResp response.APIResponse
	_ = json.Unmarshal(createRec.Body.Bytes(), &createResp)
	reqMap := createResp.Data.(map[string]any)
	createdID := reqMap["id"].(string)

	// 3. List requests as authenticated backoffice user
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/requests", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listRec := httptest.NewRecorder()
	r.ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", listRec.Code, listRec.Body.String())
	}

	// 4. Update request status
	updatePayload := map[string]any{
		"status":        "in_progress",
		"internal_note": "En contacto con el cliente",
	}
	updateBody, _ := json.Marshal(updatePayload)
	patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/requests/"+createdID+"/status", bytes.NewReader(updateBody))
	patchReq.Header.Set("Authorization", "Bearer "+token)
	patchReq.Header.Set("Content-Type", "application/json")
	patchRec := httptest.NewRecorder()
	r.ServeHTTP(patchRec, patchReq)

	if patchRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", patchRec.Code, patchRec.Body.String())
	}
}
