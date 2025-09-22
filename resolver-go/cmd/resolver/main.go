package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/opendlt/accu-did/resolver-go/handlers"
	"github.com/opendlt/accu-did/resolver-go/internal/acc"
	"github.com/opendlt/accu-did/resolver-go/internal/resolve"
	"github.com/opendlt/accu-did/resolver-go/internal/security"
)

func main() {
	// Parse command line flags
	var (
		addr            = flag.String("addr", ":8080", "listen address")
		bind            = flag.String("bind", "127.0.0.1", "bind address (security: 127.0.0.1 for localhost only)")
		real            = flag.Bool("real", false, "enable real mode (connect to Accumulate network)")
		corsAllowOrigins = flag.String("cors-allow-origins", "", "comma-separated CORS allowed origins (empty=none, *=all)")
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

	// Parse CORS origins
	var corsOrigins []string
	if *corsAllowOrigins != "" {
		corsOrigins = strings.Split(*corsAllowOrigins, ",")
		for i, origin := range corsOrigins {
			corsOrigins[i] = strings.TrimSpace(origin)
		}
	}

	// Build full bind address
	fullAddr := *bind + *addr

	// Log startup configuration
	log.Printf("Starting DID Resolver")
	log.Printf("  Mode: %s", mode)
	log.Printf("  Bind: %s", fullAddr)
	log.Printf("  CORS Origins: %v", corsOrigins)
	if *real && nodeURL != "" {
		log.Printf("  Accumulate Node: %s", nodeURL)
	}

	// Create Accumulate client
	accClient := acc.NewClient(*real, nodeURL)

	// Setup router
	r := chi.NewRouter()

	// Security middleware
	r.Use(security.RequestIDMiddleware())
	r.Use(security.CORSMiddleware(corsOrigins))

	// Standard middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Health check
	r.Get("/healthz", handlers.Healthz)

	// DID resolution
	resolveHandler := resolve.NewHandler(accClient)
	r.Get("/resolve", resolveHandler.Resolve)

	// Universal Resolver 1.0 compatibility
	r.Get("/1.0/identifiers/{did}", resolveHandler.UniversalResolve)

	// Create server
	srv := &http.Server{
		Addr:         fullAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on %s", fullAddr)
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
