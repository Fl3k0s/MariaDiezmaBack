package service_test

import (
	"context"
	"testing"

	"mariadiezmaback/internal/repository/memory"
	"mariadiezmaback/internal/service"
)

func TestPressArticleService_ListAndEnsureDefaults(t *testing.T) {
	repo := memory.NewPressArticleRepository()
	svc := service.NewPressArticleService(repo)

	ctx := context.Background()

	// 1. Initial list before seeding is empty
	initial, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(initial) != 0 {
		t.Fatalf("expected 0 articles initially, got %d", len(initial))
	}

	// 2. Ensure default press articles
	if err := svc.EnsureDefaultPressArticles(ctx); err != nil {
		t.Fatalf("failed to seed press articles: %v", err)
	}

	// 3. List should return all 5 default articles
	articles, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(articles) != 5 {
		t.Fatalf("expected 5 articles, got %d", len(articles))
	}

	for i, a := range articles {
		if a.MagazineName == "" {
			t.Errorf("article %d: empty MagazineName", i)
		}
		if a.PublicationDate == "" {
			t.Errorf("article %d: empty PublicationDate", i)
		}
		if a.Title == "" {
			t.Errorf("article %d: empty Title", i)
		}
		if a.Description == "" {
			t.Errorf("article %d: empty Description", i)
		}
		if a.ArticleURL == "" {
			t.Errorf("article %d: empty ArticleURL", i)
		}
		if a.Link == "" {
			t.Errorf("article %d: empty Link", i)
		}
	}

	// 4. Calling EnsureDefaultPressArticles again should be idempotent
	if err := svc.EnsureDefaultPressArticles(ctx); err != nil {
		t.Fatalf("subsequent EnsureDefaultPressArticles failed: %v", err)
	}
	articles2, _ := svc.List(ctx)
	if len(articles2) != 5 {
		t.Fatalf("expected 5 articles after second ensure, got %d", len(articles2))
	}
}
