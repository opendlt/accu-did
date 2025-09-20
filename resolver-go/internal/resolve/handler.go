package resolve

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/opendlt/accu-did/resolver-go/internal/acc"
)

// Handler handles DID resolution requests
type Handler struct {
	accClient acc.Client
}

// NewHandler creates a new resolve handler
func NewHandler(accClient acc.Client) *Handler {
	return &Handler{
		accClient: accClient,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error     string            `json:"error"`
	Message   string            `json:"message"`
	Details   map[string]string `json:"details,omitempty"`
	RequestID string            `json:"requestId,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// Resolve handles GET /resolve requests
func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	// Extract DID from query parameter
	did := r.URL.Query().Get("did")
	if did == "" {
		h.writeError(w, "invalidDid", "DID parameter is required", http.StatusBadRequest, nil)
		return
	}

	// Extract versionTime if provided
	var versionTime *time.Time
	if vt := r.URL.Query().Get("versionTime"); vt != "" {
		parsed, err := time.Parse(time.RFC3339, vt)
		if err != nil {
			// Try Unix timestamp
			if parsed, err = time.Parse("1704067200", vt); err != nil {
				h.writeError(w, "invalidVersionTime", "Invalid versionTime format", http.StatusBadRequest, map[string]string{
					"versionTime": vt,
					"expected":    "ISO 8601 or Unix timestamp",
				})
				return
			}
		}
		versionTime = &parsed
	}

	// Resolve DID
	result, err := ResolveDID(h.accClient, did, versionTime)
	if err != nil {
		switch err.(type) {
		case *NotFoundError:
			h.writeError(w, "notFound", err.Error(), http.StatusNotFound, nil)
		case *InvalidDIDError:
			h.writeError(w, "invalidDid", err.Error(), http.StatusBadRequest, nil)
		case *DeactivatedError:
			h.writeError(w, "deactivated", err.Error(), http.StatusGone, nil)
		default:
			h.writeError(w, "internalError", "Internal server error", http.StatusInternalServerError, nil)
		}
		return
	}

	// Return successful resolution
	w.Header().Set("Content-Type", "application/did+ld+json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.writeError(w, "internalError", "Failed to encode response", http.StatusInternalServerError, nil)
		return
	}
}

func (h *Handler) writeError(w http.ResponseWriter, errorCode, message string, status int, details map[string]string) {
	response := ErrorResponse{
		Error:     errorCode,
		Message:   message,
		Details:   details,
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(response)
}