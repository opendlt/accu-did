package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/opendlt/accu-did/registrar-go/handlers"
	"github.com/opendlt/accu-did/registrar-go/internal/acc"
	"github.com/opendlt/accu-did/registrar-go/internal/policy"
)

func main() {
	// Get port from environment
	port := os.Getenv("REGISTRAR_PORT")
	if port == "" {
		port = "8082"
	}

	// Create Accumulate client (mock for now)
	accClient := acc.NewMockClient()

	// Create authorization policy
	authPolicy := policy.NewPolicyV1()

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(30 * time.Second))

	// Health check
	r.Get("/healthz", handlers.Healthz)

	// DID registration endpoints
	createHandler := handlers.NewCreateHandler(accClient, authPolicy)
	updateHandler := handlers.NewUpdateHandler(accClient, authPolicy)
	deactivateHandler := handlers.NewDeactivateHandler(accClient, authPolicy)

	r.Post("/create", createHandler.Create)
	r.Post("/update", updateHandler.Update)
	r.Post("/deactivate", deactivateHandler.Deactivate)

	// Create server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting registrar on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
