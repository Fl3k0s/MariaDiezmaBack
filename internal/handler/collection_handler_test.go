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
	"mariadiezmaback/internal/repository/memory"
	"mariadiezmaback/internal/service"
	"mariadiezmaback/pkg/response"
)

func TestCollectionHandler_List(t *testing.T) {
	repo := memory.NewCollectionRepository()
	colSvc := service.NewCollectionService(repo)
	_ = colSvc.EnsureDefaultCollections(context.Background())

	colHandler := handler.NewCollectionHandler(colSvc)

	r := chi.NewRouter()
	r.Get("/api/v1/collections", colHandler.List)
	r.Get("/api/v1/colecciones", colHandler.List)

	// 1. Test GET /api/v1/collections
	req := httptest.NewRequest(http.MethodGet, "/api/v1/collections", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp response.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success true")
	}

	itemsBytes, _ := json.Marshal(resp.Data)
	var items []domain.CollectionResponse
	if err := json.Unmarshal(itemsBytes, &items); err != nil {
		t.Fatalf("failed to unmarshal collection items: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 collections, got %d", len(items))
	}

	first := items[0]
	if first.Name != "Atardecer Mediterráneo" {
		t.Errorf("expected newest collection 'Atardecer Mediterráneo', got '%s'", first.Name)
	}
	if first.ImagePath != "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_016.jpg" {
		t.Errorf("expected image path 'assets/images/atardecer-mediterraneo/MARIA_DIEZMA_016.jpg', got '%s'", first.ImagePath)
	}
	if first.Description == "" {
		t.Errorf("expected non-empty description")
	}

	second := items[1]
	if second.Name != "Esencia Floral" {
		t.Errorf("expected second collection 'Esencia Floral', got '%s'", second.Name)
	}

	// 2. Test GET /api/v1/colecciones (Spanish alias)
	reqES := httptest.NewRequest(http.MethodGet, "/api/v1/colecciones", nil)
	recES := httptest.NewRecorder()
	r.ServeHTTP(recES, reqES)

	if recES.Code != http.StatusOK {
		t.Fatalf("expected status 200 on /colecciones, got %d", recES.Code)
	}
}

func TestCollectionHandler_Create(t *testing.T) {
	repo := memory.NewCollectionRepository()
	colSvc := service.NewCollectionService(repo)
	colHandler := handler.NewCollectionHandler(colSvc)

	r := chi.NewRouter()
	r.Post("/api/v1/collections", colHandler.Create)
	r.Post("/api/v1/colecciones", colHandler.Create)

	// 1. Success case with Spanish keys
	payload := map[string]string{
		"nombre":      "Nueva Colección 2027",
		"imagen":      "assets/images/2027/portada.jpg",
		"descripcion": "Descripción de prueba para colección nueva.",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/collections", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp response.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Success {
		t.Errorf("expected success true")
	}

	// 2. Validation error when name is empty
	badPayload := map[string]string{
		"nombre": "",
		"imagen": "assets/images/2027/portada.jpg",
	}
	badBody, _ := json.Marshal(badPayload)
	reqBad := httptest.NewRequest(http.MethodPost, "/api/v1/colecciones", bytes.NewReader(badBody))
	reqBad.Header.Set("Content-Type", "application/json")
	recBad := httptest.NewRecorder()
	r.ServeHTTP(recBad, reqBad)

	if recBad.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for missing name, got %d", recBad.Code)
	}
}

