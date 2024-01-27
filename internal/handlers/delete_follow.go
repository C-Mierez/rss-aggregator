package handlers

import (
	"net/http"

	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/c-mierez/rss-aggregator/internal/lib/serve"
	"github.com/c-mierez/rss-aggregator/internal/middleware"
	"github.com/google/uuid"
)

type DeleteFollowHandler struct {
	db *queries.Queries
}

// New Handler
type NewDeleteFollowHandlerParams struct {
	DB *queries.Queries
}

func NewDeleteFollowHandler(params NewDeleteFollowHandlerParams) *DeleteFollowHandler {
	return &DeleteFollowHandler{
		db: params.DB,
	}
}

// ServeHTTP
type DeleteFollowInput struct {
	FollowID string `json:"follow_id" validate:"required,uuid"`
}

func (h *DeleteFollowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Get the user data from the request context
	authCTX := middleware.GetAuthCTX(r)

	// Decode the request body
	input, err := serve.JSONValidRequest[DeleteFollowInput](w, r)
	if err != nil {
		return
	}

	// Validator already makes sure the follow_id is a valid uuid
	followUUID := uuid.MustParse(input.FollowID)

	// Delete the follow
	err = h.db.DeleteUserFollowByID(r.Context(), queries.DeleteUserFollowByIDParams{
		ID:     followUUID,
		UserID: authCTX.User.ID,
	})

	if err != nil {
		serve.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	serve.JSONResponse(w, http.StatusOK, nil)

}
