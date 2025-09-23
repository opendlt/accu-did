// Package httpx provides HTTP utilities for the Accumulate DID SDK
package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Doer is the interface for HTTP request execution
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// DoJSON performs an HTTP request with JSON encoding/decoding
func DoJSON(ctx context.Context, doer Doer, method, baseURL, endpoint string, in interface{}, out interface{}) (status int, body []byte, err error) {
	// Build URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = path.Join(u.Path, endpoint)

	var reqBody io.Reader
	if in != nil {
		jsonBytes, err := json.Marshal(in)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create request: %w", err)
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := doer.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if out != nil && len(body) > 0 {
		if err := json.Unmarshal(body, out); err != nil {
			return resp.StatusCode, body, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return resp.StatusCode, body, nil
}

// DoJSONQuery performs a GET request with query parameters
func DoJSONQuery(ctx context.Context, doer Doer, baseURL, endpoint string, params map[string]string, out interface{}) (status int, body []byte, err error) {
	// Build URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = path.Join(u.Path, endpoint)

	if len(params) > 0 {
		q := u.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := doer.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if out != nil && len(body) > 0 && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if err := json.Unmarshal(body, out); err != nil {
			return resp.StatusCode, body, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return resp.StatusCode, body, nil
}

// BuildURL safely joins a base URL with an endpoint
func BuildURL(baseURL, endpoint string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Clean and join paths
	u.Path = path.Join(u.Path, endpoint)
	return u.String(), nil
}

// ParseEndpoint safely parses an endpoint with path parameters
func ParseEndpoint(template string, params map[string]string) string {
	result := template
	for k, v := range params {
		placeholder := "{" + k + "}"
		// For Universal Resolver endpoints, we need proper URL encoding
		result = strings.ReplaceAll(result, placeholder, url.PathEscape(v))
	}
	return result
}