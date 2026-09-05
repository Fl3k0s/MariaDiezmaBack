package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/handler"
	"mariadiezmaback/internal/repository/memory"
	"mariadiezmaback/internal/service"
	"mariadiezmaback/pkg/response"
)

func TestPressHandler_List(t *testing.T) {
	repo := memory.NewPressArticleRepository()
	pressSvc := service.NewPressArticleService(repo)
	_ = pressSvc.EnsureDefaultPressArticles(context.Background())
	pressH := handler.NewPressHandler(pressSvc)

	r := chi.NewRouter()
	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/press", pressH.List)
		api.Get("/prensa", pressH.List)
		api.Get("/articulos-prensa", pressH.List)
	})

	endpoints := []string{"/api/v1/press", "/api/v1/prensa", "/api/v1/articulos-prensa"}

	for _, endpoint := range endpoints {
		t.Run("Endpoint "+endpoint, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, endpoint, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("expected status 200 OK, got %d: %s", rec.Code, rec.Body.String())
			}

			var apiResp response.APIResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &apiResp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if !apiResp.Success {
				t.Errorf("expected success true")
			}

			// Unmarshal into generic slice of maps to inspect raw json keys
			var rawArticles []map[string]any
			dataBytes, _ := json.Marshal(apiResp.Data)
			if err := json.Unmarshal(dataBytes, &rawArticles); err != nil {
				t.Fatalf("failed to unmarshal articles: %v", err)
			}

			if len(rawArticles) == 0 {
				t.Fatalf("expected articles array to be non-empty")
			}

			// Also unmarshal into domain.PressArticleResponse
			var articles []domain.PressArticleResponse
			if err := json.Unmarshal(dataBytes, &articles); err != nil {
				t.Fatalf("failed to unmarshal typed articles: %v", err)
			}

			// Validate all required fields for each article
			for i, raw := range rawArticles {
				// 1. Nombre de la revista
				magazine, ok := raw["nombre_revista"].(string)
				if !ok || magazine == "" {
					t.Errorf("item %d: missing or empty 'nombre_revista'", i)
				}

				// 2. Fecha de publicacion
				pubDate, ok := raw["fecha_publicacion"].(string)
				if !ok || pubDate == "" {
					t.Errorf("item %d: missing or empty 'fecha_publicacion'", i)
				}

				// 3. Titular
				title, ok := raw["titular"].(string)
				if !ok || title == "" {
					t.Errorf("item %d: missing or empty 'titular'", i)
				}

				// 4. Pequeña descripcion
				desc, ok := raw["descripcion"].(string)
				if !ok || desc == "" {
					t.Errorf("item %d: missing or empty 'descripcion'", i)
				}

				// 5. Enlace del articulo (both 'enlace_articulo' and 'enlace' available)
				url1, ok1 := raw["enlace_articulo"].(string)
				url2, ok2 := raw["enlace"].(string)
				if (!ok1 || url1 == "") && (!ok2 || url2 == "") {
					t.Errorf("item %d: missing or empty 'enlace_articulo' or 'enlace'", i)
				}
			}
		})
	}
}
