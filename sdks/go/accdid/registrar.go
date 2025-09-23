package accdid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opendlt/accu-did/sdks/go/accdid/httpx"
	"github.com/opendlt/accu-did/sdks/go/accdid/retry"
)

// RegistrarClient provides DID registration and management functionality
type RegistrarClient struct {
	baseURL string
	doer    Doer
	logger  Logger
}

// NewRegistrarClient creates a new registrar client
func NewRegistrarClient(opts ClientOptions) (*RegistrarClient, error) {
	if opts.BaseURL == "" {
		return nil, fmt.Errorf("BaseURL is required")
	}

	opts = applyDefaults(opts)

	// Build middleware chain
	var doer Doer = opts.HTTP

	// Add retry middleware
	doer = retry.WithRetry(doer, retry.Policy{
		Max:       opts.Retries.Max,
		BaseDelay: opts.Retries.BaseDelay,
		MaxDelay:  opts.Retries.MaxDelay,
		Jitter:    opts.Retries.Jitter,
	}, func(err error) bool {
		retryable, _ := classifyError(err)
		return retryable
	})

	// Add header and timeout middleware
	doer = WithMiddleware(doer,
		headerMiddleware(opts.APIKey, opts.IdempotencyKey, opts.RequestID),
		timeoutMiddleware(opts.Timeout),
	)

	return &RegistrarClient{
		baseURL: opts.BaseURL,
		doer:    doer,
		logger:  opts.Logger,
	}, nil
}

// Register creates a new DID with the specified DID document
func (c *RegistrarClient) Register(ctx context.Context, req NativeRegisterRequest) (string, error) {
	if err := ValidateDID(req.DID); err != nil {
		return "", fmt.Errorf("invalid DID: %w", err)
	}

	c.logger.Debugf("Registering DID: %s", req.DID)

	var response RegistrarResponse
	status, body, err := httpx.DoJSON(ctx, c.doer, "POST", c.baseURL, "/register", req, &response)
	if err != nil {
		_, classified := classifyError(err)
		return "", classified
	}

	if status >= 400 {
		httpErr := decodeHTTPError(&http.Response{StatusCode: status, Status: fmt.Sprintf("%d", status)}, body)
		return "", httpErr
	}

	c.logger.Infof("Successfully registered DID: %s, txID: %s", req.DID, response.TransactionID)
	return response.TransactionID, nil
}

// Update modifies an existing DID document using a patch
func (c *RegistrarClient) Update(ctx context.Context, req NativeUpdateRequest) (string, error) {
	if err := ValidateDID(req.DID); err != nil {
		return "", fmt.Errorf("invalid DID: %w", err)
	}

	c.logger.Debugf("Updating DID: %s", req.DID)

	var response RegistrarResponse
	status, body, err := httpx.DoJSON(ctx, c.doer, "POST", c.baseURL, "/native/update", req, &response)
	if err != nil {
		_, classified := classifyError(err)
		return "", classified
	}

	if status >= 400 {
		httpErr := decodeHTTPError(&http.Response{StatusCode: status, Status: fmt.Sprintf("%d", status)}, body)
		return "", httpErr
	}

	c.logger.Infof("Successfully updated DID: %s, txID: %s", req.DID, response.TransactionID)
	return response.TransactionID, nil
}

// Deactivate marks a DID as deactivated
func (c *RegistrarClient) Deactivate(ctx context.Context, req NativeDeactivateRequest) (string, error) {
	if err := ValidateDID(req.DID); err != nil {
		return "", fmt.Errorf("invalid DID: %w", err)
	}

	c.logger.Debugf("Deactivating DID: %s", req.DID)

	var response RegistrarResponse
	status, body, err := httpx.DoJSON(ctx, c.doer, "POST", c.baseURL, "/native/deactivate", req, &response)
	if err != nil {
		_, classified := classifyError(err)
		return "", classified
	}

	if status >= 400 {
		httpErr := decodeHTTPError(&http.Response{StatusCode: status, Status: fmt.Sprintf("%d", status)}, body)
		return "", httpErr
	}

	c.logger.Infof("Successfully deactivated DID: %s, txID: %s", req.DID, response.TransactionID)
	return response.TransactionID, nil
}

// UniversalCreate creates a DID using Universal Registrar v1.0 format
func (c *RegistrarClient) UniversalCreate(ctx context.Context, didDocument interface{}) (string, error) {
	c.logger.Debugf("Creating DID via Universal Registrar")

	requestBody := map[string]interface{}{
		"didDocument": didDocument,
	}

	var response RegistrarResponse
	status, body, err := httpx.DoJSON(ctx, c.doer, "POST", c.baseURL, "/1.0/create", requestBody, &response)
	if err != nil {
		_, classified := classifyError(err)
		return "", classified
	}

	if status >= 400 {
		httpErr := decodeHTTPError(&http.Response{StatusCode: status, Status: fmt.Sprintf("%d", status)}, body)
		return "", httpErr
	}

	c.logger.Infof("Successfully created DID via Universal Registrar, txID: %s", response.TransactionID)
	return response.TransactionID, nil
}

// UniversalUpdate updates a DID using Universal Registrar v1.0 format
func (c *RegistrarClient) UniversalUpdate(ctx context.Context, did string, patch interface{}) (string, error) {
	if err := ValidateDID(did); err != nil {
		return "", fmt.Errorf("invalid DID: %w", err)
	}

	c.logger.Debugf("Updating DID via Universal Registrar: %s", did)

	requestBody := map[string]interface{}{
		"did":   did,
		"patch": patch,
	}

	var response RegistrarResponse
	status, body, err := httpx.DoJSON(ctx, c.doer, "POST", c.baseURL, "/1.0/update", requestBody, &response)
	if err != nil {
		_, classified := classifyError(err)
		return "", classified
	}

	if status >= 400 {
		httpErr := decodeHTTPError(&http.Response{StatusCode: status, Status: fmt.Sprintf("%d", status)}, body)
		return "", httpErr
	}

	c.logger.Infof("Successfully updated DID via Universal Registrar: %s, txID: %s", did, response.TransactionID)
	return response.TransactionID, nil
}

// UniversalDeactivate deactivates a DID using Universal Registrar v1.0 format
func (c *RegistrarClient) UniversalDeactivate(ctx context.Context, did string) (string, error) {
	if err := ValidateDID(did); err != nil {
		return "", fmt.Errorf("invalid DID: %w", err)
	}

	c.logger.Debugf("Deactivating DID via Universal Registrar: %s", did)

	requestBody := map[string]interface{}{
		"did": did,
	}

	var response RegistrarResponse
	status, body, err := httpx.DoJSON(ctx, c.doer, "POST", c.baseURL, "/1.0/deactivate", requestBody, &response)
	if err != nil {
		_, classified := classifyError(err)
		return "", classified
	}

	if status >= 400 {
		httpErr := decodeHTTPError(&http.Response{StatusCode: status, Status: fmt.Sprintf("%d", status)}, body)
		return "", httpErr
	}

	c.logger.Infof("Successfully deactivated DID via Universal Registrar: %s, txID: %s", did, response.TransactionID)
	return response.TransactionID, nil
}

// Health checks the registrar service health
func (c *RegistrarClient) Health(ctx context.Context) error {
	c.logger.Debugf("Checking registrar health")

	status, _, err := httpx.DoJSON(ctx, c.doer, "GET", c.baseURL, "/healthz", nil, nil)
	if err != nil {
		_, classified := classifyError(err)
		return classified
	}

	if status >= 400 {
		return fmt.Errorf("health check failed with status %d", status)
	}

	c.logger.Debugf("Registrar health check passed")
	return nil
}