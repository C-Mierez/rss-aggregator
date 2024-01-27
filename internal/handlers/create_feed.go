package handlers

import (
	"net/http"
	"time"

	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/c-mierez/rss-aggregator/internal/lib/serve"
	"github.com/c-mierez/rss-aggregator/internal/middleware"
	"github.com/c-mierez/rss-aggregator/internal/store"
	"github.com/google/uuid"
)

type CreateFeedHandler struct {
	db *queries.Queries
}

// New Handler
type NewCreateFeedHandlerParams struct {
	DB *queries.Queries
}

func NewCreateFeedHandler(params NewCreateFeedHandlerParams) *CreateFeedHandler {
	return &CreateFeedHandler{
		db: params.DB,
	}
}

// ServeHTTP
type CreateFeedInput struct {
	Name string `json:"name" validate:"required"`
	Url  string `json:"url" validate:"required"`
}

func (h *CreateFeedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the user data from the request context
	authCTX := middleware.GetAuthCTX(r)

	// Decode the request body
	input, err := serve.JSONValidRequest[CreateFeedInput](w, r)
	if err != nil {
		return
	}

	feed, err := h.db.CreateFeed(r.Context(), queries.CreateFeedParams{
		ID:        uuid.New(),
		Name:      input.Name,
		Url:       input.Url,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    authCTX.User.ID,
	})

	if err != nil {
		serve.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	serve.JSONResponse(w, http.StatusOK, store.DBToStoreFeed(feed))
}
