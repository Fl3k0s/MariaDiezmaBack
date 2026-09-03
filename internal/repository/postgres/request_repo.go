package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"mariadiezmaback/internal/domain"
)

type RequestRepo struct {
	pool *pgxpool.Pool
}

func NewRequestRepository(pool *pgxpool.Pool) *RequestRepo {
	return &RequestRepo{pool: pool}
}

func (r *RequestRepo) Create(ctx context.Context, req *domain.RequestItem) error {
	metadataBytes, err := json.Marshal(req.Metadata)
	if err != nil {
		metadataBytes = []byte("{}")
	}

	query := `
		INSERT INTO requests (
			id, type, status, priority, sender_name, sender_email, sender_phone,
			subject, message, metadata, internal_note, assigned_to, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`
	_, err = r.pool.Exec(ctx, query,
		req.ID,
		req.Type,
		req.Status,
		req.Priority,
		req.SenderName,
		req.SenderEmail,
		req.SenderPhone,
		req.Subject,
		req.Message,
		metadataBytes,
		req.InternalNote,
		req.AssignedTo,
		req.CreatedAt,
		req.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert request: %w", err)
	}
	return nil
}

func (r *RequestRepo) GetByID(ctx context.Context, id string) (*domain.RequestItem, error) {
	query := `
		SELECT 
			id, type, status, priority, sender_name, sender_email, sender_phone,
			subject, message, metadata, internal_note, assigned_to, created_at, updated_at
		FROM requests
		WHERE id = $1
	`
	var req domain.RequestItem
	var metadataBytes []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&req.ID,
		&req.Type,
		&req.Status,
		&req.Priority,
		&req.SenderName,
		&req.SenderEmail,
		&req.SenderPhone,
		&req.Subject,
		&req.Message,
		&metadataBytes,
		&req.InternalNote,
		&req.AssignedTo,
		&req.CreatedAt,
		&req.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get request by id: %w", err)
	}

	if len(metadataBytes) > 0 {
		_ = json.Unmarshal(metadataBytes, &req.Metadata)
	}

	return &req, nil
}

func (r *RequestRepo) List(ctx context.Context, filter domain.RequestFilter) ([]domain.RequestItem, int, error) {
	whereClauses := []string{"1=1"}
	args := []any{}
	argIdx := 1

	if filter.Type != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.Status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Priority != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("priority = $%d", argIdx))
		args = append(args, filter.Priority)
		argIdx++
	}
	if filter.Search != "" {
		likePattern := "%" + strings.ToLower(filter.Search) + "%"
		whereClauses = append(whereClauses, fmt.Sprintf(
			"(LOWER(sender_name) LIKE $%d OR LOWER(sender_email) LIKE $%d OR LOWER(subject) LIKE $%d OR LOWER(message) LIKE $%d)",
			argIdx, argIdx, argIdx, argIdx,
		))
		args = append(args, likePattern)
		argIdx++
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM requests WHERE %s", whereSQL)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count requests: %w", err)
	}

	// Order
	orderCol := "created_at"
	if filter.SortBy == "priority" || filter.SortBy == "status" {
		orderCol = filter.SortBy
	}
	orderDir := "DESC"
	if strings.ToLower(filter.SortOrder) == "asc" {
		orderDir = "ASC"
	}

	// Pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	perPage := filter.PerPage
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	query := fmt.Sprintf(`
		SELECT 
			id, type, status, priority, sender_name, sender_email, sender_phone,
			subject, message, metadata, internal_note, assigned_to, created_at, updated_at
		FROM requests
		WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereSQL, orderCol, orderDir, argIdx, argIdx+1)

	args = append(args, perPage, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query requests: %w", err)
	}
	defer rows.Close()

	var requests []domain.RequestItem
	for rows.Next() {
		var req domain.RequestItem
		var metadataBytes []byte
		if err := rows.Scan(
			&req.ID,
			&req.Type,
			&req.Status,
			&req.Priority,
			&req.SenderName,
			&req.SenderEmail,
			&req.SenderPhone,
			&req.Subject,
			&req.Message,
			&metadataBytes,
			&req.InternalNote,
			&req.AssignedTo,
			&req.CreatedAt,
			&req.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan request: %w", err)
		}
		if len(metadataBytes) > 0 {
			_ = json.Unmarshal(metadataBytes, &req.Metadata)
		}
		requests = append(requests, req)
	}

	return requests, total, nil
}

func (r *RequestRepo) Update(ctx context.Context, req *domain.RequestItem) error {
	metadataBytes, err := json.Marshal(req.Metadata)
	if err != nil {
		metadataBytes = []byte("{}")
	}

	query := `
		UPDATE requests SET
			type = $2,
			status = $3,
			priority = $4,
			sender_name = $5,
			sender_email = $6,
			sender_phone = $7,
			subject = $8,
			message = $9,
			metadata = $10,
			internal_note = $11,
			assigned_to = $12,
			updated_at = $13
		WHERE id = $1
	`
	tag, err := r.pool.Exec(ctx, query,
		req.ID,
		req.Type,
		req.Status,
		req.Priority,
		req.SenderName,
		req.SenderEmail,
		req.SenderPhone,
		req.Subject,
		req.Message,
		metadataBytes,
		req.InternalNote,
		req.AssignedTo,
		req.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update request: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *RequestRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM requests WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete request: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
