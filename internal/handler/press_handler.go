package handler

import (
	"net/http"

	"mariadiezmaback/internal/service"
	"mariadiezmaback/pkg/response"
)

type PressHandler struct {
	service *service.PressArticleService
}

func NewPressHandler(service *service.PressArticleService) *PressHandler {
	return &PressHandler{service: service}
}

// List handles requests to return an array of press articles.
func (h *PressHandler) List(w http.ResponseWriter, r *http.Request) {
	articles, err := h.service.List(r.Context())
	if err != nil {
		response.InternalServerError(w, "Error al obtener los artículos de prensa")
		return
	}

	response.OK(w, "Artículos de prensa obtenidos correctamente", articles)
}
