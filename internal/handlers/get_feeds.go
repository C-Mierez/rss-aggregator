package handlers

import (
	"net/http"

	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/c-mierez/rss-aggregator/internal/lib/serve"
	"github.com/c-mierez/rss-aggregator/internal/store"
)

type GetFeedsHandler struct {
	db *queries.Queries
}

// New Handler
type NewGetFeedsHandlerParams struct {
	DB *queries.Queries
}

func NewGetFeedsHandler(params NewGetFeedsHandlerParams) *GetFeedsHandler {
	return &GetFeedsHandler{
		db: params.DB,
	}
}

// ServeHTTP
func (h *GetFeedsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Fetch all feeds from the database
	feeds, err := h.db.GetFeeds(r.Context())
	if err != nil {
		serve.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return the feeds
	serve.JSONResponse(w, http.StatusOK, store.DBToStoreFeeds(feeds))
}
