package domain

// PressArticle represents a press article feature or publication.
type PressArticle struct {
	ID              string `json:"id,omitempty"`
	MagazineName    string `json:"nombre_revista"`
	PublicationDate string `json:"fecha_publicacion"`
	Title           string `json:"titular"`
	Description     string `json:"descripcion"`
	ArticleURL      string `json:"enlace_articulo"`
	URL             string `json:"enlace,omitempty"`
}

// PressArticleResponse defines the response structure for the API.
type PressArticleResponse struct {
	ID                 string `json:"id,omitempty"`
	MagazineName       string `json:"nombre_revista"`
	PublicationDate    string `json:"fecha_publicacion"`
	Title              string `json:"titular"`
	Description        string `json:"descripcion"`
	ShortDescription   string `json:"pequena_descripcion,omitempty"`
	ArticleURL         string `json:"enlace_articulo"`
	Link               string `json:"enlace"`
}

func (a *PressArticle) ToResponse() PressArticleResponse {
	url := a.ArticleURL
	if url == "" {
		url = a.URL
	}
	return PressArticleResponse{
		ID:               a.ID,
		MagazineName:     a.MagazineName,
		PublicationDate:  a.PublicationDate,
		Title:            a.Title,
		Description:      a.Description,
		ShortDescription: a.Description,
		ArticleURL:       url,
		Link:             url,
	}
}
