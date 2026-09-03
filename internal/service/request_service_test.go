package service_test

import (
	"context"
	"testing"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/repository/memory"
	"mariadiezmaback/internal/service"
)

func TestRequestService_Lifecycle(t *testing.T) {
	repo := memory.NewRequestRepository()
	svc := service.NewRequestService(repo)
	ctx := context.Background()

	// 1. Create a request
	req, err := svc.Create(ctx, domain.CreateRequestInput{
		Type:        "quote",
		Priority:    domain.PriorityHigh,
		SenderName:  "Maria Test",
		SenderEmail: "maria@example.com",
		SenderPhone: "+34 600 000 000",
		Subject:     "Presupuesto Web Backoffice",
		Message:     "Hola, necesito un presupuesto para un backoffice.",
		Metadata: map[string]any{
			"budget": 5000,
		},
	})
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	if req.ID == "" {
		t.Errorf("expected generated UUID, got empty")
	}
	if req.Status != domain.StatusPending {
		t.Errorf("expected status 'pending', got %s", req.Status)
	}
	if req.Priority != domain.PriorityHigh {
		t.Errorf("expected priority 'high', got %s", req.Priority)
	}

	// 2. Get by ID
	found, err := svc.GetByID(ctx, req.ID)
	if err != nil {
		t.Fatalf("failed to get request by ID: %v", err)
	}
	if found.SenderName != "Maria Test" {
		t.Errorf("expected sender name 'Maria Test', got %s", found.SenderName)
	}

	// 3. Update status
	updated, err := svc.UpdateStatus(ctx, req.ID, domain.UpdateStatusInput{
		Status:       domain.StatusInProgress,
		InternalNote: "Asignado para revisión técnica",
	})
	if err != nil {
		t.Fatalf("failed to update status: %v", err)
	}
	if updated.Status != domain.StatusInProgress {
		t.Errorf("expected status in_progress, got %s", updated.Status)
	}
	if updated.InternalNote != "Asignado para revisión técnica" {
		t.Errorf("expected internal note updated, got %s", updated.InternalNote)
	}

	// 4. List with filter
	list, total, err := svc.List(ctx, domain.RequestFilter{
		Status: domain.StatusInProgress,
	})
	if err != nil {
		t.Fatalf("failed to list requests: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Errorf("expected 1 item, got total=%d len=%d", total, len(list))
	}

	// 5. Delete request
	err = svc.Delete(ctx, req.ID)
	if err != nil {
		t.Fatalf("failed to delete request: %v", err)
	}

	_, err = svc.GetByID(ctx, req.ID)
	if err == nil {
		t.Fatalf("expected error when getting deleted request, got nil")
	}
}
