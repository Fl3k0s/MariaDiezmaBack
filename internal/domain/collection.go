package domain

import (
	"encoding/json"
	"strings"
	"time"
)

type Collection struct {
	ID          string    `json:"id"`
	Name        string    `json:"nombre"`
	ImagePath   string    `json:"imagen"`
	Description string    `json:"descripcion"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CollectionResponse struct {
	ID          string `json:"id"`
	Name        string `json:"nombre"`
	ImagePath   string `json:"imagen"`
	Description string `json:"descripcion"`
}

func (c *Collection) ToResponse() CollectionResponse {
	return CollectionResponse{
		ID:          c.ID,
		Name:        c.Name,
		ImagePath:   c.ImagePath,
		Description: c.Description,
	}
}

type CreateCollectionInput struct {
	Name        string `json:"nombre"`
	ImagePath   string `json:"imagen"`
	Description string `json:"descripcion"`
}

func (c *CreateCollectionInput) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name        string `json:"name"`
		Nombre      string `json:"nombre"`
		ImagePath   string `json:"image_path"`
		Imagen      string `json:"imagen"`
		Image       string `json:"image"`
		RutaImagen  string `json:"ruta_imagen"`
		Description string `json:"description"`
		Descripcion string `json:"descripcion"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	c.Name = strings.TrimSpace(raw.Nombre)
	if c.Name == "" {
		c.Name = strings.TrimSpace(raw.Name)
	}

	c.ImagePath = strings.TrimSpace(raw.Imagen)
	if c.ImagePath == "" {
		c.ImagePath = strings.TrimSpace(raw.ImagePath)
	}
	if c.ImagePath == "" {
		c.ImagePath = strings.TrimSpace(raw.Image)
	}
	if c.ImagePath == "" {
		c.ImagePath = strings.TrimSpace(raw.RutaImagen)
	}

	c.Description = strings.TrimSpace(raw.Descripcion)
	if c.Description == "" {
		c.Description = strings.TrimSpace(raw.Description)
	}

	return nil
}

