package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"mariadiezmaback/internal/domain"
)

type CollectionRepo struct {
	pool *pgxpool.Pool
}

func NewCollectionRepository(pool *pgxpool.Pool) *CollectionRepo {
	return &CollectionRepo{pool: pool}
}

func (r *CollectionRepo) Create(ctx context.Context, col *domain.Collection) error {
	query := `
		INSERT INTO collections (id, name, image_path, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		col.ID,
		col.Name,
		col.ImagePath,
		col.Description,
		col.CreatedAt,
		col.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert collection: %w", err)
	}
	return nil
}

func (r *CollectionRepo) GetByID(ctx context.Context, id string) (*domain.Collection, error) {
	query := `
		SELECT id, name, image_path, description, created_at, updated_at
		FROM collections
		WHERE id = $1
	`
	var col domain.Collection
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&col.ID,
		&col.Name,
		&col.ImagePath,
		&col.Description,
		&col.CreatedAt,
		&col.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get collection by id: %w", err)
	}
	return &col, nil
}

func (r *CollectionRepo) List(ctx context.Context) ([]domain.Collection, error) {
	query := `
		SELECT id, name, image_path, description, created_at, updated_at
		FROM collections
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query collections: %w", err)
	}
	defer rows.Close()

	var collections []domain.Collection
	for rows.Next() {
		var col domain.Collection
		if err := rows.Scan(
			&col.ID,
			&col.Name,
			&col.ImagePath,
			&col.Description,
			&col.CreatedAt,
			&col.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan collection: %w", err)
		}
		collections = append(collections, col)
	}

	return collections, nil
}
