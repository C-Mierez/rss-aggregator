package handlers

import (
	"net/http"

	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/c-mierez/rss-aggregator/internal/lib/serve"
	"github.com/c-mierez/rss-aggregator/internal/middleware"
	"github.com/c-mierez/rss-aggregator/internal/store"
)

type GetUserFollowsHandler struct {
	db *queries.Queries
}

// New Handler
type NewGetUserFollowsHandlerParams struct {
	DB *queries.Queries
}

func NewGetUserFollowsHandler(params NewGetUserFollowsHandlerParams) *GetUserFollowsHandler {
	return &GetUserFollowsHandler{
		db: params.DB,
	}
}

// ServeHTTP
func (h *GetUserFollowsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Get the user data from the request context
	authCTX := middleware.GetAuthCTX(r)

	// Fetch all follows from the database
	follows, err := h.db.GetFollowsByUserID(r.Context(), authCTX.User.ID)
	if err != nil {
		serve.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return the feeds
	serve.JSONResponse(w, http.StatusOK, store.DBToStoreFollows(follows))
}
