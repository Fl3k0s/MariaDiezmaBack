package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/service"
	"mariadiezmaback/pkg/response"
)

type DressHandler struct {
	service *service.DressService
}

func NewDressHandler(service *service.DressService) *DressHandler {
	return &DressHandler{service: service}
}

func (h *DressHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	collectionFilter := q.Get("coleccion")
	if collectionFilter == "" {
		collectionFilter = q.Get("collection")
	}

	dresses, err := h.service.List(r.Context(), collectionFilter)
	if err != nil {
		response.InternalServerError(w, "Error al obtener los vestidos")
		return
	}

	response.OK(w, "Vestidos obtenidos correctamente", dresses)
}

func (h *DressHandler) GetDetail(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// Extract name (supports query param or URL param)
	name := strings.TrimSpace(q.Get("nombre"))
	if name == "" {
		name = strings.TrimSpace(q.Get("name"))
	}
	if name == "" {
		name = strings.TrimSpace(chi.URLParam(r, "nombre"))
	}

	// Extract collection (supports query param or URL param)
	collection := strings.TrimSpace(q.Get("coleccion"))
	if collection == "" {
		collection = strings.TrimSpace(q.Get("collection"))
	}
	if collection == "" {
		collection = strings.TrimSpace(chi.URLParam(r, "coleccion"))
	}

	if name == "" || collection == "" {
		response.BadRequest(w, "Se requiere el nombre y la colección del vestido", "parámetros 'nombre' y 'coleccion' obligatorios")
		return
	}

	detail, err := h.service.GetByNameAndCollection(r.Context(), name, collection)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "Vestido no encontrado para el nombre y colección indicados")
			return
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			response.BadRequest(w, "Datos de búsqueda inválidos", err.Error())
			return
		}
		response.InternalServerError(w, "Error al obtener los detalles del vestido")
		return
	}

	response.OK(w, "Detalle del vestido obtenido correctamente", detail)
}

func (h *DressHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateDressInput
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
		response.InternalServerError(w, "Error al crear el vestido")
		return
	}

	response.Created(w, "Vestido creado correctamente", created)
}

