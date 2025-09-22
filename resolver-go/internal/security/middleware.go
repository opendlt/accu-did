package security

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// SecurityConfig holds security configuration for resolver
type SecurityConfig struct {
	CORSAllowOrigins []string // CORS allowed origins
	BindAddr         string   // bind address
}

// DefaultSecurityConfig returns default configuration
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		BindAddr: "127.0.0.1:8080",
	}
}

// CORSMiddleware handles CORS configuration
func CORSMiddleware(allowOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// If no origins configured, no CORS headers
			if len(allowOrigins) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range allowOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-Id")
				w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
			}

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				if allowed {
					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequestIDMiddleware adds X-Request-Id header if not present
func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = generateRequestID()
				r.Header.Set("X-Request-Id", requestID)
			}

			// Add to response headers
			w.Header().Set("X-Request-Id", requestID)

			// Add to context for logging
			ctx := context.WithValue(r.Context(), "request_id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// generateRequestID creates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// writeErrorResponse writes a canonical error response
func writeErrorResponse(w http.ResponseWriter, status int, errorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := fmt.Sprintf(`{"code":%d,"error":"%s","message":"%s"}`, status, errorCode, message)
	w.Write([]byte(response))
}