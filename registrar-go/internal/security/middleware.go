package security

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

// SecurityConfig holds security configuration
type SecurityConfig struct {
	APIKey     string
	AllowList  []string // CIDR/IP list
	RateRPS    int      // requests per second
	RateBurst  int      // burst capacity
	BindAddr   string   // bind address
}

// DefaultSecurityConfig returns default configuration
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		RateRPS:   50,
		RateBurst: 100,
		BindAddr:  "127.0.0.1:8081",
	}
}

// APIKeyMiddleware enforces API key authentication if configured
func APIKeyMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip API key check if not configured
			if apiKey == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check Authorization header
			auth := r.Header.Get("Authorization")
			if auth == "" {
				writeErrorResponse(w, http.StatusUnauthorized, "unauthorized", "API key required")
				return
			}

			// Expect "Bearer <api-key>" or "ApiKey <api-key>" format
			var providedKey string
			if strings.HasPrefix(auth, "Bearer ") {
				providedKey = strings.TrimPrefix(auth, "Bearer ")
			} else if strings.HasPrefix(auth, "ApiKey ") {
				providedKey = strings.TrimPrefix(auth, "ApiKey ")
			} else {
				writeErrorResponse(w, http.StatusUnauthorized, "unauthorized", "Invalid authorization format")
				return
			}

			if providedKey != apiKey {
				writeErrorResponse(w, http.StatusUnauthorized, "unauthorized", "Invalid API key")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// IPAllowListMiddleware enforces IP allowlisting if configured
func IPAllowListMiddleware(allowList []string) func(http.Handler) http.Handler {
	// Parse CIDR blocks and IPs at startup
	var cidrs []*net.IPNet
	var ips []net.IP

	for _, item := range allowList {
		if strings.Contains(item, "/") {
			// CIDR block
			_, cidr, err := net.ParseCIDR(item)
			if err != nil {
				// Log error but don't fail startup
				continue
			}
			cidrs = append(cidrs, cidr)
		} else {
			// Single IP
			ip := net.ParseIP(item)
			if ip != nil {
				ips = append(ips, ip)
			}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if no allowlist configured
			if len(allowList) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			// Get client IP (handle X-Forwarded-For, X-Real-IP)
			clientIP := getClientIP(r)
			if clientIP == nil {
				writeErrorResponse(w, http.StatusForbidden, "forbidden", "Cannot determine client IP")
				return
			}

			// Check against allowed IPs
			allowed := false
			for _, ip := range ips {
				if ip.Equal(clientIP) {
					allowed = true
					break
				}
			}

			// Check against CIDR blocks
			if !allowed {
				for _, cidr := range cidrs {
					if cidr.Contains(clientIP) {
						allowed = true
						break
					}
				}
			}

			if !allowed {
				writeErrorResponse(w, http.StatusForbidden, "forbidden", "IP address not allowed")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware implements token bucket rate limiting
func RateLimitMiddleware(rps int, burst int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				writeErrorResponse(w, http.StatusTooManyRequests, "rateLimitExceeded", "Rate limit exceeded")
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

// getClientIP extracts the real client IP from request
func getClientIP(r *http.Request) net.IP {
	// Try X-Forwarded-For first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Take first IP from comma-separated list
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := net.ParseIP(strings.TrimSpace(ips[0]))
			if ip != nil {
				return ip
			}
		}
	}

	// Try X-Real-IP
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		ip := net.ParseIP(xri)
		if ip != nil {
			return ip
		}
	}

	// Fall back to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil
	}
	return net.ParseIP(host)
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