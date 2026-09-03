package handler

import (
	"net/http"
	"time"

	"mariadiezmaback/pkg/response"
)

type HealthHandler struct {
	startTime time.Time
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{startTime: time.Now()}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.OK(w, "service is healthy", map[string]any{
		"status": "UP",
		"uptime": time.Since(h.startTime).String(),
		"time":   time.Now().UTC(),
	})
}
