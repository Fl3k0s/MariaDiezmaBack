package memory

import (
	"context"
	"sort"
	"sync"

	"mariadiezmaback/internal/domain"
)

type HolidayRepo struct {
	mu       sync.RWMutex
	holidays map[string]*domain.Holiday
}

func NewHolidayRepository() *HolidayRepo {
	return &HolidayRepo{
		holidays: make(map[string]*domain.Holiday),
	}
}

func (r *HolidayRepo) Create(ctx context.Context, holiday *domain.Holiday) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if a holiday with same date already exists
	for _, h := range r.holidays {
		if h.Date == holiday.Date && h.ID != holiday.ID {
			return domain.ErrAlreadyExists
		}
	}

	hCopy := *holiday
	r.holidays[holiday.ID] = &hCopy
	return nil
}

func (r *HolidayRepo) CreateBatch(ctx context.Context, holidays []domain.Holiday) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, holiday := range holidays {
		hCopy := holiday
		r.holidays[holiday.ID] = &hCopy
	}
	return nil
}

func (r *HolidayRepo) GetByID(ctx context.Context, id string) (*domain.Holiday, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	h, exists := r.holidays[id]
	if !exists {
		return nil, domain.ErrNotFound
	}
	hCopy := *h
	return &hCopy, nil
}

func (r *HolidayRepo) GetByDate(ctx context.Context, date string) (*domain.Holiday, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, h := range r.holidays {
		if h.Date == date {
			hCopy := *h
			return &hCopy, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (r *HolidayRepo) List(ctx context.Context, year int) ([]domain.Holiday, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.Holiday, 0)
	for _, h := range r.holidays {
		if year > 0 && h.Year != year {
			continue
		}
		list = append(list, *h)
	}

	// Order by Date ASC (chronological)
	sort.Slice(list, func(i, j int) bool {
		return list[i].Date < list[j].Date
	})

	return list, nil
}

func (r *HolidayRepo) Update(ctx context.Context, holiday *domain.Holiday) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.holidays[holiday.ID]
	if !exists {
		return domain.ErrNotFound
	}

	// Check if changing date conflicts with another record
	for id, h := range r.holidays {
		if id != holiday.ID && h.Date == holiday.Date {
			return domain.ErrAlreadyExists
		}
	}

	existing.Date = holiday.Date
	existing.Year = holiday.Year
	existing.Name = holiday.Name
	existing.Description = holiday.Description
	existing.Scope = holiday.Scope
	existing.UpdatedAt = holiday.UpdatedAt

	return nil
}

func (r *HolidayRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.holidays[id]; !exists {
		return domain.ErrNotFound
	}

	delete(r.holidays, id)
	return nil
}

func (r *HolidayRepo) DeleteByYear(ctx context.Context, year int) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var deletedCount int64
	for id, h := range r.holidays {
		if h.Year == year {
			delete(r.holidays, id)
			deletedCount++
		}
	}

	return deletedCount, nil
}
