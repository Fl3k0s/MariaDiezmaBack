package service_test

import (
	"context"
	"testing"
	"time"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/repository/memory"
	"mariadiezmaback/internal/service"
)

func TestDressService_ListAndEnsureDefaults(t *testing.T) {
	repo := memory.NewDressRepository()
	svc := service.NewDressService(repo)
	ctx := context.Background()

	// 1. Initial list when empty
	initialList, err := svc.List(ctx, "")
	if err != nil {
		t.Fatalf("expected no error listing empty dresses, got %v", err)
	}
	if len(initialList) != 0 {
		t.Fatalf("expected 0 dresses, got %d", len(initialList))
	}

	// 2. Ensure defaults
	err = svc.EnsureDefaultDresses(ctx)
	if err != nil {
		t.Fatalf("expected no error seeding default dresses, got %v", err)
	}

	// 3. List all dresses
	list, err := svc.List(ctx, "")
	if err != nil {
		t.Fatalf("expected no error listing seeded dresses, got %v", err)
	}
	if len(list) != 10 {
		t.Fatalf("expected 10 seeded dresses, got %d", len(list))
	}

	for _, d := range list {
		if d.Name == "" {
			t.Errorf("expected dress to have a name")
		}
		if d.Collection == "" {
			t.Errorf("expected dress to have a collection")
		}
		if d.ImagePath == "" {
			t.Errorf("expected dress to have an image path")
		}
	}

	// 4. Filter by collection
	filtered, err := svc.List(ctx, "Esencia Floral")
	if err != nil {
		t.Fatalf("failed to filter dresses by collection: %v", err)
	}
	if len(filtered) != 5 {
		t.Fatalf("expected 5 dresses for 'Esencia Floral', got %d", len(filtered))
	}

	// 5. Add custom dress
	custom := &domain.Dress{
		ID:         "custom-dress-1",
		Name:       "Vestido Nupcial Sirena",
		Collection: "Colección Novias",
		ImagePath:  "assets/images/novias/MARIA_DIEZMA_031.jpg",
		Image1Path: "assets/images/novias/MARIA_DIEZMA_031.jpg",
		Image2Path: "assets/images/novias/MARIA_DIEZMA_032.jpg",
		Image3Path: "assets/images/novias/MARIA_DIEZMA_033.jpg",
		CreatedAt:  time.Now().UTC().Add(1 * time.Hour),
		UpdatedAt:  time.Now().UTC().Add(1 * time.Hour),
	}
	if err := repo.Create(ctx, custom); err != nil {
		t.Fatalf("failed to create custom dress: %v", err)
	}

	listAfterCustom, err := svc.List(ctx, "")
	if err != nil {
		t.Fatalf("failed to list after adding custom dress: %v", err)
	}
	if len(listAfterCustom) != 11 {
		t.Fatalf("expected 11 dresses, got %d", len(listAfterCustom))
	}
	if listAfterCustom[0].Name != "Vestido Nupcial Sirena" {
		t.Errorf("expected newest dress 'Vestido Nupcial Sirena' to be first, got '%s'", listAfterCustom[0].Name)
	}
}

func TestDressService_GetByNameAndCollection(t *testing.T) {
	repo := memory.NewDressRepository()
	svc := service.NewDressService(repo)
	ctx := context.Background()

	_ = svc.EnsureDefaultDresses(ctx)

	// 1. Success case: find "Vestido Magnolia" in "Esencia Floral"
	detail, err := svc.GetByNameAndCollection(ctx, "Vestido Magnolia", "Esencia Floral")
	if err != nil {
		t.Fatalf("expected to find dress, got error: %v", err)
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

	// 2. Case insensitive match
	detailCI, err := svc.GetByNameAndCollection(ctx, "vestido magnolia", "esencia floral")
	if err != nil {
		t.Fatalf("expected case-insensitive match, got error: %v", err)
	}
	if detailCI.Name != "Vestido Magnolia" {
		t.Errorf("expected name 'Vestido Magnolia', got '%s'", detailCI.Name)
	}

	// 3. Not found
	_, err = svc.GetByNameAndCollection(ctx, "Vestido Inexistente", "Esencia Floral")
	if err == nil {
		t.Fatalf("expected error for non-existent dress, got nil")
	}

	// 4. Missing parameters
	_, err = svc.GetByNameAndCollection(ctx, "", "Esencia Floral")
	if err == nil {
		t.Fatalf("expected validation error for empty name, got nil")
	}
}

func TestDressService_Create(t *testing.T) {
	repo := memory.NewDressRepository()
	svc := service.NewDressService(repo)
	ctx := context.Background()

	// 1. Success
	resp, err := svc.Create(ctx, domain.CreateDressInput{
		Name:        "Vestido Violeta",
		Collection:  "Esencia Floral",
		ImagePath:   "assets/images/esencia-floral/MARIA_DIEZMA_110.jpg",
		Description: "Vestido violeta artesanal.",
	})
	if err != nil {
		t.Fatalf("unexpected error creating dress: %v", err)
	}
	if resp.ID == "" {
		t.Errorf("expected non-empty ID")
	}
	if resp.Name != "Vestido Violeta" {
		t.Errorf("expected name 'Vestido Violeta', got '%s'", resp.Name)
	}
	if resp.Collection != "Esencia Floral" {
		t.Errorf("expected collection 'Esencia Floral', got '%s'", resp.Collection)
	}
	if resp.Image1Path != "assets/images/esencia-floral/MARIA_DIEZMA_110.jpg" {
		t.Errorf("expected image1 path fallback to ImagePath, got '%s'", resp.Image1Path)
	}

	// 2. Validation error - missing collection
	_, err = svc.Create(ctx, domain.CreateDressInput{
		Name: "Vestido Sin Colección",
	})
	if err == nil {
		t.Fatalf("expected validation error for missing collection, got nil")
	}

	// 3. Validation error - missing name
	_, err = svc.Create(ctx, domain.CreateDressInput{
		Collection: "Esencia Floral",
	})
	if err == nil {
		t.Fatalf("expected validation error for missing name, got nil")
	}
}


