package handlers

import (
	"net/http"

	"github.com/c-mierez/rss-aggregator/internal/lib/res"
)

type ErrorHandler struct{}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (h *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res.JSONError(w, http.StatusInternalServerError, "Something went wrong")
}
