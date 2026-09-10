package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
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

type mockTestMailer struct{}

func (m *mockTestMailer) SendAppointmentNotification(ctx context.Context, toEmail string, appt *domain.Appointment) error {
	return nil
}

func setupTestRouter() (*chi.Mux, *memory.RequestRepo) {
	reqRepo := memory.NewRequestRepository()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	apptSvc := service.NewAppointmentService(reqRepo, &mockTestMailer{}, "admin@test.com", logger)
	apptHandler := handler.NewAppointmentHandler(apptSvc)

	r := chi.NewRouter()
	r.Post("/api/v1/citas", apptHandler.Create)
	r.Post("/api/v1/appointments", apptHandler.Create)
	return r, reqRepo
}

func TestAppointmentHandler_CreateWithSpanishKeys(t *testing.T) {
	r, _ := setupTestRouter()

	spanishPayload := map[string]any{
		"nombre":        "Ana Gomez",
		"email":         "ana@example.com",
		"telefono":      "+34 600 123 456",
		"fecha":         "2026-11-20",
		"tramo_horario": "10:00 - 11:00",
		"tipo_cita":     "Primera Consulta",
	}

	body, _ := json.Marshal(spanishPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/citas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var res response.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success true")
	}

	dataMap := res.Data.(map[string]any)
	if dataMap["name"] != "Ana Gomez" {
		t.Errorf("expected name 'Ana Gomez', got %v", dataMap["name"])
	}
	if dataMap["time_slot"] != "10:00 - 11:00" {
		t.Errorf("expected time_slot '10:00 - 11:00', got %v", dataMap["time_slot"])
	}
}

func TestAppointmentHandler_CreateWithWebFormFields(t *testing.T) {
	r, _ := setupTestRouter()

	// All 7 web fields as requested:
	// Tipo de cita (string), Fecha (string), Franja horaria (string),
	// Nombre y apellidos (string), Telefono de contacto (string),
	// Fecha estimada (date), Detalles (string)
	webPayload := map[string]any{
		"tipo_cita":         "Novia a medida",
		"fecha":             "2026-10-25",
		"franja_horaria":    "Tarde (16:00 - 19:00)",
		"nombre_apellidos":  "Lucía Domínguez",
		"telefono_contacto": "+34 678 901 234",
		"fecha_estimada":    "2027-05-15",
		"detalles":          "Interesada en telas de seda natural y encaje floral",
	}

	body, _ := json.Marshal(webPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/citas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var res response.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	dataMap := res.Data.(map[string]any)
	if dataMap["nombre_apellidos"] != "Lucía Domínguez" {
		t.Errorf("expected nombre_apellidos 'Lucía Domínguez', got %v", dataMap["nombre_apellidos"])
	}
	if dataMap["name"] != "Lucía Domínguez" {
		t.Errorf("expected name 'Lucía Domínguez', got %v", dataMap["name"])
	}
	if dataMap["telefono_contacto"] != "+34 678 901 234" {
		t.Errorf("expected telefono_contacto '+34 678 901 234', got %v", dataMap["telefono_contacto"])
	}
	if dataMap["franja_horaria"] != "Tarde (16:00 - 19:00)" {
		t.Errorf("expected franja_horaria 'Tarde (16:00 - 19:00)', got %v", dataMap["franja_horaria"])
	}
	if dataMap["tipo_cita"] != "Novia a medida" {
		t.Errorf("expected tipo_cita 'Novia a medida', got %v", dataMap["tipo_cita"])
	}
	if dataMap["fecha_estimada"] != "2027-05-15" {
		t.Errorf("expected fecha_estimada '2027-05-15', got %v", dataMap["fecha_estimada"])
	}
	if dataMap["estimated_date"] != "2027-05-15" {
		t.Errorf("expected estimated_date '2027-05-15', got %v", dataMap["estimated_date"])
	}
	if dataMap["detalles"] != "Interesada en telas de seda natural y encaje floral" {
		t.Errorf("expected detalles 'Interesada en telas de seda natural y encaje floral', got %v", dataMap["detalles"])
	}
}

func TestAppointmentHandler_CreateWithNullableFechaEstimada(t *testing.T) {
	r, _ := setupTestRouter()

	// fecha_estimada can come as null, and email is omitted
	payload := map[string]any{
		"tipo_cita":         "Madrina",
		"fecha":             "2026-11-10",
		"franja_horaria":    "Mañana (10:00 - 13:00)",
		"nombre_apellidos":  "Carmen Navarro",
		"telefono_contacto": "611223344",
		"fecha_estimada":    nil,
		"detalles":          "Sin mangas",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/citas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var res response.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	dataMap := res.Data.(map[string]any)
	if dataMap["fecha_estimada"] != nil {
		t.Errorf("expected fecha_estimada to be nil/null, got %v", dataMap["fecha_estimada"])
	}
	if dataMap["estimated_date"] != nil {
		t.Errorf("expected estimated_date to be nil/null, got %v", dataMap["estimated_date"])
	}
}

func TestAppointmentHandler_CreateWithFranaHorariaTypo(t *testing.T) {
	r, _ := setupTestRouter()

	// Specifically test the "frana_horaria" typo
	payload := map[string]any{
		"tipo_cita":         "Invitada",
		"fecha":             "2026-12-05",
		"frana_horaria":     "11:00 - 12:00",
		"nombre_apellidos":  "Elena Santos",
		"telefono_contacto": "655443322",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/citas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var res response.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	dataMap := res.Data.(map[string]any)
	if dataMap["franja_horaria"] != "11:00 - 12:00" {
		t.Errorf("expected franja_horaria '11:00 - 12:00', got %v", dataMap["franja_horaria"])
	}
}

func TestAppointmentHandler_ValidationFailure(t *testing.T) {
	r, _ := setupTestRouter()

	badPayload := map[string]any{
		"nombre_apellidos": "Solo nombre",
	}

	body, _ := json.Marshal(badPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/citas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d: %s", rec.Code, rec.Body.String())
	}
}
