package handlers

import (
	"net/http"

	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/c-mierez/rss-aggregator/internal/lib/serve"
	"github.com/c-mierez/rss-aggregator/internal/middleware"
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
	// Get the user data from the request context
	authCTX := middleware.GetAuthCTX(r)

	serve.JSONResponse(w, http.StatusOK, authCTX.User)
}
