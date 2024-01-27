package handlers

import (
	"net/http"

	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/c-mierez/rss-aggregator/internal/lib/serve"
	"github.com/c-mierez/rss-aggregator/internal/middleware"
	"github.com/c-mierez/rss-aggregator/internal/store"
)

type GetUserPostsHandler struct {
	db *queries.Queries
}

// New Handler
type NewGetUserPostsHandlerParams struct {
	DB *queries.Queries
}

func NewGetUserPostsHandler(params NewGetUserPostsHandlerParams) *GetUserPostsHandler {
	return &GetUserPostsHandler{
		db: params.DB,
	}
}

// ServeHTTP
func (h *GetUserPostsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Get the user data from the request context
	authCTX := middleware.GetAuthCTX(r)

	// Get the user posts
	posts, err := h.db.GetUserPosts(r.Context(), queries.GetUserPostsParams{
		UserID: authCTX.User.ID,
		Limit:  10,
	})
	if err != nil {
		serve.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return the posts
	serve.JSONResponse(w, http.StatusOK, store.DBToStorePosts(posts))

}
