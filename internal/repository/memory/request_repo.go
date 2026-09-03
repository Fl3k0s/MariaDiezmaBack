package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	"mariadiezmaback/internal/domain"
)

type RequestRepo struct {
	mu       sync.RWMutex
	requests map[string]*domain.RequestItem
}

func NewRequestRepository() *RequestRepo {
	return &RequestRepo{
		requests: make(map[string]*domain.RequestItem),
	}
}

func (r *RequestRepo) Create(ctx context.Context, req *domain.RequestItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.requests[req.ID] = req
	return nil
}

func (r *RequestRepo) GetByID(ctx context.Context, id string) (*domain.RequestItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	req, exists := r.requests[id]
	if !exists {
		return nil, domain.ErrNotFound
	}
	reqCopy := *req
	return &reqCopy, nil
}

func (r *RequestRepo) List(ctx context.Context, filter domain.RequestFilter) ([]domain.RequestItem, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]*domain.RequestItem, 0)
	searchTerm := strings.ToLower(strings.TrimSpace(filter.Search))

	for _, item := range r.requests {
		if filter.Type != "" && item.Type != filter.Type {
			continue
		}
		if filter.Status != "" && item.Status != filter.Status {
			continue
		}
		if filter.Priority != "" && item.Priority != filter.Priority {
			continue
		}
		if searchTerm != "" {
			nameMatch := strings.Contains(strings.ToLower(item.SenderName), searchTerm)
			emailMatch := strings.Contains(strings.ToLower(item.SenderEmail), searchTerm)
			subjectMatch := strings.Contains(strings.ToLower(item.Subject), searchTerm)
			messageMatch := strings.Contains(strings.ToLower(item.Message), searchTerm)
			if !nameMatch && !emailMatch && !subjectMatch && !messageMatch {
				continue
			}
		}
		filtered = append(filtered, item)
	}

	total := len(filtered)

	// Sort default by CreatedAt desc
	sort.Slice(filtered, func(i, j int) bool {
		if strings.ToLower(filter.SortOrder) == "asc" {
			return filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
		}
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	// Pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	perPage := filter.PerPage
	if perPage < 1 {
		perPage = 10
	}

	start := (page - 1) * perPage
	if start > total {
		return []domain.RequestItem{}, total, nil
	}

	end := start + perPage
	if end > total {
		end = total
	}

	paginated := filtered[start:end]
	result := make([]domain.RequestItem, len(paginated))
	for i, item := range paginated {
		result[i] = *item
	}

	return result, total, nil
}

func (r *RequestRepo) Update(ctx context.Context, req *domain.RequestItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.requests[req.ID]; !exists {
		return domain.ErrNotFound
	}
	r.requests[req.ID] = req
	return nil
}

func (r *RequestRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.requests[id]; !exists {
		return domain.ErrNotFound
	}
	delete(r.requests, id)
	return nil
}
