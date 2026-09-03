package handler

import (
	"net/http"

	"mariadiezmaback/internal/service"
	"mariadiezmaback/pkg/response"
)

type CollectionHandler struct {
	service *service.CollectionService
}

func NewCollectionHandler(service *service.CollectionService) *CollectionHandler {
	return &CollectionHandler{service: service}
}

func (h *CollectionHandler) List(w http.ResponseWriter, r *http.Request) {
	collections, err := h.service.List(r.Context())
	if err != nil {
		response.InternalServerError(w, "Error al obtener las colecciones")
		return
	}

	response.OK(w, "Colecciones obtenidas correctamente", collections)
}
