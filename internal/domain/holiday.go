package domain

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Holiday represents an official holiday for Madrid defined year by year.
type Holiday struct {
	ID          string    `json:"id"`
	Date        string    `json:"date"`        // YYYY-MM-DD
	Year        int       `json:"year"`        // e.g. 2026
	Name        string    `json:"name"`        // e.g. "San Isidro"
	Description string    `json:"description"` // e.g. "Patrón de Madrid"
	Scope       string    `json:"scope"`       // e.g. "local", "autonomico", "nacional"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// HolidayResponse is the structure serialized for frontend and web clients.
// It provides both English and Spanish field keys for seamless web integration.
type HolidayResponse struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	Fecha       string `json:"fecha"`
	Year        int    `json:"year"`
	Anio        int    `json:"anio"`
	Name        string `json:"name"`
	Nombre      string `json:"nombre"`
	Description string `json:"description,omitempty"`
	Descripcion string `json:"descripcion,omitempty"`
	Scope       string `json:"scope,omitempty"`
	Ambito      string `json:"ambito,omitempty"`
}

func (h *Holiday) ToResponse() HolidayResponse {
	return HolidayResponse{
		ID:          h.ID,
		Date:        h.Date,
		Fecha:       h.Date,
		Year:        h.Year,
		Anio:        h.Year,
		Name:        h.Name,
		Nombre:      h.Name,
		Description: h.Description,
		Descripcion: h.Description,
		Scope:       h.Scope,
		Ambito:      h.Scope,
	}
}

// CreateHolidayInput supports payloads in Spanish and English (snake_case and camelCase).
type CreateHolidayInput struct {
	Date        string `json:"date"`
	Year        int    `json:"year"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Scope       string `json:"scope"`
}

func (c *CreateHolidayInput) UnmarshalJSON(data []byte) error {
	var raw struct {
		Date             string `json:"date"`
		Fecha            string `json:"fecha"`
		Year             any    `json:"year"`
		Anio             any    `json:"anio"`
		Ano              any    `json:"ano"`
		Name             string `json:"name"`
		Nombre           string `json:"nombre"`
		Description      string `json:"description"`
		Descripcion      string `json:"descripcion"`
		Scope            string `json:"scope"`
		Ambito           string `json:"ambito"`
		Tipo             string `json:"tipo"`
		Festivo          string `json:"festivo"`
		NombreFestivo    string `json:"nombre_festivo"`
		NombreFestivoCam string `json:"nombreFestivo"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Date / Fecha
	c.Date = strings.TrimSpace(raw.Date)
	if c.Date == "" {
		c.Date = strings.TrimSpace(raw.Fecha)
	}

	// Name / Nombre
	for _, n := range []string{raw.Name, raw.Nombre, raw.Festivo, raw.NombreFestivo, raw.NombreFestivoCam} {
		if t := strings.TrimSpace(n); t != "" {
			c.Name = t
			break
		}
	}

	// Description / Descripcion
	c.Description = strings.TrimSpace(raw.Description)
	if c.Description == "" {
		c.Description = strings.TrimSpace(raw.Descripcion)
	}

	// Scope / Ambito
	c.Scope = strings.TrimSpace(raw.Scope)
	if c.Scope == "" {
		c.Scope = strings.TrimSpace(raw.Ambito)
	}
	if c.Scope == "" {
		c.Scope = strings.TrimSpace(raw.Tipo)
	}
	if c.Scope == "" {
		c.Scope = "madrid"
	}

	// Year parsing from any type (int, float64, string)
	parseYear := func(v any) int {
		if v == nil {
			return 0
		}
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
		case string:
			clean := strings.TrimSpace(val)
			if y, err := strconv.Atoi(clean); err == nil {
				return y
			}
		}
		return 0
	}

	c.Year = parseYear(raw.Year)
	if c.Year == 0 {
		c.Year = parseYear(raw.Anio)
	}
	if c.Year == 0 {
		c.Year = parseYear(raw.Ano)
	}

	// Infer year from date if date is formatted like "YYYY-MM-DD"
	if c.Year == 0 && len(c.Date) >= 4 {
		if y, err := strconv.Atoi(c.Date[:4]); err == nil {
			c.Year = y
		}
	}

	return nil
}

// Validate validates input fields for a holiday.
func (c *CreateHolidayInput) Validate() error {
	if strings.TrimSpace(c.Date) == "" {
		return fmt.Errorf("%w: la fecha del festivo es obligatoria (formato YYYY-MM-DD)", ErrInvalidInput)
	}

	// Ensure valid date format YYYY-MM-DD
	parsedTime, err := time.Parse("2006-01-02", c.Date)
	if err != nil {
		return fmt.Errorf("%w: el formato de fecha debe ser YYYY-MM-DD (ejemplo: 2026-05-15)", ErrInvalidInput)
	}

	if c.Year == 0 {
		c.Year = parsedTime.Year()
	} else if c.Year != parsedTime.Year() {
		// Keep year synchronized with date
		c.Year = parsedTime.Year()
	}

	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("%w: el nombre del festivo es obligatorio", ErrInvalidInput)
	}

	return nil
}

// UpdateHolidayInput provides fields for updating a holiday.
type UpdateHolidayInput struct {
	Date        string `json:"date"`
	Year        int    `json:"year"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Scope       string `json:"scope"`
}

func (u *UpdateHolidayInput) UnmarshalJSON(data []byte) error {
	var c CreateHolidayInput
	if err := c.UnmarshalJSON(data); err != nil {
		return err
	}
	u.Date = c.Date
	u.Year = c.Year
	u.Name = c.Name
	u.Description = c.Description
	u.Scope = c.Scope
	return nil
}

// DefaultMadridHolidays returns standard official holidays for Madrid (years 2025 and 2026).
func DefaultMadridHolidays() []Holiday {
	now := time.Now().UTC()
	return []Holiday{
		// --- AÑO 2025 ---
		{ID: "50000000-0000-0000-2025-000000000001", Date: "2025-01-01", Year: 2025, Name: "Año Nuevo", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000002", Date: "2025-01-06", Year: 2025, Name: "Epifanía del Señor (Reyes Magos)", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000003", Date: "2025-04-17", Year: 2025, Name: "Jueves Santo", Description: "Festivo autonómico Comunidad de Madrid", Scope: "autonomico", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000004", Date: "2025-04-18", Year: 2025, Name: "Viernes Santo", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000005", Date: "2025-05-01", Year: 2025, Name: "Fiesta del Trabajo", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000006", Date: "2025-05-02", Year: 2025, Name: "Fiesta de la Comunidad de Madrid", Description: "Día de la Comunidad de Madrid", Scope: "autonomico", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000007", Date: "2025-05-15", Year: 2025, Name: "San Isidro Labrador", Description: "Patrón de Madrid capital", Scope: "local", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000008", Date: "2025-07-25", Year: 2025, Name: "Santiago Apóstol", Description: "Festivo autonómico Comunidad de Madrid", Scope: "autonomico", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000009", Date: "2025-08-15", Year: 2025, Name: "Asunción de la Virgen", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000010", Date: "2025-10-12", Year: 2025, Name: "Fiesta Nacional de España", Description: "Día de la Hispanidad", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000011", Date: "2025-11-01", Year: 2025, Name: "Todos los Santos", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000012", Date: "2025-11-10", Year: 2025, Name: "Lunes siguiente a La Almudena", Description: "Patrona de Madrid capital (traslado)", Scope: "local", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000013", Date: "2025-12-06", Year: 2025, Name: "Día de la Constitución Española", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000014", Date: "2025-12-08", Year: 2025, Name: "Inmaculada Concepción", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2025-000000000015", Date: "2025-12-25", Year: 2025, Name: "Natividad del Señor (Navidad)", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},

		// --- AÑO 2026 ---
		{ID: "50000000-0000-0000-2026-000000000001", Date: "2026-01-01", Year: 2026, Name: "Año Nuevo", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2026-000000000002", Date: "2026-01-06", Year: 2026, Name: "Epifanía del Señor (Reyes Magos)", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2026-000000000003", Date: "2026-04-02", Year: 2026, Name: "Jueves Santo", Description: "Festivo autonómico Comunidad de Madrid", Scope: "autonomico", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2026-000000000004", Date: "2026-04-03", Year: 2026, Name: "Viernes Santo", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2026-000000000005", Date: "2026-05-01", Year: 2026, Name: "Fiesta del Trabajo", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2026-000000000006", Date: "2026-05-02", Year: 2026, Name: "Fiesta de la Comunidad de Madrid", Description: "Día de la Comunidad de Madrid", Scope: "autonomico", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2026-000000000007", Date: "2026-05-15", Year: 2026, Name: "San Isidro Labrador", Description: "Patrón de Madrid capital", Scope: "local", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2026-000000000008", Date: "2026-08-15", Year: 2026, Name: "Asunción de la Virgen", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2026-000000000009", Date: "2026-10-12", Year: 2026, Name: "Fiesta Nacional de España", Description: "Día de la Hispanidad", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2026-000000000010", Date: "2026-11-02", Year: 2026, Name: "Lunes siguiente a Todos los Santos", Description: "Festivo nacional/autonómico (traslado)", Scope: "autonomico", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2026-000000000011", Date: "2026-11-09", Year: 2026, Name: "Nuestra Señora de la Almudena", Description: "Patrona de Madrid capital", Scope: "local", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2026-000000000012", Date: "2026-12-07", Year: 2026, Name: "Lunes siguiente a la Constitución Española", Description: "Festivo autonómico Comunidad de Madrid (traslado)", Scope: "autonomico", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2026-000000000013", Date: "2026-12-08", Year: 2026, Name: "Inmaculada Concepción", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
		{ID: "50000000-0000-0000-2026-000000000014", Date: "2026-12-25", Year: 2026, Name: "Natividad del Señor (Navidad)", Description: "Festivo nacional", Scope: "nacional", CreatedAt: now, UpdatedAt: now},
	}
}
