package memory

import (
	"context"
	"sort"
	"sync"

	"mariadiezmaback/internal/domain"
)

type CollectionRepo struct {
	mu          sync.RWMutex
	collections map[string]*domain.Collection
}

func NewCollectionRepository() *CollectionRepo {
	return &CollectionRepo{
		collections: make(map[string]*domain.Collection),
	}
}

func (r *CollectionRepo) Create(ctx context.Context, col *domain.Collection) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.collections[col.ID] = col
	return nil
}

func (r *CollectionRepo) GetByID(ctx context.Context, id string) (*domain.Collection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	col, exists := r.collections[id]
	if !exists {
		return nil, domain.ErrNotFound
	}
	colCopy := *col
	return &colCopy, nil
}

func (r *CollectionRepo) List(ctx context.Context) ([]domain.Collection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.Collection, 0, len(r.collections))
	for _, col := range r.collections {
		list = append(list, *col)
	}

	// Order by CreatedAt DESC (newest first)
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.After(list[j].CreatedAt)
	})

	return list, nil
}
