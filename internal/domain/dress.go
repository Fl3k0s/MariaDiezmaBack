package domain

import (
	"time"
)

type Dress struct {
	ID          string    `json:"id"`
	Name        string    `json:"nombre"`
	Collection  string    `json:"coleccion"`
	ImagePath   string    `json:"ruta_imagen"`   // Portada / Imagen 1
	Image1Path  string    `json:"ruta_imagen_1"` // Imagen 1
	Image2Path  string    `json:"ruta_imagen_2"` // Imagen 2
	Image3Path  string    `json:"ruta_imagen_3"` // Imagen 3
	Description string    `json:"descripcion"`   // Descripción del vestido
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DressResponse struct {
	ID         string `json:"id"`
	Name       string `json:"nombre"`
	Collection string `json:"coleccion"`
	ImagePath  string `json:"ruta_imagen"` // Imagen 1 para catálogo
}

type DressDetailResponse struct {
	ID          string `json:"id"`
	Name        string `json:"nombre"`
	Collection  string `json:"coleccion"`
	Image1Path  string `json:"ruta_imagen_1"`
	Image2Path  string `json:"ruta_imagen_2"`
	Image3Path  string `json:"ruta_imagen_3"`
	Description string `json:"descripcion"`
}

func (d *Dress) ToResponse() DressResponse {
	img := d.ImagePath
	if img == "" {
		img = d.Image1Path
	}
	return DressResponse{
		ID:         d.ID,
		Name:       d.Name,
		Collection: d.Collection,
		ImagePath:  img,
	}
}

func (d *Dress) ToDetailResponse() DressDetailResponse {
	img1 := d.Image1Path
	if img1 == "" {
		img1 = d.ImagePath
	}
	img3 := d.Image3Path
	if img3 == "" {
		img3 = d.ImagePath
	}
	return DressDetailResponse{
		ID:          d.ID,
		Name:        d.Name,
		Collection:  d.Collection,
		Image1Path:  img1,
		Image2Path:  d.Image2Path,
		Image3Path:  img3,
		Description: d.Description,
	}
}
