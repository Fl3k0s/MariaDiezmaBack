package service_test

import (
	"context"
	"testing"
	"time"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/repository/memory"
	"mariadiezmaback/internal/service"
)

func TestCollectionService_ListAndEnsureDefaults(t *testing.T) {
	repo := memory.NewCollectionRepository()
	svc := service.NewCollectionService(repo)
	ctx := context.Background()

	// 1. Initial list when empty
	initialList, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("expected no error listing empty collections, got %v", err)
	}
	if len(initialList) != 0 {
		t.Fatalf("expected 0 collections, got %d", len(initialList))
	}

	// 2. Ensure defaults
	err = svc.EnsureDefaultCollections(ctx)
	if err != nil {
		t.Fatalf("expected no error seeding defaults, got %v", err)
	}

	// 3. List should now contain defaults
	list, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("expected no error listing seeded collections, got %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 seeded collections, got %d", len(list))
	}

	for _, c := range list {
		if c.Name == "" {
			t.Errorf("expected collection to have a name")
		}
		if c.ImagePath == "" {
			t.Errorf("expected collection to have an image path")
		}
		if c.Description == "" {
			t.Errorf("expected collection to have a description")
		}
	}

	// 4. Adding custom collection
	custom := &domain.Collection{
		ID:          "custom-1",
		Name:        "Colección Exclusiva",
		ImagePath:   "assets/images/exclusiva/MARIA_DIEZMA_031.jpg",
		Description: "Edición limitada de piezas artesanales.",
		CreatedAt:   time.Now().UTC().Add(1 * time.Hour),
		UpdatedAt:   time.Now().UTC().Add(1 * time.Hour),
	}
	if err := repo.Create(ctx, custom); err != nil {
		t.Fatalf("failed to create custom collection: %v", err)
	}

	listAfterCustom, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("failed to list collections after adding custom: %v", err)
	}
	if len(listAfterCustom) != 3 {
		t.Fatalf("expected 3 collections, got %d", len(listAfterCustom))
	}
	if listAfterCustom[0].Name != "Colección Exclusiva" {
		t.Errorf("expected newest collection 'Colección Exclusiva' to be first, got '%s'", listAfterCustom[0].Name)
	}
}
