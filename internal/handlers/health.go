package handlers

import (
	"net/http"

	"github.com/c-mierez/rss-aggregator/internal/lib/res"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res.JSONResponse(w, http.StatusOK, "OK")
}
