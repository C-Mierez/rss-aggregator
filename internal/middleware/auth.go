package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/c-mierez/rss-aggregator/internal/lib/auth"
	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/c-mierez/rss-aggregator/internal/lib/serve"
	"github.com/c-mierez/rss-aggregator/internal/store"
)

type AuthMiddlewareCtx struct {
	User store.User
}

func GetAuthCTX(r *http.Request) AuthMiddlewareCtx {
	data := r.Context().Value(AuthMiddlewareCTXKey)
	if data == nil {
		// This is a panic because it means it is being used without the middleware
		log.Fatalf("AuthMiddlewareCtx not found in request context")
	}

	authCTX, ok := data.(AuthMiddlewareCtx)
	if !ok {
		// This is a panic because it means the data is not of the correct type
		log.Fatalf("AuthMiddlewareCtx not of correct type")
	}

	return authCTX
}

func AuthMiddleware(handler http.Handler, db *queries.Queries) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Decode the API key from the request header
		apiKey, err := auth.GetAPIKeyHeader(r.Header)
		if err != nil {
			serve.JSONError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Fetch the user from the database using the API key
		user, err := db.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			serve.JSONError(w, http.StatusBadRequest, fmt.Sprintf("Error fetching user: %v", err))
			return
		}

		// Add the user to the request context
		ctx := context.WithValue(r.Context(), AuthMiddlewareCTXKey, AuthMiddlewareCtx{
			User: store.DBToStoreUser(user),
		})

		// Call the next handler
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
