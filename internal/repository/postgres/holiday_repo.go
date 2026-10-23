package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"mariadiezmaback/internal/domain"
)

type HolidayRepo struct {
	pool *pgxpool.Pool
}

func NewHolidayRepository(pool *pgxpool.Pool) *HolidayRepo {
	return &HolidayRepo{pool: pool}
}

func (r *HolidayRepo) Create(ctx context.Context, h *domain.Holiday) error {
	query := `
		INSERT INTO holidays (id, date, year, name, description, scope, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (date) DO UPDATE
		SET name = EXCLUDED.name,
		    year = EXCLUDED.year,
		    description = EXCLUDED.description,
		    scope = EXCLUDED.scope,
		    updated_at = EXCLUDED.updated_at
	`
	_, err := r.pool.Exec(ctx, query,
		h.ID,
		h.Date,
		h.Year,
		h.Name,
		h.Description,
		h.Scope,
		h.CreatedAt,
		h.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert holiday: %w", err)
	}
	return nil
}

func (r *HolidayRepo) CreateBatch(ctx context.Context, holidays []domain.Holiday) error {
	if len(holidays) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO holidays (id, date, year, name, description, scope, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (date) DO UPDATE
		SET name = EXCLUDED.name,
		    year = EXCLUDED.year,
		    description = EXCLUDED.description,
		    scope = EXCLUDED.scope,
		    updated_at = EXCLUDED.updated_at
	`

	for _, h := range holidays {
		_, err := tx.Exec(ctx, query,
			h.ID,
			h.Date,
			h.Year,
			h.Name,
			h.Description,
			h.Scope,
			h.CreatedAt,
			h.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert batch holiday %s: %w", h.Date, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit holidays batch: %w", err)
	}

	return nil
}

func (r *HolidayRepo) GetByID(ctx context.Context, id string) (*domain.Holiday, error) {
	query := `
		SELECT id, date, year, name, description, scope, created_at, updated_at
		FROM holidays
		WHERE id = $1
	`
	var h domain.Holiday
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&h.ID,
		&h.Date,
		&h.Year,
		&h.Name,
		&h.Description,
		&h.Scope,
		&h.CreatedAt,
		&h.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get holiday by id: %w", err)
	}
	return &h, nil
}

func (r *HolidayRepo) GetByDate(ctx context.Context, date string) (*domain.Holiday, error) {
	query := `
		SELECT id, date, year, name, description, scope, created_at, updated_at
		FROM holidays
		WHERE date = $1
	`
	var h domain.Holiday
	err := r.pool.QueryRow(ctx, query, date).Scan(
		&h.ID,
		&h.Date,
		&h.Year,
		&h.Name,
		&h.Description,
		&h.Scope,
		&h.CreatedAt,
		&h.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get holiday by date: %w", err)
	}
	return &h, nil
}

func (r *HolidayRepo) List(ctx context.Context, year int) ([]domain.Holiday, error) {
	var (
		query string
		rows  pgx.Rows
		err   error
	)

	if year > 0 {
		query = `
			SELECT id, date, year, name, description, scope, created_at, updated_at
			FROM holidays
			WHERE year = $1
			ORDER BY date ASC
		`
		rows, err = r.pool.Query(ctx, query, year)
	} else {
		query = `
			SELECT id, date, year, name, description, scope, created_at, updated_at
			FROM holidays
			ORDER BY date ASC
		`
		rows, err = r.pool.Query(ctx, query)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query holidays: %w", err)
	}
	defer rows.Close()

	holidays := make([]domain.Holiday, 0)
	for rows.Next() {
		var h domain.Holiday
		if err := rows.Scan(
			&h.ID,
			&h.Date,
			&h.Year,
			&h.Name,
			&h.Description,
			&h.Scope,
			&h.CreatedAt,
			&h.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan holiday row: %w", err)
		}
		holidays = append(holidays, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during holiday rows iteration: %w", err)
	}

	return holidays, nil
}

func (r *HolidayRepo) Update(ctx context.Context, h *domain.Holiday) error {
	query := `
		UPDATE holidays
		SET date = $2,
		    year = $3,
		    name = $4,
		    description = $5,
		    scope = $6,
		    updated_at = $7
		WHERE id = $1
	`
	tag, err := r.pool.Exec(ctx, query,
		h.ID,
		h.Date,
		h.Year,
		h.Name,
		h.Description,
		h.Scope,
		h.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update holiday: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *HolidayRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM holidays WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete holiday: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *HolidayRepo) DeleteByYear(ctx context.Context, year int) (int64, error) {
	query := `DELETE FROM holidays WHERE year = $1`
	tag, err := r.pool.Exec(ctx, query, year)
	if err != nil {
		return 0, fmt.Errorf("failed to delete holidays for year %d: %w", year, err)
	}
	return tag.RowsAffected(), nil
}
