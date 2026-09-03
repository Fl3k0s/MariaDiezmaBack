package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/service"
	"mariadiezmaback/pkg/response"
)

type RequestHandler struct {
	reqService *service.RequestService
}

func NewRequestHandler(reqService *service.RequestService) *RequestHandler {
	return &RequestHandler{reqService: reqService}
}

func (h *RequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.BadRequest(w, "invalid request body", err.Error())
		return
	}

	item, err := h.reqService.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			response.BadRequest(w, "validation failed", err.Error())
			return
		}
		response.InternalServerError(w, "failed to create request")
		return
	}

	response.Created(w, "request submitted successfully", item)
}

func (h *RequestHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("per_page"))

	filter := domain.RequestFilter{
		Type:      q.Get("type"),
		Status:    domain.RequestStatus(q.Get("status")),
		Priority:  domain.RequestPriority(q.Get("priority")),
		Search:    q.Get("search"),
		Page:      page,
		PerPage:   perPage,
		SortBy:    q.Get("sort_by"),
		SortOrder: q.Get("sort_order"),
	}

	items, total, err := h.reqService.List(r.Context(), filter)
	if err != nil {
		response.InternalServerError(w, "failed to fetch requests")
		return
	}

	p := filter.Page
	if p < 1 {
		p = 1
	}
	pp := filter.PerPage
	if pp < 1 {
		pp = 20
	}

	response.Paginated(w, items, total, p, pp)
}

func (h *RequestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, err := h.reqService.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "request not found")
			return
		}
		response.InternalServerError(w, "failed to fetch request")
		return
	}

	response.OK(w, "request retrieved", item)
}

func (h *RequestHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input domain.UpdateStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.BadRequest(w, "invalid request body", err.Error())
		return
	}

	item, err := h.reqService.UpdateStatus(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "request not found")
			return
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			response.BadRequest(w, "invalid status provided", err.Error())
			return
		}
		response.InternalServerError(w, "failed to update status")
		return
	}

	response.OK(w, "status updated successfully", item)
}

func (h *RequestHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input domain.UpdateRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.BadRequest(w, "invalid request body", err.Error())
		return
	}

	item, err := h.reqService.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "request not found")
			return
		}
		response.InternalServerError(w, "failed to update request")
		return
	}

	response.OK(w, "request updated successfully", item)
}

func (h *RequestHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.reqService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "request not found")
			return
		}
		response.InternalServerError(w, "failed to delete request")
		return
	}

	response.OK(w, "request deleted successfully", nil)
}
