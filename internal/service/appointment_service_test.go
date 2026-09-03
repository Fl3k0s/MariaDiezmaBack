package service_test

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/repository/memory"
	"mariadiezmaback/internal/service"
)

type MockMailer struct {
	mu            sync.Mutex
	sentCalls     []*domain.Appointment
	sentRecipient []string
}

func (m *MockMailer) SendAppointmentNotification(ctx context.Context, toEmail string, appt *domain.Appointment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sentCalls = append(m.sentCalls, appt)
	m.sentRecipient = append(m.sentRecipient, toEmail)
	return nil
}

func TestAppointmentService_Create(t *testing.T) {
	reqRepo := memory.NewRequestRepository()
	mockMail := &MockMailer{}
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	notificationEmail := "citas@mariadiezma.com"
	svc := service.NewAppointmentService(reqRepo, mockMail, notificationEmail, discardLogger)
	ctx := context.Background()

	// 1. Successful appointment creation
	input := domain.CreateAppointmentInput{
		Name:     "Laura Martínez",
		Email:    "laura@example.com",
		Phone:    "+34 654 321 987",
		Date:     "2026-10-15",
		TimeSlot: "16:00 - 17:00",
		Type:     "Consulta Nutricional Presencial",
	}

	appt, err := svc.Create(ctx, input)
	if err != nil {
		t.Fatalf("expected successful appointment creation, got error: %v", err)
	}

	if appt.ID == "" {
		t.Errorf("expected generated UUID")
	}
	if appt.Name != input.Name {
		t.Errorf("expected name %s, got %s", input.Name, appt.Name)
	}
	if appt.Email != input.Email {
		t.Errorf("expected email %s, got %s", input.Email, appt.Email)
	}
	if appt.Phone != input.Phone {
		t.Errorf("expected phone %s, got %s", input.Phone, appt.Phone)
	}
	if appt.Date != input.Date {
		t.Errorf("expected date %s, got %s", input.Date, appt.Date)
	}
	if appt.TimeSlot != input.TimeSlot {
		t.Errorf("expected time slot %s, got %s", input.TimeSlot, appt.TimeSlot)
	}
	if appt.Type != input.Type {
		t.Errorf("expected type %s, got %s", input.Type, appt.Type)
	}

	// Verify saved in backoffice repo
	savedItem, err := reqRepo.GetByID(ctx, appt.ID)
	if err != nil {
		t.Fatalf("expected request item to be saved in repository: %v", err)
	}
	if savedItem.Type != "appointment" {
		t.Errorf("expected item type 'appointment', got %s", savedItem.Type)
	}

	// Give goroutine a moment to dispatch email
	time.Sleep(50 * time.Millisecond)

	mockMail.mu.Lock()
	if len(mockMail.sentCalls) != 1 {
		t.Errorf("expected 1 email dispatched, got %d", len(mockMail.sentCalls))
	} else {
		if mockMail.sentRecipient[0] != notificationEmail {
			t.Errorf("expected recipient %s, got %s", notificationEmail, mockMail.sentRecipient[0])
		}
		if mockMail.sentCalls[0].Name != input.Name {
			t.Errorf("expected appointment name %s in email, got %s", input.Name, mockMail.sentCalls[0].Name)
		}
	}
	mockMail.mu.Unlock()

	// 2. Validation failure: missing fields
	badInput := domain.CreateAppointmentInput{
		Name:  "Incomplete",
		Email: "not-an-email",
	}
	_, err = svc.Create(ctx, badInput)
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}
}
