package memory

import (
	"context"
	"sort"
	"sync"

	"mariadiezmaback/internal/domain"
)

type PressArticleRepo struct {
	mu       sync.RWMutex
	articles map[string]*domain.PressArticle
}

func NewPressArticleRepository() *PressArticleRepo {
	return &PressArticleRepo{
		articles: make(map[string]*domain.PressArticle),
	}
}

func (r *PressArticleRepo) Create(ctx context.Context, article *domain.PressArticle) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.articles[article.ID] = article
	return nil
}

func (r *PressArticleRepo) GetByID(ctx context.Context, id string) (*domain.PressArticle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	article, exists := r.articles[id]
	if !exists {
		return nil, domain.ErrNotFound
	}
	copyArt := *article
	return &copyArt, nil
}

func (r *PressArticleRepo) List(ctx context.Context) ([]domain.PressArticle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.PressArticle, 0, len(r.articles))
	for _, a := range r.articles {
		list = append(list, *a)
	}

	// Order by ID ASC or created order
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})

	return list, nil
}
