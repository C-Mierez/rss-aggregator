package handlers

import (
	"net/http"

	"github.com/c-mierez/rss-aggregator/internal/lib/serve"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serve.JSONResponse(w, http.StatusOK, "OK")
}
