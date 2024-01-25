package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/c-mierez/rss-aggregator/internal/env"
	"github.com/c-mierez/rss-aggregator/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func init() {
	// Load .env file
	env.LoadAndCheck(true)
}

func main() {
	// Create a global context
	globalCtx, globalCtxCancel := context.WithCancel(context.Background())

	// Connect to database
	db := setUpDatabase(globalCtx)

	// Create a new router
	router := setUpRouter()

	// Start the server
	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:" + env.Get(env.PORT),
	}

	gracefulShutdownWrapper(globalCtx, func() {
		log.Printf("Starting server on PORT: %+v\n", env.Get(env.PORT))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s\n", err.Error())
		}
	}, func(shutdownCtx context.Context) {
		// Shutdown the server
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Could not gracefully shutdown server: %s\n", err.Error())
		}

		// Close the database connection
		db.Close()

		log.Println("Graceful shutdown complete.")

		// Cancel the global context
		globalCtxCancel()
	})

	<-globalCtx.Done() // Wait for the global context to be cancelled

}

func gracefulShutdownWrapper(globalCtx context.Context, gracefulExecution func(), gracefulShutdown func(shutdownCtx context.Context)) {

	// Graceful shutdown signal channel
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Routine to handle graceful shutdown
	go func() {
		sig := <-shutdownChan
		log.Printf("Received signal %s. Shutting down...\n", sig.String())

		// Create a failsafe trigger to shut down after 10 seconds if the graceful shutdown does not finish in time
		failsafeCtx, failsafeCtxCancel := context.WithTimeout(globalCtx, 10*time.Second)
		defer failsafeCtxCancel()

		go func() {
			<-failsafeCtx.Done() // Wait for the failsafe context to be cancelled after it times out

			if failsafeCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out. Shutting down immediately...")
			}
		}()

		// Perform all the graceful shutdown tasks
		gracefulShutdown(failsafeCtx)
	}()

	// Routine in which the graceful execution will take place
	go gracefulExecution()
}

func setUpDatabase(ctx context.Context) *pgxpool.Pool {
	db, err := pgxpool.New(ctx, env.Get(env.DATABASE_URL))
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
