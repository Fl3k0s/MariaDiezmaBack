package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Appointment struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email,omitempty"`
	Phone         string    `json:"phone"`
	Date          string    `json:"date"`
	TimeSlot      string    `json:"time_slot"`
	Type          string    `json:"type"`
	EstimatedDate *string   `json:"estimated_date"`
	Details       string    `json:"details,omitempty"`
	Status        string    `json:"status"` // pending, confirmed, cancelled
	CreatedAt     time.Time `json:"created_at"`
}

type CreateAppointmentInput struct {
	Name          string  `json:"name"`
	Email         string  `json:"email,omitempty"`
	Phone         string  `json:"phone"`
	Date          string  `json:"date"`
	TimeSlot      string  `json:"time_slot"`
	Type          string  `json:"type"`
	EstimatedDate *string `json:"estimated_date"`
	Details       string  `json:"details,omitempty"`
}

// UnmarshalJSON allows accepting both Spanish and English JSON field names
func (c *CreateAppointmentInput) UnmarshalJSON(data []byte) error {
	var raw struct {
		// Spanish field names (snake_case)
		TipoCita         string `json:"tipo_cita"`
		Tipo             string `json:"tipo"`
		Fecha            string `json:"fecha"`
		FranjaHoraria    string `json:"franja_horaria"`
		FranaHoraria     string `json:"frana_horaria"`
		TramoHorario     string `json:"tramo_horario"`
		NombreApellidos  string `json:"nombre_apellidos"`
		NombreYApellidos string `json:"nombre_y_apellidos"`
		Nombre           string `json:"nombre"`
		TelefonoContacto string `json:"telefono_contacto"`
		Telefono         string `json:"telefono"`
		FechaEstimada    any    `json:"fecha_estimada"`
		Detalles         string `json:"detalles"`
		Comentarios      string `json:"comentarios"`
		Mensaje          string `json:"mensaje"`
		EmailES                string `json:"email"`
		Mail                   string `json:"mail"`
		Correo                 string `json:"correo"`
		CorreoElectronico      string `json:"correo_electronico"`
		EmailContacto          string `json:"email_contacto"`

		// Spanish field names (camelCase)
		TipoCitaCamel          string `json:"tipoCita"`
		FranjaHorariaCamel     string `json:"franjaHoraria"`
		FranaHorariaCamel      string `json:"franaHoraria"`
		TramoHorarioCamel      string `json:"tramoHorario"`
		NombreApellidosCamel   string `json:"nombreApellidos"`
		NombreYApellidosCamel  string `json:"nombreYApellidos"`
		TelefonoContactoCamel  string `json:"telefonoContacto"`
		FechaEstimadaCamel     any    `json:"fechaEstimada"`
		CorreoElectronicoCamel string `json:"correoElectronico"`
		EmailContactoCamel     string `json:"emailContacto"`

		// English field names
		Type               string `json:"type"`
		Date               string `json:"date"`
		TimeSlot           string `json:"time_slot"`
		TimeSlotCamel      string `json:"timeSlot"`
		Name               string `json:"name"`
		Phone              string `json:"phone"`
		EstimatedDate      any    `json:"estimated_date"`
		EstimatedDateCamel any    `json:"estimatedDate"`
		Details            string `json:"details"`
		EmailEN            string `json:"email_address"`
		EmailAddressCamel  string `json:"emailAddress"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// 1. Tipo de cita
	for _, t := range []string{
		raw.TipoCita,
		raw.TipoCitaCamel,
		raw.Tipo,
		raw.Type,
	} {
		if trimmed := strings.TrimSpace(t); trimmed != "" {
			c.Type = trimmed
			break
		}
	}

	// 2. Fecha
	for _, d := range []string{
		raw.Fecha,
		raw.Date,
	} {
		if trimmed := strings.TrimSpace(d); trimmed != "" {
			c.Date = trimmed
			break
		}
	}

	// 3. Franja horaria (o tramo horario)
	for _, slot := range []string{
		raw.FranjaHoraria,
		raw.FranjaHorariaCamel,
		raw.FranaHoraria,
		raw.FranaHorariaCamel,
		raw.TramoHorario,
		raw.TramoHorarioCamel,
		raw.TimeSlot,
		raw.TimeSlotCamel,
	} {
		if trimmed := strings.TrimSpace(slot); trimmed != "" {
			c.TimeSlot = trimmed
			break
		}
	}

	// 4. Nombre y apellidos
	for _, n := range []string{
		raw.NombreApellidos,
		raw.NombreApellidosCamel,
		raw.NombreYApellidos,
		raw.NombreYApellidosCamel,
		raw.Nombre,
		raw.Name,
	} {
		if trimmed := strings.TrimSpace(n); trimmed != "" {
			c.Name = trimmed
			break
		}
	}

	// 5. Teléfono de contacto
	for _, p := range []string{
		raw.TelefonoContacto,
		raw.TelefonoContactoCamel,
		raw.Telefono,
		raw.Phone,
	} {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			c.Phone = trimmed
			break
		}
	}

	// 6. Fecha estimada (puede venir null o fecha)
	for _, ed := range []any{
		raw.FechaEstimada,
		raw.FechaEstimadaCamel,
		raw.EstimatedDate,
		raw.EstimatedDateCamel,
	} {
		if parsed := parseEstimatedDate(ed); parsed != nil {
			c.EstimatedDate = parsed
			break
		}
	}

	// 7. Detalles
	for _, det := range []string{
		raw.Detalles,
		raw.Details,
		raw.Comentarios,
		raw.Mensaje,
	} {
		if trimmed := strings.TrimSpace(det); trimmed != "" {
			c.Details = trimmed
			break
		}
	}

	// Email / Mail / Correo
	for _, e := range []string{
		raw.EmailES,
		raw.Mail,
		raw.Correo,
		raw.CorreoElectronico,
		raw.CorreoElectronicoCamel,
		raw.EmailContacto,
		raw.EmailContactoCamel,
		raw.EmailEN,
		raw.EmailAddressCamel,
	} {
		if trimmed := strings.TrimSpace(e); trimmed != "" {
			c.Email = trimmed
			break
		}
	}

	return nil
}

func parseEstimatedDate(val any) *string {
	if val == nil {
		return nil
	}
	switch v := val.(type) {
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" || strings.EqualFold(trimmed, "null") || strings.EqualFold(trimmed, "undefined") {
			return nil
		}
		return &trimmed
	case fmt.Stringer:
		trimmed := strings.TrimSpace(v.String())
		if trimmed == "" || strings.EqualFold(trimmed, "null") {
			return nil
		}
		return &trimmed
	default:
		str := fmt.Sprintf("%v", v)
		trimmed := strings.TrimSpace(str)
		if trimmed == "" || strings.EqualFold(trimmed, "null") {
			return nil
		}
		return &trimmed
	}
}
