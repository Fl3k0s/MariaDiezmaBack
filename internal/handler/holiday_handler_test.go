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

func setupHolidayTestRouter() (*chi.Mux, *service.HolidayService) {
	repo := memory.NewHolidayRepository()
	svc := service.NewHolidayService(repo)
	_ = svc.EnsureDefaultHolidays(context.Background())

	h := handler.NewHolidayHandler(svc)

	r := chi.NewRouter()
	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/holidays", h.List)
		api.Get("/festivos", h.List)
		api.Get("/holidays/{year}", h.GetByYear)
		api.Get("/festivos/{year}", h.GetByYear)
		api.Get("/holidays/id/{id}", h.GetByID)
		api.Get("/festivos/id/{id}", h.GetByID)

		api.Post("/holidays", h.Create)
		api.Post("/festivos", h.Create)
		api.Put("/holidays/{id}", h.Update)
		api.Put("/festivos/{id}", h.Update)
		api.Delete("/holidays/{id}", h.Delete)
		api.Delete("/festivos/{id}", h.Delete)
		api.Delete("/holidays/year/{year}", h.DeleteByYear)
		api.Delete("/festivos/year/{year}", h.DeleteByYear)
	})

	return r, svc
}

func TestHolidayHandler_ListWithResults(t *testing.T) {
	router, _ := setupHolidayTestRouter()

	// 1. GET /api/v1/festivos?year=2026
	req := httptest.NewRequest(http.MethodGet, "/api/v1/festivos?year=2026", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}

	var res response.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed unmarshaling response: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success true")
	}

	dataBytes, _ := json.Marshal(res.Data)
	var holidays []domain.HolidayResponse
	if err := json.Unmarshal(dataBytes, &holidays); err != nil {
		t.Fatalf("failed unmarshaling data: %v", err)
	}

	if len(holidays) != 14 {
		t.Fatalf("expected 14 holidays in 2026, got %d", len(holidays))
	}

	// Verify both date/fecha and name/nombre fields are present
	first := holidays[0]
	if first.Date == "" || first.Fecha == "" || first.Name == "" || first.Nombre == "" {
		t.Errorf("expected bilingual fields populated in holiday response: %+v", first)
	}

	// 2. GET /api/v1/holidays/2026 (Path param)
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/holidays/2026", nil)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec2.Code)
	}
}

func TestHolidayHandler_ListWithoutResults_ReturnsEmptyArrayForNormalCalendarDisplay(t *testing.T) {
	router, _ := setupHolidayTestRouter()

	// Consultation for a year that has NO defined holidays (e.g. 2035)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/festivos?year=2035", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// MUST be 200 OK so web doesn't crash
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK even without results, got %d", rec.Code)
	}

	var res response.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed decoding json: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success: true")
	}

	// Verify that data is serialized as "[]" and not "null"
	var rawMap map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &rawMap)
	dataSlice, ok := rawMap["data"].([]any)
	if !ok {
		t.Fatalf("expected data to be an array, got %T: %v", rawMap["data"], rawMap["data"])
	}
	if len(dataSlice) != 0 {
		t.Fatalf("expected empty array [], got %d items", len(dataSlice))
	}
}

func TestHolidayHandler_CreateSingleAndBatch(t *testing.T) {
	router, _ := setupHolidayTestRouter()

	// 1. Create single holiday with Spanish keys
	payloadSingle := []byte(`{
		"fecha": "2027-05-15",
		"nombre": "San Isidro",
		"descripcion": "Fiesta patronal de Madrid",
		"ambito": "local"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/festivos", bytes.NewReader(payloadSingle))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	// 2. Create batch holidays for year 2029
	payloadBatch := []byte(`[
		{"date": "2029-01-01", "name": "Año Nuevo", "year": 2029},
		{"date": "2029-01-06", "name": "Reyes", "year": 2029}
	]`)

	reqBatch := httptest.NewRequest(http.MethodPost, "/api/v1/festivos", bytes.NewReader(payloadBatch))
	reqBatch.Header.Set("Content-Type", "application/json")
	recBatch := httptest.NewRecorder()
	router.ServeHTTP(recBatch, reqBatch)

	if recBatch.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for batch, got %d: %s", recBatch.Code, recBatch.Body.String())
	}

	// Check that 2029 now has 2 holidays
	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/festivos/2029", nil)
	recGet := httptest.NewRecorder()
	router.ServeHTTP(recGet, reqGet)

	var res response.APIResponse
	_ = json.Unmarshal(recGet.Body.Bytes(), &res)
	dataBytes, _ := json.Marshal(res.Data)
	var items []domain.HolidayResponse
	_ = json.Unmarshal(dataBytes, &items)
	if len(items) != 2 {
		t.Fatalf("expected 2 holidays in 2029, got %d", len(items))
	}
}

func TestHolidayHandler_DeleteAndYearReset(t *testing.T) {
	router, _ := setupHolidayTestRouter()

	// 1. Delete all holidays for 2026
	reqDel := httptest.NewRequest(http.MethodDelete, "/api/v1/festivos/year/2026", nil)
	recDel := httptest.NewRecorder()
	router.ServeHTTP(recDel, reqDel)

	if recDel.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on year deletion, got %d: %s", recDel.Code, recDel.Body.String())
	}

	// 2. Verify subsequent GET /api/v1/festivos/2026 returns 200 OK with empty array []
	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/festivos/2026", nil)
	recGet := httptest.NewRecorder()
	router.ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusOK {
		t.Fatalf("expected 200 OK after deletion, got %d", recGet.Code)
	}

	var res response.APIResponse
	_ = json.Unmarshal(recGet.Body.Bytes(), &res)
	dataBytes, _ := json.Marshal(res.Data)
	var items []domain.HolidayResponse
	_ = json.Unmarshal(dataBytes, &items)

	if len(items) != 0 {
		t.Fatalf("expected 0 holidays after deletion, got %d", len(items))
	}
}
