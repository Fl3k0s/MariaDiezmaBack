package handler_test

import (
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
	if first.Name != "Esencia Floral" {
		t.Errorf("expected name 'Esencia Floral', got '%s'", first.Name)
	}
	if first.ImagePath != "assets/images/esencia-floral/MARIA_DIEZMA_001.jpg" {
		t.Errorf("expected image path 'assets/images/esencia-floral/MARIA_DIEZMA_001.jpg', got '%s'", first.ImagePath)
	}
	if first.Description == "" {
		t.Errorf("expected non-empty description")
	}

	// 2. Test GET /api/v1/colecciones (Spanish alias)
	reqES := httptest.NewRequest(http.MethodGet, "/api/v1/colecciones", nil)
	recES := httptest.NewRecorder()
	r.ServeHTTP(recES, reqES)

	if recES.Code != http.StatusOK {
		t.Fatalf("expected status 200 on /colecciones, got %d", recES.Code)
	}
}
