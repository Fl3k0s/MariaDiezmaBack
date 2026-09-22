package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"mariadiezmaback/internal/domain"
)

type DressRepo struct {
	pool *pgxpool.Pool
}

func NewDressRepository(pool *pgxpool.Pool) *DressRepo {
	return &DressRepo{pool: pool}
}

func (r *DressRepo) Create(ctx context.Context, dress *domain.Dress) error {
	img1 := dress.Image1Path
	if img1 == "" {
		img1 = dress.ImagePath
	}
	img := dress.ImagePath
	if img == "" {
		img = img1
	}
	img3 := dress.Image3Path
	if img3 == "" {
		img3 = dress.Image2Path
	}

	query := `
		INSERT INTO dresses (id, name, collection, image_path, image1_path, image2_path, image3_path, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query,
		dress.ID,
		dress.Name,
		dress.Collection,
		img,
		img1,
		dress.Image2Path,
		img3,
		dress.Description,
		dress.CreatedAt,
		dress.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert dress: %w", err)
	}
	return nil
}

func (r *DressRepo) GetByID(ctx context.Context, id string) (*domain.Dress, error) {
	query := `
		SELECT id, name, collection, image_path, image1_path, image2_path, image3_path, description, created_at, updated_at
		FROM dresses
		WHERE id = $1
	`
	var d domain.Dress
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&d.ID,
		&d.Name,
		&d.Collection,
		&d.ImagePath,
		&d.Image1Path,
		&d.Image2Path,
		&d.Image3Path,
		&d.Description,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get dress by id: %w", err)
	}
	return &d, nil
}

func (r *DressRepo) GetByNameAndCollection(ctx context.Context, name, collection string) (*domain.Dress, error) {
	query := `
		SELECT id, name, collection, image_path, image1_path, image2_path, image3_path, description, created_at, updated_at
		FROM dresses
		WHERE LOWER(name) = LOWER($1) AND LOWER(collection) = LOWER($2)
		LIMIT 1
	`
	var d domain.Dress
	err := r.pool.QueryRow(ctx, query, strings.TrimSpace(name), strings.TrimSpace(collection)).Scan(
		&d.ID,
		&d.Name,
		&d.Collection,
		&d.ImagePath,
		&d.Image1Path,
		&d.Image2Path,
		&d.Image3Path,
		&d.Description,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get dress by name and collection: %w", err)
	}
	return &d, nil
}

func (r *DressRepo) List(ctx context.Context, collectionFilter string) ([]domain.Dress, error) {
	query := `
		SELECT id, name, collection, image_path, image1_path, image2_path, image3_path, description, created_at, updated_at
		FROM dresses
	`
	var args []any

	if strings.TrimSpace(collectionFilter) != "" {
		query += ` WHERE LOWER(collection) = LOWER($1)`
		args = append(args, strings.TrimSpace(collectionFilter))
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query dresses: %w", err)
	}
	defer rows.Close()

	var dresses []domain.Dress
	for rows.Next() {
		var d domain.Dress
		if err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Collection,
			&d.ImagePath,
			&d.Image1Path,
			&d.Image2Path,
			&d.Image3Path,
			&d.Description,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan dress: %w", err)
		}
		dresses = append(dresses, d)
	}

	return dresses, nil
}
