package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	"mariadiezmaback/internal/domain"
)

type DressRepo struct {
	mu      sync.RWMutex
	dresses map[string]*domain.Dress
}

func NewDressRepository() *DressRepo {
	return &DressRepo{
		dresses: make(map[string]*domain.Dress),
	}
}

func (r *DressRepo) Create(ctx context.Context, dress *domain.Dress) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.dresses[dress.ID] = dress
	return nil
}

func (r *DressRepo) GetByID(ctx context.Context, id string) (*domain.Dress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	d, exists := r.dresses[id]
	if !exists {
		return nil, domain.ErrNotFound
	}
	dressCopy := *d
	return &dressCopy, nil
}

func (r *DressRepo) List(ctx context.Context, collectionFilter string) ([]domain.Dress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filter := strings.TrimSpace(strings.ToLower(collectionFilter))
	list := make([]domain.Dress, 0)

	for _, d := range r.dresses {
		if filter != "" && !strings.EqualFold(d.Collection, filter) {
			continue
		}
		list = append(list, *d)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})

	return list, nil
}

func (r *DressRepo) GetByNameAndCollection(ctx context.Context, name, collection string) (*domain.Dress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	normName := strings.TrimSpace(strings.ToLower(name))
	normCol := strings.TrimSpace(strings.ToLower(collection))

	for _, d := range r.dresses {
		if strings.EqualFold(strings.TrimSpace(d.Name), normName) &&
			strings.EqualFold(strings.TrimSpace(d.Collection), normCol) {
			dressCopy := *d
			return &dressCopy, nil
		}
	}

	return nil, domain.ErrNotFound
}

