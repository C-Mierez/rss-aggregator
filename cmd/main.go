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
	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/c-mierez/rss-aggregator/internal/middleware"
	"github.com/c-mierez/rss-aggregator/internal/scrapper"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"database/sql"

	_ "github.com/lib/pq"
)

func init() {
	// Load .env file
	env.LoadAndCheck(true)
}

func main() {

	// Create a global context
	globalCtx, globalCtxCancel := context.WithCancel(context.Background())

	// Connect to database
	db, err := sql.Open("postgres", env.Get(env.DATABASE_URL))
	if err != nil {
		log.Fatalf("Could not connect to database: %s\n", err.Error())
	}

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

	q := queries.New(db)

	// Routes
	router.Get("/health", handlers.NewHealthHandler().ServeHTTP)
	router.Get("/error", handlers.NewErrorHandler().ServeHTTP)
	router.Post("/createUser", handlers.NewCreateUserHandler(handlers.NewCreateUserHandlerParams{
		DB: q,
	}).ServeHTTP)
	router.Get("/getUser", middleware.AuthMiddleware(handlers.NewGetUserHandler(handlers.NewGetUserHandlerParams{
		DB: q,
	}), q))
	router.Post("/createFeed", middleware.AuthMiddleware(handlers.NewCreateFeedHandler(handlers.NewCreateFeedHandlerParams{
		DB: q,
	}), q))
	router.Get("/getFeeds", handlers.NewGetFeedsHandler(handlers.NewGetFeedsHandlerParams{
		DB: q,
	}).ServeHTTP)
	router.Post("/createFollow", middleware.AuthMiddleware(handlers.NewCreateFollowHandler(handlers.NewCreateFollowHandlerParams{
		DB: q,
	}), q))
	router.Get("/getUserFollows", middleware.AuthMiddleware(handlers.NewGetUserFollowsHandler(handlers.NewGetUserFollowsHandlerParams{
		DB: q,
	}), q))
	router.Delete("/deleteFollow", middleware.AuthMiddleware(handlers.NewDeleteFollowHandler(handlers.NewDeleteFollowHandlerParams{
		DB: q,
	}), q))
	router.Get("/getUserPosts", middleware.AuthMiddleware(handlers.NewGetUserPostsHandler(handlers.NewGetUserPostsHandlerParams{
		DB: q,
	}), q))

	// Start the server
	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:" + env.Get(env.PORT),
	}

	graceful(
		globalCtx,
		func() {
			// Start Scraping Routine
			go scrapper.StartScraping(q, 10, 10*time.Second)

			// Server HTTP Server
			log.Printf("Starting server on PORT: %+v\n", env.Get(env.PORT))
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Could not start server: %s\n", err.Error())
			}
		},
		func(shutdownCtx context.Context) {
			// Shutdown the server
			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Fatalf("Could not gracefully shutdown server: %s\n", err.Error())
			}

			// Close the database connection
			db.Close()

			log.Println("Graceful shutdown complete.")

			// Cancel the global context
			globalCtxCancel()
		},
	)

	<-globalCtx.Done() // Wait for the global context to be cancelled

}

// Wrapper for handling graceful execution and shutdown
func graceful(globalCtx context.Context, gracefulExecution func(), gracefulShutdown func(shutdownCtx context.Context)) {

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
