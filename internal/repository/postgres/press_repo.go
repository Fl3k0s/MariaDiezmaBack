package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"mariadiezmaback/internal/domain"
)

type PressArticleRepo struct {
	pool *pgxpool.Pool
}

func NewPressArticleRepository(pool *pgxpool.Pool) *PressArticleRepo {
	return &PressArticleRepo{pool: pool}
}

func (r *PressArticleRepo) Create(ctx context.Context, article *domain.PressArticle) error {
	query := `
		INSERT INTO press_articles (id, magazine_name, publication_date, title, description, article_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	url := article.ArticleURL
	if url == "" {
		url = article.URL
	}
	_, err := r.pool.Exec(ctx, query,
		article.ID,
		article.MagazineName,
		article.PublicationDate,
		article.Title,
		article.Description,
		url,
	)
	if err != nil {
		return fmt.Errorf("failed to insert press article: %w", err)
	}
	return nil
}

func (r *PressArticleRepo) GetByID(ctx context.Context, id string) (*domain.PressArticle, error) {
	query := `
		SELECT id, magazine_name, publication_date, title, description, article_url
		FROM press_articles
		WHERE id = $1
	`
	var a domain.PressArticle
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&a.ID,
		&a.MagazineName,
		&a.PublicationDate,
		&a.Title,
		&a.Description,
		&a.ArticleURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get press article by id: %w", err)
	}
	a.URL = a.ArticleURL
	return &a, nil
}

func (r *PressArticleRepo) List(ctx context.Context) ([]domain.PressArticle, error) {
	query := `
		SELECT id, magazine_name, publication_date, title, description, article_url
		FROM press_articles
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query press articles: %w", err)
	}
	defer rows.Close()

	var articles []domain.PressArticle
	for rows.Next() {
		var a domain.PressArticle
		if err := rows.Scan(
			&a.ID,
			&a.MagazineName,
			&a.PublicationDate,
			&a.Title,
			&a.Description,
			&a.ArticleURL,
		); err != nil {
			return nil, fmt.Errorf("failed to scan press article: %w", err)
		}
		a.URL = a.ArticleURL
		articles = append(articles, a)
	}

	return articles, nil
}
