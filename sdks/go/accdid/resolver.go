package accdid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opendlt/accu-did/sdks/go/accdid/httpx"
	"github.com/opendlt/accu-did/sdks/go/accdid/retry"
)

// ResolverClient provides DID resolution functionality
type ResolverClient struct {
	baseURL string
	doer    Doer
	logger  Logger
}

// NewResolverClient creates a new resolver client
func NewResolverClient(opts ClientOptions) (*ResolverClient, error) {
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

	return &ResolverClient{
		baseURL: opts.BaseURL,
		doer:    doer,
		logger:  opts.Logger,
	}, nil
}

// Resolve resolves a DID to its DID document using the native resolver API
func (c *ResolverClient) Resolve(ctx context.Context, did string) (*ResolutionResult, error) {
	if err := ValidateDID(did); err != nil {
		return nil, fmt.Errorf("invalid DID: %w", err)
	}

	c.logger.Debugf("Resolving DID: %s", did)

	// Use query parameters for the native resolver endpoint
	params := map[string]string{
		"did": did,
	}

	var result ResolutionResult
	status, body, err := httpx.DoJSONQuery(ctx, c.doer, c.baseURL, "/resolve", params, &result)
	if err != nil {
		_, classified := classifyError(err)
		return nil, classified
	}

	if status >= 400 {
		httpErr := decodeHTTPError(&http.Response{StatusCode: status, Status: fmt.Sprintf("%d", status)}, body)
		return nil, httpErr
	}

	c.logger.Debugf("Successfully resolved DID: %s", did)
	return &result, nil
}

// UniversalResolve resolves a DID using the Universal Resolver API format
func (c *ResolverClient) UniversalResolve(ctx context.Context, did string) (*ResolutionResult, error) {
	if err := ValidateDID(did); err != nil {
		return nil, fmt.Errorf("invalid DID: %w", err)
	}

	c.logger.Debugf("Resolving DID via Universal Resolver: %s", did)

	// Build endpoint with DID as path parameter
	endpoint := httpx.ParseEndpoint("/1.0/identifiers/{did}", map[string]string{
		"did": did,
	})

	var response UniversalResolveResponse
	status, body, err := httpx.DoJSON(ctx, c.doer, "GET", c.baseURL, endpoint, nil, &response)
	if err != nil {
		_, classified := classifyError(err)
		return nil, classified
	}

	if status >= 400 {
		httpErr := decodeHTTPError(&http.Response{StatusCode: status, Status: fmt.Sprintf("%d", status)}, body)
		return nil, httpErr
	}

	// Convert Universal Resolver response to our standard format
	result := &ResolutionResult{
		DIDDocument:      response.DIDDocument,
		Metadata:         response.DIDResolutionMetadata,
		DocumentMetadata: response.DIDDocumentMetadata,
	}

	c.logger.Debugf("Successfully resolved DID via Universal Resolver: %s", did)
	return result, nil
}

// Health checks the resolver service health
func (c *ResolverClient) Health(ctx context.Context) error {
	c.logger.Debugf("Checking resolver health")

	status, _, err := httpx.DoJSON(ctx, c.doer, "GET", c.baseURL, "/healthz", nil, nil)
	if err != nil {
		_, classified := classifyError(err)
		return classified
	}

	if status >= 400 {
		return fmt.Errorf("health check failed with status %d", status)
	}

	c.logger.Debugf("Resolver health check passed")
	return nil
}