package main

import (
	"context"
	"log"
	"net/http"

	"github.com/c-mierez/rss-aggregator/internal/env"
	"github.com/c-mierez/rss-aggregator/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
)

func init() {
	// Load .env file
	env.LoadAndCheck(true)
}

func main() {

	// Connect to database
	db := setUpDatabase()
	defer db.Close(context.Background())

	// Create a new router
	router := setUpRouter()

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

func setUpDatabase() *pgx.Conn {
	db, err := pgx.Connect(context.Background(), env.Get(env.DATABASE_URL))
	if err != nil {
		log.Fatalf("Could not connect to database: %s\n", err.Error())
	}

	return db
}

func setUpRouter() *chi.Mux {
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

	return router
}
