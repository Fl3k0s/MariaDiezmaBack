package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/repository"
)

type CollectionService struct {
	repo repository.CollectionRepository
}

func NewCollectionService(repo repository.CollectionRepository) *CollectionService {
	return &CollectionService{repo: repo}
}

func (s *CollectionService) List(ctx context.Context) ([]domain.CollectionResponse, error) {
	collections, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collections: %w", err)
	}

	result := make([]domain.CollectionResponse, len(collections))
	for i, c := range collections {
		result[i] = c.ToResponse()
	}

	return result, nil
}

func (s *CollectionService) EnsureDefaultCollections(ctx context.Context) error {
	existing, err := s.repo.List(ctx)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	now := time.Now().UTC()
	defaults := []domain.Collection{
		{
			ID:          uuid.NewString(),
			Name:        "Esencia Floral",
			ImagePath:   "assets/images/esencia-floral/MARIA_DIEZMA_001.jpg",
			Description: "Diseños inspirados en la delicadeza botánica y tonalidades primaverales.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.NewString(),
			Name:        "Atardecer Mediterráneo",
			ImagePath:   "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_016.jpg",
			Description: "Colección cálida con texturas fluidas y tonos terracota y dorados.",
			CreatedAt:   now.Add(1 * time.Second),
			UpdatedAt:   now.Add(1 * time.Second),
		},
	}

	for _, c := range defaults {
		col := c
		if err := s.repo.Create(ctx, &col); err != nil {
			return fmt.Errorf("failed to seed collection %s: %w", c.Name, err)
		}
	}

	return nil
}
