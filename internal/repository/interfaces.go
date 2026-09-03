package repository

import (
	"context"

	"mariadiezmaback/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
}

type RequestRepository interface {
	Create(ctx context.Context, req *domain.RequestItem) error
	GetByID(ctx context.Context, id string) (*domain.RequestItem, error)
	List(ctx context.Context, filter domain.RequestFilter) ([]domain.RequestItem, int, error)
	Update(ctx context.Context, req *domain.RequestItem) error
	Delete(ctx context.Context, id string) error
}

type CollectionRepository interface {
	List(ctx context.Context) ([]domain.Collection, error)
	GetByID(ctx context.Context, id string) (*domain.Collection, error)
	Create(ctx context.Context, col *domain.Collection) error
}

type DressRepository interface {
	List(ctx context.Context, collectionFilter string) ([]domain.Dress, error)
	GetByID(ctx context.Context, id string) (*domain.Dress, error)
	GetByNameAndCollection(ctx context.Context, name, collection string) (*domain.Dress, error)
	Create(ctx context.Context, dress *domain.Dress) error
}


