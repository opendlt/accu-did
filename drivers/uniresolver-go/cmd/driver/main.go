package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/opendlt/accu-did/drivers/uniresolver-go/internal/proxy"
)

type Config struct {
	ResolverURL string `envconfig:"RESOLVER_URL" default:"http://resolver:8080"`
	Port        int    `envconfig:"PORT" default:"8081"`
}

func main() {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("Failed to process environment config: %v", err)
	}

	// Create proxy
	p := proxy.New(config.ResolverURL)

	// Setup routes
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", healthHandler).Methods("GET")

	// Universal Resolver Driver API
	r.HandleFunc("/1.0/identifiers/{did}", p.ResolveHandler).Methods("GET")

	// Driver metadata
	r.HandleFunc("/", driverInfoHandler).Methods("GET")

	// Start server
	addr := fmt.Sprintf(":%d", config.Port)
	log.Printf("Universal Resolver Driver starting on %s", addr)
	log.Printf("Proxying to resolver at %s", config.ResolverURL)

	server := &http.Server{
		Addr:         addr,
		Handler:      loggingMiddleware(r),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "uniresolver-driver-did-acc",
	})
}

func driverInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"driver": "did-acc",
		"version": "1.0.0",
		"methods": []string{"acc"},
		"endpoints": map[string]string{
			"resolve": "/1.0/identifiers/{did}",
		},
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Call next handler
		next.ServeHTTP(w, r)

		// Log response time
		log.Printf("Completed in %v", time.Since(start))
	})
}