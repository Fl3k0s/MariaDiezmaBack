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

func TestDressHandler_List(t *testing.T) {
	repo := memory.NewDressRepository()
	dressSvc := service.NewDressService(repo)
	_ = dressSvc.EnsureDefaultDresses(context.Background())

	dressHandler := handler.NewDressHandler(dressSvc)

	r := chi.NewRouter()
	r.Get("/api/v1/dresses", dressHandler.List)
	r.Get("/api/v1/vestidos", dressHandler.List)

	// 1. Test GET /api/v1/vestidos
	req := httptest.NewRequest(http.MethodGet, "/api/v1/vestidos", nil)
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
	var items []domain.DressResponse
	if err := json.Unmarshal(itemsBytes, &items); err != nil {
		t.Fatalf("failed to unmarshal dress items: %v", err)
	}

	if len(items) != 10 {
		t.Fatalf("expected 10 dresses, got %d", len(items))
	}

	first := items[0]
	if first.Name != "Vestido Magnolia" {
		t.Errorf("expected name 'Vestido Magnolia', got '%s'", first.Name)
	}
	if first.Collection != "Esencia Floral" {
		t.Errorf("expected collection 'Esencia Floral', got '%s'", first.Collection)
	}
	if first.ImagePath != "assets/images/esencia-floral/MARIA_DIEZMA_001.jpg" {
		t.Errorf("expected image path 'assets/images/esencia-floral/MARIA_DIEZMA_001.jpg', got '%s'", first.ImagePath)
	}

	// 2. Test GET /api/v1/vestidos?coleccion=Esencia+Floral
	reqFilter := httptest.NewRequest(http.MethodGet, "/api/v1/vestidos?coleccion=Esencia+Floral", nil)
	recFilter := httptest.NewRecorder()
	r.ServeHTTP(recFilter, reqFilter)

	if recFilter.Code != http.StatusOK {
		t.Fatalf("expected status 200 on filtered dresses, got %d", recFilter.Code)
	}

	var respFilter response.APIResponse
	_ = json.Unmarshal(recFilter.Body.Bytes(), &respFilter)
	filterItemsBytes, _ := json.Marshal(respFilter.Data)
	var filteredItems []domain.DressResponse
	_ = json.Unmarshal(filterItemsBytes, &filteredItems)

	if len(filteredItems) != 5 {
		t.Fatalf("expected 5 filtered dresses, got %d", len(filteredItems))
	}
	for _, item := range filteredItems {
		if item.Collection != "Esencia Floral" {
			t.Errorf("expected collection 'Esencia Floral', got '%s'", item.Collection)
		}
	}
}

func TestDressHandler_GetDetail(t *testing.T) {
	repo := memory.NewDressRepository()
	dressSvc := service.NewDressService(repo)
	_ = dressSvc.EnsureDefaultDresses(context.Background())

	dressHandler := handler.NewDressHandler(dressSvc)

	r := chi.NewRouter()
	r.Get("/api/v1/vestidos/detalle", dressHandler.GetDetail)
	r.Get("/api/v1/dresses/detail", dressHandler.GetDetail)

	// 1. Success case: GET /api/v1/vestidos/detalle?nombre=Vestido+Magnolia&coleccion=Esencia+Floral
	req := httptest.NewRequest(http.MethodGet, "/api/v1/vestidos/detalle?nombre=Vestido+Magnolia&coleccion=Esencia+Floral", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp response.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	detailBytes, _ := json.Marshal(resp.Data)
	var detail domain.DressDetailResponse
	if err := json.Unmarshal(detailBytes, &detail); err != nil {
		t.Fatalf("failed to unmarshal detail: %v", err)
	}

	if detail.Name != "Vestido Magnolia" {
		t.Errorf("expected name 'Vestido Magnolia', got '%s'", detail.Name)
	}
	if detail.Collection != "Esencia Floral" {
		t.Errorf("expected collection 'Esencia Floral', got '%s'", detail.Collection)
	}
	if detail.Image1Path != "assets/images/esencia-floral/MARIA_DIEZMA_001.jpg" {
		t.Errorf("expected image 1 'assets/images/esencia-floral/MARIA_DIEZMA_001.jpg', got '%s'", detail.Image1Path)
	}
	if detail.Image2Path != "assets/images/esencia-floral/MARIA_DIEZMA_002.jpg" {
		t.Errorf("expected image 2 'assets/images/esencia-floral/MARIA_DIEZMA_002.jpg', got '%s'", detail.Image2Path)
	}
	if detail.Image3Path != "assets/images/esencia-floral/MARIA_DIEZMA_003.jpg" {
		t.Errorf("expected image 3 'assets/images/esencia-floral/MARIA_DIEZMA_003.jpg', got '%s'", detail.Image3Path)
	}
	if detail.Description == "" {
		t.Errorf("expected non-empty description")
	}

	// 2. Missing params: 400 Bad Request
	reqBad := httptest.NewRequest(http.MethodGet, "/api/v1/vestidos/detalle?nombre=Vestido+Magnolia", nil)
	recBad := httptest.NewRecorder()
	r.ServeHTTP(recBad, reqBad)

	if recBad.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 when collection is missing, got %d", recBad.Code)
	}

	// 3. Not found: 404 Not Found
	reqNF := httptest.NewRequest(http.MethodGet, "/api/v1/vestidos/detalle?nombre=Inexistente&coleccion=Esencia+Floral", nil)
	recNF := httptest.NewRecorder()
	r.ServeHTTP(recNF, reqNF)

	if recNF.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for non-existent dress, got %d", recNF.Code)
	}
}

