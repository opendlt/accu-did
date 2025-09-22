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

	"github.com/opendlt/accu-did/registrar-go/handlers"
	"github.com/opendlt/accu-did/registrar-go/internal/acc"
	"github.com/opendlt/accu-did/registrar-go/internal/policy"
	"github.com/opendlt/accu-did/registrar-go/internal/security"
)

func main() {
	// Parse command line flags
	var (
		addr        = flag.String("addr", ":8081", "listen address")
		bind        = flag.String("bind", "127.0.0.1", "bind address (security: 127.0.0.1 for localhost only)")
		real        = flag.Bool("real", false, "enable real mode (connect to Accumulate network)")
		authAPIKey  = flag.String("auth-api-key", "", "API key for authentication (env: REGISTRAR_API_KEY)")
		allowlist   = flag.String("allowlist", "", "comma-separated CIDR/IP allowlist (env: REGISTRAR_ALLOWLIST)")
		rateRPS     = flag.Int("rate-rps", 50, "rate limit requests per second")
		rateBurst   = flag.Int("rate-burst", 100, "rate limit burst capacity")
	)
	flag.Parse()

	// Override from environment variables
	if envAPIKey := os.Getenv("REGISTRAR_API_KEY"); envAPIKey != "" {
		*authAPIKey = envAPIKey
	}
	if envAllowlist := os.Getenv("REGISTRAR_ALLOWLIST"); envAllowlist != "" {
		*allowlist = envAllowlist
	}

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

	// Parse allowlist
	var allowList []string
	if *allowlist != "" {
		allowList = strings.Split(*allowlist, ",")
		for i, item := range allowList {
			allowList[i] = strings.TrimSpace(item)
		}
	}

	// Build full bind address
	fullAddr := *bind + *addr

	// Security configuration logged above
	_ = &security.SecurityConfig{
		APIKey:    *authAPIKey,
		AllowList: allowList,
		RateRPS:   *rateRPS,
		RateBurst: *rateBurst,
		BindAddr:  fullAddr,
	}

	// Log startup configuration
	log.Printf("Starting DID Registrar")
	log.Printf("  Mode: %s", mode)
	log.Printf("  Bind: %s", fullAddr)
	log.Printf("  API Key: %s", func() string {
		if *authAPIKey != "" {
			return "configured"
		}
		return "disabled"
	}())
	log.Printf("  IP Allowlist: %v", allowList)
	log.Printf("  Rate Limit: %d RPS, %d burst", *rateRPS, *rateBurst)
	if *real && nodeURL != "" {
		log.Printf("  Accumulate Node: %s", nodeURL)
	}

	// Create Accumulate submitter
	accSubmitter := acc.NewSubmitter(*real, nodeURL)

	// Create authorization policy
	authPolicy := policy.NewPolicyV1()

	// Setup router
	r := chi.NewRouter()

	// Security middleware (applied to all routes)
	r.Use(security.RequestIDMiddleware())
	r.Use(security.IPAllowListMiddleware(allowList))
	r.Use(security.RateLimitMiddleware(*rateRPS, *rateBurst))
	r.Use(security.APIKeyMiddleware(*authAPIKey))

	// Standard middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
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
