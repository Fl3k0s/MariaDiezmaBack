package domain

import (
	"encoding/json"
	"strings"
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

type CreateDressInput struct {
	Name        string `json:"nombre"`
	Collection  string `json:"coleccion"`
	ImagePath   string `json:"ruta_imagen"`
	Image1Path  string `json:"ruta_imagen_1"`
	Image2Path  string `json:"ruta_imagen_2"`
	Image3Path  string `json:"ruta_imagen_3"`
	Description string `json:"descripcion"`
}

func (d *CreateDressInput) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name         string `json:"name"`
		Nombre       string `json:"nombre"`
		Collection   string `json:"collection"`
		Coleccion    string `json:"coleccion"`
		ImagePath    string `json:"image_path"`
		Imagen       string `json:"imagen"`
		RutaImagen   string `json:"ruta_imagen"`
		Image1Path   string `json:"image1_path"`
		Imagen1      string `json:"imagen_1"`
		RutaImagen1  string `json:"ruta_imagen_1"`
		Image2Path   string `json:"image2_path"`
		Imagen2      string `json:"imagen_2"`
		RutaImagen2  string `json:"ruta_imagen_2"`
		Image3Path   string `json:"image3_path"`
		Imagen3      string `json:"imagen_3"`
		RutaImagen3  string `json:"ruta_imagen_3"`
		Description  string `json:"description"`
		Descripcion  string `json:"descripcion"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	d.Name = strings.TrimSpace(raw.Nombre)
	if d.Name == "" {
		d.Name = strings.TrimSpace(raw.Name)
	}

	d.Collection = strings.TrimSpace(raw.Coleccion)
	if d.Collection == "" {
		d.Collection = strings.TrimSpace(raw.Collection)
	}

	d.ImagePath = strings.TrimSpace(raw.RutaImagen)
	if d.ImagePath == "" {
		d.ImagePath = strings.TrimSpace(raw.ImagePath)
	}
	if d.ImagePath == "" {
		d.ImagePath = strings.TrimSpace(raw.Imagen)
	}

	d.Image1Path = strings.TrimSpace(raw.RutaImagen1)
	if d.Image1Path == "" {
		d.Image1Path = strings.TrimSpace(raw.Image1Path)
	}
	if d.Image1Path == "" {
		d.Image1Path = strings.TrimSpace(raw.Imagen1)
	}

	d.Image2Path = strings.TrimSpace(raw.RutaImagen2)
	if d.Image2Path == "" {
		d.Image2Path = strings.TrimSpace(raw.Image2Path)
	}
	if d.Image2Path == "" {
		d.Image2Path = strings.TrimSpace(raw.Imagen2)
	}

	d.Image3Path = strings.TrimSpace(raw.RutaImagen3)
	if d.Image3Path == "" {
		d.Image3Path = strings.TrimSpace(raw.Image3Path)
	}
	if d.Image3Path == "" {
		d.Image3Path = strings.TrimSpace(raw.Imagen3)
	}

	d.Description = strings.TrimSpace(raw.Descripcion)
	if d.Description == "" {
		d.Description = strings.TrimSpace(raw.Description)
	}

	return nil
}

