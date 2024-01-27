package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/c-mierez/rss-aggregator/internal/lib/serve"
	"github.com/c-mierez/rss-aggregator/internal/store"
	"github.com/google/uuid"
)

type CreateUserHandler struct {
	db *queries.Queries
}

// New Handler
type NewCreateUserHandlerParams struct {
	DB *queries.Queries
}

func NewCreateUserHandler(params NewCreateUserHandlerParams) *CreateUserHandler {
	return &CreateUserHandler{
		db: params.DB,
	}
}

// ServeHTTP
type CreateUserInput struct {
	Name string `json:"name" validate:"required"`
}

func (h *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	data, err := serve.JSONValidRequest[CreateUserInput](w, r)
	if err != nil {
		return
	}

	user, err := h.db.CreateUser(r.Context(), queries.CreateUserParams{
		ID:        uuid.New(),
		Name:      data.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("Error creating user: %+v\n", err)
		serve.JSONError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	serve.JSONResponse(w, http.StatusOK, store.DBToStoreUser(user))

}
