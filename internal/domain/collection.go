package domain

import (
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
