package service

import (
	"context"
	"fmt"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/repository"
)

type PressArticleService struct {
	repo repository.PressArticleRepository
}

func NewPressArticleService(repo repository.PressArticleRepository) *PressArticleService {
	return &PressArticleService{repo: repo}
}

func (s *PressArticleService) List(ctx context.Context) ([]domain.PressArticleResponse, error) {
	articles, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch press articles: %w", err)
	}

	result := make([]domain.PressArticleResponse, len(articles))
	for i, a := range articles {
		result[i] = a.ToResponse()
	}
	return result, nil
}

func (s *PressArticleService) EnsureDefaultPressArticles(ctx context.Context) error {
	existing, err := s.repo.List(ctx)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	defaults := DefaultPressArticles()
	for _, a := range defaults {
		article := a
		if err := s.repo.Create(ctx, &article); err != nil {
			return fmt.Errorf("failed to seed press article %s: %w", a.Title, err)
		}
	}

	return nil
}

func DefaultPressArticles() []domain.PressArticle {
	return []domain.PressArticle{
		{
			ID:              "50000000-0000-0000-0000-000000000001",
			MagazineName:    "Vogue España",
			PublicationDate: "2024-05-15",
			Title:           "María Diezma: La nueva era de la alta costura nupcial y la artesanía contemporánea",
			Description:     "Un recorrido íntimo por el atelier madrileño de María Diezma, donde cada puntada rinde homenaje a la tradición y al patronaje a medida.",
			ArticleURL:      "https://www.vogue.es/novias/articulos/maria-diezma-alta-costura-nupcial",
		},
		{
			ID:              "50000000-0000-0000-0000-000000000002",
			MagazineName:    "Telva Novias",
			PublicationDate: "2024-03-20",
			Title:           "Las diseñadoras que están transformando los vestidos de novia en España",
			Description:     "Descubre las propuestas vanguardistas y elegantes de María Diezma en su última colección, destacando por sus tejidos nobles y caídas impecables.",
			ArticleURL:      "https://www.telva.com/novias/disenadoras-vestidos-novia-espana-maria-diezma",
		},
		{
			ID:              "50000000-0000-0000-0000-000000000003",
			MagazineName:    "Harper's Bazaar",
			PublicationDate: "2023-11-10",
			Title:           "El romanticismo reinventado: Los detalles joya y bordados botánicos de María Diezma",
			Description:     "Un análisis de los vestidos de fiesta y ceremonia que combinan cortes clásicos con detalles de alta costura únicos.",
			ArticleURL:      "https://www.harpersbazaar.com/es/moda/noticias-moda/maria-diezma-romanticismo-reinventado",
		},
		{
			ID:              "50000000-0000-0000-0000-000000000004",
			MagazineName:    "Elle Gourmet & Lifestyle",
			PublicationDate: "2023-09-05",
			Title:           "Guía de estilo nupcial: Claves para acertar con un vestido a medida firmado por María Diezma",
			Description:     "Consejos exclusivos y el proceso de creación personalizada desde el boceto inicial hasta la última prueba en el taller.",
			ArticleURL:      "https://www.elle.com/es/bodas/novias/vestidos-a-medida-maria-diezma",
		},
		{
			ID:              "50000000-0000-0000-0000-000000000005",
			MagazineName:    "¡HOLA! Novias",
			PublicationDate: "2023-06-18",
			Title:           "Espaldas infinitas y tejidos de ensueño en la nueva propuesta de María Diezma",
			Description:     "La creadora desvela las claves de su colección nupcial pensada para novias que buscan autenticidad, feminidad y elegancia atemporal.",
			ArticleURL:      "https://www.hola.com/novias/articulos/maria-diezma-coleccion-espaldas-tejidos",
		},
	}
}
