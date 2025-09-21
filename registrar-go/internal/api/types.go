package api

import "time"

type DeactivateRequest struct {
	DID     string                 `json:"did"`
	Options map[string]interface{} `json:"options,omitempty"`
	Secret  map[string]interface{} `json:"secret,omitempty"`
}

type UpdateRequest struct {
	DID      string         `json:"did"`
	Document map[string]any `json:"document"`
	Options  map[string]any `json:"options,omitempty"`
	Secret   map[string]any `json:"secret,omitempty"`
}

type ErrorResponse struct {
	Error     string            `json:"error"`
	Message   string            `json:"message"`
	Details   map[string]string `json:"details,omitempty"`
	RequestID string            `json:"requestId,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}
