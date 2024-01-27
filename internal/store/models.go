package store

import (
	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/google/uuid"
)

/* ---------------------------------- User ---------------------------------- */
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	APIKey    string    `json:"api_key"`
}

func DBToStoreUser(dbUser queries.User) User {
	return User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		CreatedAt: dbUser.CreatedAt.String(),
		UpdatedAt: dbUser.UpdatedAt.String(),
		APIKey:    dbUser.ApiKey,
	}
}