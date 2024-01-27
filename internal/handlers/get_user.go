package handlers

import (
	"fmt"
	"net/http"

	"github.com/c-mierez/rss-aggregator/internal/lib/auth"
	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/c-mierez/rss-aggregator/internal/lib/serve"
	"github.com/c-mierez/rss-aggregator/internal/store"
)

type GetUserHandler struct {
	db *queries.Queries
}

// New Handler
type NewGetUserHandlerParams struct {
	DB *queries.Queries
}

func NewGetUserHandler(params NewGetUserHandlerParams) *GetUserHandler {
	return &GetUserHandler{
		db: params.DB,
	}
}

// ServeHTTP
type GetUserInput struct {
	ID string `json:"id" validate:"required"`
}

func (h *GetUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		serve.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := h.db.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		serve.JSONError(w, http.StatusBadRequest, fmt.Sprintf("Error fetching user: %v", err))
		return
	}

	serve.JSONResponse(w, http.StatusOK, store.DBToStoreUser(user))
}
