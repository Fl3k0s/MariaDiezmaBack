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

type HolidayService struct {
	repo repository.HolidayRepository
}

func NewHolidayService(repo repository.HolidayRepository) *HolidayService {
	return &HolidayService{repo: repo}
}

// List returns all holidays or filters by year.
// If no holidays are found, it returns an empty slice (never null) so the frontend renders normal days without errors.
func (s *HolidayService) List(ctx context.Context, year int) ([]domain.HolidayResponse, error) {
	holidays, err := s.repo.List(ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch holidays: %w", err)
	}

	result := make([]domain.HolidayResponse, 0, len(holidays))
	for _, h := range holidays {
		result = append(result, h.ToResponse())
	}

	return result, nil
}

// GetByID returns a holiday by its unique ID.
func (s *HolidayService) GetByID(ctx context.Context, id string) (*domain.HolidayResponse, error) {
	h, err := s.repo.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	resp := h.ToResponse()
	return &resp, nil
}

// Create creates a new holiday with validation.
func (s *HolidayService) Create(ctx context.Context, input domain.CreateHolidayInput) (*domain.HolidayResponse, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	holiday := &domain.Holiday{
		ID:          uuid.NewString(),
		Date:        input.Date,
		Year:        input.Year,
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
		Scope:       strings.TrimSpace(input.Scope),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, holiday); err != nil {
		return nil, fmt.Errorf("failed to create holiday: %w", err)
	}

	resp := holiday.ToResponse()
	return &resp, nil
}

// CreateBatch creates or registers multiple holidays at once (e.g. for an entire year).
func (s *HolidayService) CreateBatch(ctx context.Context, inputs []domain.CreateHolidayInput) ([]domain.HolidayResponse, error) {
	if len(inputs) == 0 {
		return make([]domain.HolidayResponse, 0), nil
	}

	now := time.Now().UTC()
	holidays := make([]domain.Holiday, 0, len(inputs))
	responses := make([]domain.HolidayResponse, 0, len(inputs))

	for i, in := range inputs {
		if err := in.Validate(); err != nil {
			return nil, fmt.Errorf("error in item %d: %w", i+1, err)
		}

		h := domain.Holiday{
			ID:          uuid.NewString(),
			Date:        in.Date,
			Year:        in.Year,
			Name:        strings.TrimSpace(in.Name),
			Description: strings.TrimSpace(in.Description),
			Scope:       strings.TrimSpace(in.Scope),
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		holidays = append(holidays, h)
		responses = append(responses, h.ToResponse())
	}

	if err := s.repo.CreateBatch(ctx, holidays); err != nil {
		return nil, fmt.Errorf("failed to create batch holidays: %w", err)
	}

	return responses, nil
}

// Update updates an existing holiday by ID.
func (s *HolidayService) Update(ctx context.Context, id string, input domain.UpdateHolidayInput) (*domain.HolidayResponse, error) {
	createIn := domain.CreateHolidayInput{
		Date:        input.Date,
		Year:        input.Year,
		Name:        input.Name,
		Description: input.Description,
		Scope:       input.Scope,
	}
	if err := createIn.Validate(); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Date = createIn.Date
	existing.Year = createIn.Year
	existing.Name = strings.TrimSpace(createIn.Name)
	existing.Description = strings.TrimSpace(createIn.Description)
	existing.Scope = strings.TrimSpace(createIn.Scope)
	existing.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update holiday: %w", err)
	}

	resp := existing.ToResponse()
	return &resp, nil
}

// Delete removes a holiday by ID.
func (s *HolidayService) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: el ID del festivo es obligatorio", domain.ErrInvalidInput)
	}
	return s.repo.Delete(ctx, id)
}

// DeleteByYear removes all holidays for a given year.
func (s *HolidayService) DeleteByYear(ctx context.Context, year int) (int64, error) {
	if year < 1900 || year > 2100 {
		return 0, fmt.Errorf("%w: año no válido (%d)", domain.ErrInvalidInput, year)
	}
	return s.repo.DeleteByYear(ctx, year)
}

// EnsureDefaultHolidays seeds official Madrid holidays for 2025 and 2026 if the repository is empty.
func (s *HolidayService) EnsureDefaultHolidays(ctx context.Context) error {
	existing, err := s.repo.List(ctx, 0)
	if err != nil {
		return fmt.Errorf("failed checking holidays: %w", err)
	}

	if len(existing) > 0 {
		return nil // Already has holidays
	}

	defaults := domain.DefaultMadridHolidays()
	if err := s.repo.CreateBatch(ctx, defaults); err != nil {
		return fmt.Errorf("failed seeding default holidays: %w", err)
	}

	return nil
}
