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

func TestAppointmentHandler_CreateWithSpanishKeys(t *testing.T) {
	reqRepo := memory.NewRequestRepository()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	apptSvc := service.NewAppointmentService(reqRepo, &mockTestMailer{}, "admin@test.com", logger)
	apptHandler := handler.NewAppointmentHandler(apptSvc)

	r := chi.NewRouter()
	r.Post("/api/v1/citas", apptHandler.Create)

	// Test sending Spanish JSON payload as requested by user
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
