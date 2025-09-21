package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/opendlt/accu-did/drivers/uniregistrar-go/internal/proxy"
)

type Config struct {
	RegistrarURL string `envconfig:"REGISTRAR_URL" default:"http://registrar:8082"`
	Port         int    `envconfig:"PORT" default:"8083"`
}

func main() {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("Failed to process environment config: %v", err)
	}

	// Create proxy
	p := proxy.New(config.RegistrarURL)

	// Setup routes
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", healthHandler).Methods("GET")

	// Universal Registrar Driver API
	r.HandleFunc("/1.0/create", p.CreateHandler).Methods("POST")
	r.HandleFunc("/1.0/update", p.UpdateHandler).Methods("POST")
	r.HandleFunc("/1.0/deactivate", p.DeactivateHandler).Methods("POST")

	// Driver metadata
	r.HandleFunc("/", driverInfoHandler).Methods("GET")

	// Start server
	addr := fmt.Sprintf(":%d", config.Port)
	log.Printf("Universal Registrar Driver starting on %s", addr)
	log.Printf("Proxying to registrar at %s", config.RegistrarURL)

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
		"status":  "healthy",
		"service": "uniregistrar-driver-did-acc",
	})
}

func driverInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"driver":     "did-acc",
		"version":    "1.0.0",
		"methods":    []string{"acc"},
		"operations": []string{"create", "update", "deactivate"},
		"endpoints": map[string]string{
			"create":     "/1.0/create?method=acc",
			"update":     "/1.0/update?method=acc",
			"deactivate": "/1.0/deactivate?method=acc",
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
