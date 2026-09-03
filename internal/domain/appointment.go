package domain

import (
	"encoding/json"
	"strings"
	"time"
)

type Appointment struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Date      string    `json:"date"`
	TimeSlot  string    `json:"time_slot"`
	Type      string    `json:"type"`
	Status    string    `json:"status"` // pending, confirmed, cancelled
	CreatedAt time.Time `json:"created_at"`
}

type CreateAppointmentInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Date     string `json:"date"`
	TimeSlot string `json:"time_slot"`
	Type     string `json:"type"`
}

// UnmarshalJSON allows accepting both Spanish and English JSON field names
func (c *CreateAppointmentInput) UnmarshalJSON(data []byte) error {
	var raw struct {
		// Spanish field names
		Nombre       string `json:"nombre"`
		EmailES      string `json:"email"`
		Telefono     string `json:"telefono"`
		Fecha        string `json:"fecha"`
		TramoHorario string `json:"tramo_horario"`
		TipoCita     string `json:"tipo_cita"`

		// English field names
		Name     string `json:"name"`
		EmailEN  string `json:"email_address"`
		Phone    string `json:"phone"`
		Date     string `json:"date"`
		TimeSlot string `json:"time_slot"`
		Type     string `json:"type"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Normalize values
	c.Name = strings.TrimSpace(raw.Nombre)
	if c.Name == "" {
		c.Name = strings.TrimSpace(raw.Name)
	}

	c.Email = strings.TrimSpace(raw.EmailES)
	if c.Email == "" {
		c.Email = strings.TrimSpace(raw.EmailEN)
	}

	c.Phone = strings.TrimSpace(raw.Telefono)
	if c.Phone == "" {
		c.Phone = strings.TrimSpace(raw.Phone)
	}

	c.Date = strings.TrimSpace(raw.Fecha)
	if c.Date == "" {
		c.Date = strings.TrimSpace(raw.Date)
	}

	c.TimeSlot = strings.TrimSpace(raw.TramoHorario)
	if c.TimeSlot == "" {
		c.TimeSlot = strings.TrimSpace(raw.TimeSlot)
	}

	c.Type = strings.TrimSpace(raw.TipoCita)
	if c.Type == "" {
		c.Type = strings.TrimSpace(raw.Type)
	}

	return nil
}
