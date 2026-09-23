package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"mariadiezmaback/internal/domain"
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

func (h *CollectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateCollectionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.BadRequest(w, "JSON inválido en el cuerpo de la petición", err.Error())
		return
	}

	created, err := h.service.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			response.BadRequest(w, err.Error(), nil)
			return
		}
		response.InternalServerError(w, "Error al crear la colección")
		return
	}

	response.Created(w, "Colección creada correctamente", created)
}

