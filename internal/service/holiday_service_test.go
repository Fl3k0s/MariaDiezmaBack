package service_test

import (
	"context"
	"errors"
	"testing"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/repository/memory"
	"mariadiezmaback/internal/service"
)

func TestHolidayService_ListEmptyReturnsEmptySlice(t *testing.T) {
	repo := memory.NewHolidayRepository()
	svc := service.NewHolidayService(repo)
	ctx := context.Background()

	// List when completely empty
	list, err := svc.List(ctx, 2026)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if list == nil {
		t.Fatalf("expected non-nil slice, got nil")
	}

	if len(list) != 0 {
		t.Fatalf("expected 0 items, got %d", len(list))
	}
}

func TestHolidayService_EnsureDefaultMadridHolidays(t *testing.T) {
	repo := memory.NewHolidayRepository()
	svc := service.NewHolidayService(repo)
	ctx := context.Background()

	// Ensure defaults
	if err := svc.EnsureDefaultHolidays(ctx); err != nil {
		t.Fatalf("failed seeding defaults: %v", err)
	}

	// List all
	all, err := svc.List(ctx, 0)
	if err != nil {
		t.Fatalf("failed listing all: %v", err)
	}
	if len(all) == 0 {
		t.Fatalf("expected holidays seeded, got 0")
	}

	// List for 2025
	h2025, err := svc.List(ctx, 2025)
	if err != nil {
		t.Fatalf("failed listing 2025: %v", err)
	}
	if len(h2025) != 15 {
		t.Errorf("expected 15 holidays in Madrid for 2025, got %d", len(h2025))
	}

	// List for 2026
	h2026, err := svc.List(ctx, 2026)
	if err != nil {
		t.Fatalf("failed listing 2026: %v", err)
	}
	if len(h2026) != 14 {
		t.Errorf("expected 14 holidays in Madrid for 2026, got %d", len(h2026))
	}

	// Verify San Isidro exists in 2026
	var foundSanIsidro bool
	for _, h := range h2026 {
		if h.Date == "2026-05-15" && h.Name == "San Isidro Labrador" {
			foundSanIsidro = true
			if h.Fecha != "2026-05-15" || h.Nombre != "San Isidro Labrador" {
				t.Errorf("expected Spanish aliases populated")
			}
		}
	}
	if !foundSanIsidro {
		t.Errorf("expected San Isidro on 2026-05-15")
	}

	// List for unknown year 2040 -> should return empty slice without error
	h2040, err := svc.List(ctx, 2040)
	if err != nil {
		t.Fatalf("expected no error for unconfigured year, got %v", err)
	}
	if len(h2040) != 0 {
		t.Fatalf("expected 0 holidays for 2040, got %d", len(h2040))
	}
}

func TestHolidayService_CreateAndValidate(t *testing.T) {
	repo := memory.NewHolidayRepository()
	svc := service.NewHolidayService(repo)
	ctx := context.Background()

	// 1. Missing date
	_, err := svc.Create(ctx, domain.CreateHolidayInput{
		Name: "Día Festivo",
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for missing date, got %v", err)
	}

	// 2. Invalid date format
	_, err = svc.Create(ctx, domain.CreateHolidayInput{
		Date: "15/05/2027",
		Name: "San Isidro",
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for non YYYY-MM-DD date, got %v", err)
	}

	// 3. Missing name
	_, err = svc.Create(ctx, domain.CreateHolidayInput{
		Date: "2027-05-15",
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for missing name, got %v", err)
	}

	// 4. Valid holiday creation (year automatically inferred from date)
	created, err := svc.Create(ctx, domain.CreateHolidayInput{
		Date:        "2027-05-15",
		Name:        "San Isidro Labrador",
		Description: "Patrón de Madrid",
		Scope:       "local",
	})
	if err != nil {
		t.Fatalf("expected successful creation, got %v", err)
	}
	if created.Year != 2027 {
		t.Errorf("expected inferred year 2027, got %d", created.Year)
	}
	if created.ID == "" {
		t.Errorf("expected generated UUID")
	}

	// 5. Get by ID
	fetched, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("failed to fetch by ID: %v", err)
	}
	if fetched.Name != "San Isidro Labrador" {
		t.Errorf("expected San Isidro Labrador, got %s", fetched.Name)
	}
}

func TestHolidayService_BatchCreate(t *testing.T) {
	repo := memory.NewHolidayRepository()
	svc := service.NewHolidayService(repo)
	ctx := context.Background()

	batch := []domain.CreateHolidayInput{
		{Date: "2028-01-01", Name: "Año Nuevo"},
		{Date: "2028-01-06", Name: "Reyes Magos"},
		{Date: "2028-05-01", Name: "Fiesta del Trabajo"},
	}

	created, err := svc.CreateBatch(ctx, batch)
	if err != nil {
		t.Fatalf("failed batch create: %v", err)
	}
	if len(created) != 3 {
		t.Fatalf("expected 3 items created, got %d", len(created))
	}

	list2028, err := svc.List(ctx, 2028)
	if err != nil {
		t.Fatalf("failed to list 2028: %v", err)
	}
	if len(list2028) != 3 {
		t.Fatalf("expected 3 items in 2028, got %d", len(list2028))
	}
}

func TestHolidayService_DeleteAndEliminationFlow(t *testing.T) {
	repo := memory.NewHolidayRepository()
	svc := service.NewHolidayService(repo)
	ctx := context.Background()

	_ = svc.EnsureDefaultHolidays(ctx)

	// 1. Delete single holiday
	list2026, _ := svc.List(ctx, 2026)
	targetID := list2026[0].ID
	initialCount := len(list2026)

	if err := svc.Delete(ctx, targetID); err != nil {
		t.Fatalf("failed to delete single holiday: %v", err)
	}

	afterDelete, _ := svc.List(ctx, 2026)
	if len(afterDelete) != initialCount-1 {
		t.Errorf("expected %d items after single deletion, got %d", initialCount-1, len(afterDelete))
	}

	// 2. Delete all holidays for year 2025
	countDeleted, err := svc.DeleteByYear(ctx, 2025)
	if err != nil {
		t.Fatalf("failed to delete by year: %v", err)
	}
	if countDeleted != 15 {
		t.Errorf("expected 15 holidays deleted for 2025, got %d", countDeleted)
	}

	// After deletion, listing 2025 MUST return an empty slice (not null, no error)
	list2025, err := svc.List(ctx, 2025)
	if err != nil {
		t.Fatalf("expected no error listing deleted year, got %v", err)
	}
	if list2025 == nil {
		t.Fatalf("expected non-nil empty slice")
	}
	if len(list2025) != 0 {
		t.Fatalf("expected 0 holidays for deleted 2025, got %d", len(list2025))
	}
}
