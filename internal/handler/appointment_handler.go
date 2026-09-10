package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/service"
	"mariadiezmaback/pkg/response"
)

// AppointmentHandler handles web requests for booking appointments.
// Supports web fields:
// - Tipo de cita (string)
// - Fecha (string)
// - Franja horaria (string)
// - Nombre y apellidos (string)
// - Teléfono de contacto (string)
// - Fecha estimada (date, nullable)
// - Detalles (string)
type AppointmentHandler struct {
	service *service.AppointmentService
}

func NewAppointmentHandler(service *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{service: service}
}

// Create handles POST /api/v1/citas or /api/v1/appointments.
// Accepts both Spanish and English JSON payloads from the web form.
func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateAppointmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.BadRequest(w, "JSON inválido en el cuerpo de la petición", err.Error())
		return
	}

	appt, err := h.service.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			response.BadRequest(w, "Error de validación en los datos de la cita", err.Error())
			return
		}
		response.InternalServerError(w, "Error al procesar la solicitud de cita")
		return
	}

	response.Created(w, "Cita solicitada correctamente", appt)
}
