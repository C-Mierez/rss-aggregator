package main

import (
	"log"
	"net/http"

	"github.com/c-mierez/rss-aggregator/internal/env"
	"github.com/c-mierez/rss-aggregator/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func init() {
	// Load .env file
	env.LoadAndCheck(true)
}

func main() {

	// Create a new router
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Routes
	router.Get("/health", handlers.NewHealthHandler().ServeHTTP)

	router.Get("/error", handlers.NewErrorHandler().ServeHTTP)

	// Start the server
	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:" + env.Get(env.PORT),
	}

	log.Printf("Starting server on PORT: %+v\n", env.Get(env.PORT))
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

}
