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

type CreateFollowHandler struct {
	db *queries.Queries
}

// New Handler
type NewCreateFollowHandlerParams struct {
	DB *queries.Queries
}

func NewCreateFollowHandler(params NewCreateFollowHandlerParams) *CreateFollowHandler {
	return &CreateFollowHandler{
		db: params.DB,
	}
}

// ServeHTTP
type CreateFollowInput struct {
	FeedID string `json:"feed_id" validate:"required,uuid"`
}

func (h *CreateFollowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Get the user data from the request context
	authCTX := middleware.GetAuthCTX(r)

	// Decode the request body
	input, err := serve.JSONValidRequest[CreateFollowInput](w, r)
	if err != nil {
		return
	}

	// Verify feed exists
	// Validator already makes sure the feed_id is a valid uuid
	feedUUID := uuid.MustParse(input.FeedID)

	_, err = h.db.GetFeed(r.Context(), feedUUID)
	if err != nil {
		serve.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create a new follow
	follow, err := h.db.CreateFollow(r.Context(), queries.CreateFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    authCTX.User.ID,
		FeedID:    feedUUID,
	})
	if err != nil {
		serve.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return the follow
	serve.JSONResponse(w, http.StatusCreated, store.DBToStoreFollow(follow))

}
