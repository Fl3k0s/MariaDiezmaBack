package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/service"
	"mariadiezmaback/pkg/response"
)

type HolidayHandler struct {
	service *service.HolidayService
}

func NewHolidayHandler(service *service.HolidayService) *HolidayHandler {
	return &HolidayHandler{service: service}
}

// List handles GET /api/v1/holidays and /api/v1/festivos.
// Optional query parameter: ?year=2026 (or ?anio=2026 / ?ano=2026).
// If no holidays are found, returns 200 OK with data: [] so the web renders standard dates without issue.
func (h *HolidayHandler) List(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		yearStr = r.URL.Query().Get("anio")
	}
	if yearStr == "" {
		yearStr = r.URL.Query().Get("ano")
	}

	year := 0
	if yearStr != "" {
		y, err := strconv.Atoi(strings.TrimSpace(yearStr))
		if err != nil || y < 1900 || y > 2100 {
			response.BadRequest(w, "El parámetro año debe ser un año válido (ej. 2026)", nil)
			return
		}
		year = y
	}

	holidays, err := h.service.List(r.Context(), year)
	if err != nil {
		response.InternalServerError(w, "Error al obtener los festivos")
		return
	}

	response.OK(w, "Festivos obtenidos correctamente", holidays)
}

// GetByYear handles GET /api/v1/holidays/{year} and /api/v1/festivos/{year}.
// Returns 200 OK with data: [] if no holidays exist for the requested year.
func (h *HolidayHandler) GetByYear(w http.ResponseWriter, r *http.Request) {
	yearParam := chi.URLParam(r, "year")
	year, err := strconv.Atoi(strings.TrimSpace(yearParam))
	if err != nil || year < 1900 || year > 2100 {
		response.BadRequest(w, "El parámetro año debe ser un número válido (ej. 2026)", nil)
		return
	}

	holidays, err := h.service.List(r.Context(), year)
	if err != nil {
		response.InternalServerError(w, "Error al obtener los festivos del año")
		return
	}

	response.OK(w, "Festivos obtenidos correctamente", holidays)
}

// GetByID handles GET /api/v1/holidays/id/{id} and /api/v1/festivos/id/{id}.
func (h *HolidayHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	holiday, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "Festivo no encontrado")
			return
		}
		response.InternalServerError(w, "Error al buscar el festivo")
		return
	}

	response.OK(w, "Festivo obtenido correctamente", holiday)
}

// Create handles POST /api/v1/holidays and /api/v1/festivos.
// Supports both a single holiday object or an array/batch of holidays.
func (h *HolidayHandler) Create(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		response.BadRequest(w, "Error al leer el cuerpo de la petición", err.Error())
		return
	}

	trimmed := bytes.TrimSpace(bodyBytes)
	if len(trimmed) == 0 {
		response.BadRequest(w, "El cuerpo de la petición no puede estar vacío", nil)
		return
	}

	// 1. Array of holidays: [...]
	if trimmed[0] == '[' {
		var inputs []domain.CreateHolidayInput
		if err := json.Unmarshal(trimmed, &inputs); err != nil {
			response.BadRequest(w, "JSON inválido en array de festivos", err.Error())
			return
		}
		if len(inputs) == 0 {
			response.BadRequest(w, "El array de festivos no puede estar vacío", nil)
			return
		}

		created, err := h.service.CreateBatch(r.Context(), inputs)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidInput) {
				response.BadRequest(w, err.Error(), nil)
				return
			}
			response.InternalServerError(w, "Error al crear el lote de festivos")
			return
		}

		response.Created(w, "Lote de festivos creado correctamente", created)
		return
	}

	// 2. Object with wrapped list: { "holidays": [...] } or { "festivos": [...] }
	var wrapped struct {
		Holidays []domain.CreateHolidayInput `json:"holidays"`
		Festivos []domain.CreateHolidayInput `json:"festivos"`
	}
	if err := json.Unmarshal(trimmed, &wrapped); err == nil && (len(wrapped.Holidays) > 0 || len(wrapped.Festivos) > 0) {
		inputs := wrapped.Holidays
		if len(inputs) == 0 {
			inputs = wrapped.Festivos
		}
		created, err := h.service.CreateBatch(r.Context(), inputs)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidInput) {
				response.BadRequest(w, err.Error(), nil)
				return
			}
			response.InternalServerError(w, "Error al crear el lote de festivos")
			return
		}

		response.Created(w, "Lote de festivos creado correctamente", created)
		return
	}

	// 3. Single holiday object: { "date": "...", "name": "..." }
	var single domain.CreateHolidayInput
	if err := json.Unmarshal(trimmed, &single); err != nil {
		response.BadRequest(w, "JSON inválido en festivo", err.Error())
		return
	}

	created, err := h.service.Create(r.Context(), single)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			response.BadRequest(w, err.Error(), nil)
			return
		}
		if errors.Is(err, domain.ErrAlreadyExists) {
			response.BadRequest(w, "Ya existe un festivo registrado para esa fecha", nil)
			return
		}
		response.InternalServerError(w, "Error al crear el festivo")
		return
	}

	response.Created(w, "Festivo creado correctamente", created)
}

// Update handles PUT /api/v1/holidays/{id} and /api/v1/festivos/{id}.
func (h *HolidayHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, "El ID del festivo es obligatorio", nil)
		return
	}

	var input domain.UpdateHolidayInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.BadRequest(w, "JSON inválido en el cuerpo de la petición", err.Error())
		return
	}

	updated, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "Festivo no encontrado")
			return
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			response.BadRequest(w, err.Error(), nil)
			return
		}
		if errors.Is(err, domain.ErrAlreadyExists) {
			response.BadRequest(w, "Ya existe otro festivo registrado con esa fecha", nil)
			return
		}
		response.InternalServerError(w, "Error al actualizar el festivo")
		return
	}

	response.OK(w, "Festivo actualizado correctamente", updated)
}

// Delete handles DELETE /api/v1/holidays/{id} and /api/v1/festivos/{id}.
// Allows deleting a specific holiday at any time.
func (h *HolidayHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, "El ID del festivo es obligatorio", nil)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "Festivo no encontrado")
			return
		}
		response.InternalServerError(w, "Error al eliminar el festivo")
		return
	}

	response.OK(w, "Festivo eliminado correctamente", nil)
}

// DeleteByYear handles DELETE /api/v1/holidays/year/{year} and /api/v1/festivos/year/{year}.
// Allows deleting all holidays for a whole year at any time.
func (h *HolidayHandler) DeleteByYear(w http.ResponseWriter, r *http.Request) {
	yearParam := chi.URLParam(r, "year")
	if yearParam == "" {
		yearParam = r.URL.Query().Get("year")
	}
	if yearParam == "" {
		yearParam = r.URL.Query().Get("anio")
	}

	year, err := strconv.Atoi(strings.TrimSpace(yearParam))
	if err != nil || year < 1900 || year > 2100 {
		response.BadRequest(w, "El parámetro año debe ser un año válido (ej. 2026)", nil)
		return
	}

	count, err := h.service.DeleteByYear(r.Context(), year)
	if err != nil {
		response.InternalServerError(w, "Error al eliminar los festivos del año")
		return
	}

	response.OK(w, fmt.Sprintf("Se han eliminado %d festivo(s) del año %d", count, year), map[string]any{
		"deleted_count": count,
		"year":          year,
	})
}
