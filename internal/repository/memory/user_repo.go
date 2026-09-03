package memory

import (
	"context"
	"strings"
	"sync"

	"mariadiezmaback/internal/domain"
)

type UserRepo struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewUserRepository() *UserRepo {
	return &UserRepo{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.users {
		if strings.EqualFold(existing.Email, user.Email) {
			return domain.ErrAlreadyExists
		}
	}

	r.users[user.ID] = user
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, exists := r.users[id]
	if !exists {
		return nil, domain.ErrNotFound
	}
	// return a copy
	userCopy := *u
	return &userCopy, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if strings.EqualFold(u.Email, email) {
			userCopy := *u
			return &userCopy, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (r *UserRepo) List(ctx context.Context) ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.User, 0, len(r.users))
	for _, u := range r.users {
		list = append(list, *u)
	}
	return list, nil
}
