package main

import (
	"context"
	"flag"
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
	// Parse command line flags
	var (
		addr = flag.String("addr", ":8081", "listen address")
		real = flag.Bool("real", false, "enable real mode (connect to Accumulate network)")
	)
	flag.Parse()

	// Get Accumulate node URL from environment if in real mode
	var nodeURL string
	if *real {
		nodeURL = os.Getenv("ACC_NODE_URL")
		if nodeURL == "" {
			log.Fatal("ACC_NODE_URL environment variable is required when using --real mode")
		}
	}

	// Determine mode for logging
	mode := "FAKE"
	if *real {
		mode = "REAL"
	}

	// Create Accumulate submitter
	accSubmitter := acc.NewSubmitter(*real, nodeURL)

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

	// Legacy DID registration endpoints (Universal Registrar v0.x compatibility)
	createHandler := handlers.NewCreateHandler(accSubmitter, authPolicy)
	updateHandler := handlers.NewUpdateHandler(accSubmitter, authPolicy)
	deactivateHandler := handlers.NewDeactivateHandler(accSubmitter, authPolicy)

	r.Post("/create", createHandler.Create)
	r.Post("/update", updateHandler.Update)
	r.Post("/deactivate", deactivateHandler.Deactivate)

	// Native DID registration endpoints (clean internal API)
	nativeHandler := handlers.NewNativeHandler(accSubmitter)
	r.Post("/register", nativeHandler.Register)
	r.Post("/native/update", nativeHandler.Update)
	r.Post("/native/deactivate", nativeHandler.Deactivate)

	// Universal Registrar v1.0 compatibility endpoints
	universalHandler := handlers.NewUniversalHandler(accSubmitter, authPolicy)
	r.Post("/1.0/create", universalHandler.UniversalCreate)
	r.Post("/1.0/update", universalHandler.UniversalUpdate)
	r.Post("/1.0/deactivate", universalHandler.UniversalDeactivate)

	// Create server
	srv := &http.Server{
		Addr:         *addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting DID registrar on %s (mode: %s)", *addr, mode)
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
