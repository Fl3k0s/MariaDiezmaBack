package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/repository"
)

type DressService struct {
	repo repository.DressRepository
}

func NewDressService(repo repository.DressRepository) *DressService {
	return &DressService{repo: repo}
}

func (s *DressService) List(ctx context.Context, collectionFilter string) ([]domain.DressResponse, error) {
	dresses, err := s.repo.List(ctx, collectionFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dresses: %w", err)
	}

	result := make([]domain.DressResponse, len(dresses))
	for i, d := range dresses {
		result[i] = d.ToResponse()
	}

	return result, nil
}

func (s *DressService) GetByNameAndCollection(ctx context.Context, name, collection string) (*domain.DressDetailResponse, error) {
	cleanName := strings.TrimSpace(name)
	cleanCol := strings.TrimSpace(collection)

	if cleanName == "" || cleanCol == "" {
		return nil, fmt.Errorf("%w: nombre y colección son campos requeridos", domain.ErrInvalidInput)
	}

	dress, err := s.repo.GetByNameAndCollection(ctx, cleanName, cleanCol)
	if err != nil {
		return nil, err
	}

	detail := dress.ToDetailResponse()
	return &detail, nil
}

func (s *DressService) Create(ctx context.Context, input domain.CreateDressInput) (*domain.DressDetailResponse, error) {
	cleanName := strings.TrimSpace(input.Name)
	cleanCol := strings.TrimSpace(input.Collection)
	if cleanName == "" || cleanCol == "" {
		return nil, fmt.Errorf("%w: el nombre y la colección del vestido son obligatorios", domain.ErrInvalidInput)
	}

	img := strings.TrimSpace(input.ImagePath)
	img1 := strings.TrimSpace(input.Image1Path)
	if img == "" {
		img = img1
	}
	if img1 == "" {
		img1 = img
	}

	img2 := strings.TrimSpace(input.Image2Path)
	img3 := strings.TrimSpace(input.Image3Path)
	if img3 == "" {
		img3 = img2
	}

	now := time.Now().UTC()
	dress := &domain.Dress{
		ID:          uuid.NewString(),
		Name:        cleanName,
		Collection:  cleanCol,
		ImagePath:   img,
		Image1Path:  img1,
		Image2Path:  img2,
		Image3Path:  img3,
		Description: strings.TrimSpace(input.Description),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, dress); err != nil {
		return nil, fmt.Errorf("failed to create dress: %w", err)
	}

	detail := dress.ToDetailResponse()
	return &detail, nil
}

func (s *DressService) EnsureDefaultDresses(ctx context.Context) error {
	existing, err := s.repo.List(ctx, "")
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	now := time.Now().UTC()
	defaults := []domain.Dress{
		// -------------------------------------------------------------
		// Colección 1: Esencia Floral (5 vestidos)
		// -------------------------------------------------------------
		{
			ID:          uuid.NewString(),
			Name:        "Vestido Magnolia",
			Collection:  "Esencia Floral",
			ImagePath:   "assets/images/esencia-floral/MARIA_DIEZMA_001.jpg",
			Image1Path:  "assets/images/esencia-floral/MARIA_DIEZMA_001.jpg",
			Image2Path:  "assets/images/esencia-floral/MARIA_DIEZMA_002.jpg",
			Image3Path:  "assets/images/esencia-floral/MARIA_DIEZMA_003.jpg",
			Description: "Vestido de corte sirena con bordados florales artesanales en tul y escote corazón.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.NewString(),
			Name:        "Vestido Jazmín",
			Collection:  "Esencia Floral",
			ImagePath:   "assets/images/esencia-floral/MARIA_DIEZMA_004.jpg",
			Image1Path:  "assets/images/esencia-floral/MARIA_DIEZMA_004.jpg",
			Image2Path:  "assets/images/esencia-floral/MARIA_DIEZMA_005.jpg",
			Image3Path:  "assets/images/esencia-floral/MARIA_DIEZMA_006.jpg",
			Description: "Diseño evasé con mangas abullonadas y aplicaciones de pétalos en relieve.",
			CreatedAt:   now.Add(1 * time.Second),
			UpdatedAt:   now.Add(1 * time.Second),
		},
		{
			ID:          uuid.NewString(),
			Name:        "Vestido Dalia",
			Collection:  "Esencia Floral",
			ImagePath:   "assets/images/esencia-floral/MARIA_DIEZMA_007.jpg",
			Image1Path:  "assets/images/esencia-floral/MARIA_DIEZMA_007.jpg",
			Image2Path:  "assets/images/esencia-floral/MARIA_DIEZMA_008.jpg",
			Image3Path:  "assets/images/esencia-floral/MARIA_DIEZMA_009.jpg",
			Description: "Silueta en A con cuerpo drapeado en gasa de seda y estampado floral sutil.",
			CreatedAt:   now.Add(2 * time.Second),
			UpdatedAt:   now.Add(2 * time.Second),
		},
		{
			ID:          uuid.NewString(),
			Name:        "Vestido Camelia",
			Collection:  "Esencia Floral",
			ImagePath:   "assets/images/esencia-floral/MARIA_DIEZMA_010.jpg",
			Image1Path:  "assets/images/esencia-floral/MARIA_DIEZMA_010.jpg",
			Image2Path:  "assets/images/esencia-floral/MARIA_DIEZMA_011.jpg",
			Image3Path:  "assets/images/esencia-floral/MARIA_DIEZMA_012.jpg",
			Description: "Vestido midi con falda de vuelo y cinturón joya bordado a mano.",
			CreatedAt:   now.Add(3 * time.Second),
			UpdatedAt:   now.Add(3 * time.Second),
		},
		{
			ID:          uuid.NewString(),
			Name:        "Vestido Azahar",
			Collection:  "Esencia Floral",
			ImagePath:   "assets/images/esencia-floral/MARIA_DIEZMA_013.jpg",
			Image1Path:  "assets/images/esencia-floral/MARIA_DIEZMA_013.jpg",
			Image2Path:  "assets/images/esencia-floral/MARIA_DIEZMA_014.jpg",
			Image3Path:  "assets/images/esencia-floral/MARIA_DIEZMA_015.jpg",
			Description: "Elegante confección en encaje chantilly y espalda baja con botonadura forrada.",
			CreatedAt:   now.Add(4 * time.Second),
			UpdatedAt:   now.Add(4 * time.Second),
		},

		// -------------------------------------------------------------
		// Colección 2: Atardecer Mediterráneo (5 vestidos)
		// -------------------------------------------------------------
		{
			ID:          uuid.NewString(),
			Name:        "Vestido Siena",
			Collection:  "Atardecer Mediterráneo",
			ImagePath:   "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_016.jpg",
			Image1Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_016.jpg",
			Image2Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_017.jpg",
			Image3Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_018.jpg",
			Description: "Vestido fluido de gasa de seda con espalda descubierta en tono cálido terracota.",
			CreatedAt:   now.Add(5 * time.Second),
			UpdatedAt:   now.Add(5 * time.Second),
		},
		{
			ID:          uuid.NewString(),
			Name:        "Vestido Aurora",
			Collection:  "Atardecer Mediterráneo",
			ImagePath:   "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_019.jpg",
			Image1Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_019.jpg",
			Image2Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_020.jpg",
			Image3Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_021.jpg",
			Description: "Confeccionado en crepé satinado con drapeado asimétrico que evoca la luz dorada.",
			CreatedAt:   now.Add(6 * time.Second),
			UpdatedAt:   now.Add(6 * time.Second),
		},
		{
			ID:          uuid.NewString(),
			Name:        "Vestido Coral",
			Collection:  "Atardecer Mediterráneo",
			ImagePath:   "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_022.jpg",
			Image1Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_022.jpg",
			Image2Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_023.jpg",
			Image3Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_024.jpg",
			Description: "Vestido midi con falda plisada soleil y escote halter cruzado.",
			CreatedAt:   now.Add(7 * time.Second),
			UpdatedAt:   now.Add(7 * time.Second),
		},
		{
			ID:          uuid.NewString(),
			Name:        "Vestido Terracota",
			Collection:  "Atardecer Mediterráneo",
			ImagePath:   "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_025.jpg",
			Image1Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_025.jpg",
			Image2Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_026.jpg",
			Image3Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_027.jpg",
			Description: "Diseño minimalista en lino sedoso con abertura lateral y tirantes espagueti.",
			CreatedAt:   now.Add(8 * time.Second),
			UpdatedAt:   now.Add(8 * time.Second),
		},
		{
			ID:          uuid.NewString(),
			Name:        "Vestido Sol Poniente",
			Collection:  "Atardecer Mediterráneo",
			ImagePath:   "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_028.jpg",
			Image1Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_028.jpg",
			Image2Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_029.jpg",
			Image3Path:  "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_030.jpg",
			Description: "Vestido de fiesta en mikado tornasolado con manga capa desestructurada.",
			CreatedAt:   now.Add(9 * time.Second),
			UpdatedAt:   now.Add(9 * time.Second),
		},
	}

	for _, d := range defaults {
		item := d
		if err := s.repo.Create(ctx, &item); err != nil {
			return fmt.Errorf("failed to seed dress %s: %w", d.Name, err)
		}
	}

	return nil
}
