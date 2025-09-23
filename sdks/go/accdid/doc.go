// Package accdid provides a Go SDK for interacting with Accumulate DID resolver and registrar services.
//
// The SDK offers a clean, production-ready HTTP client with automatic retries, comprehensive error handling,
// and support for both development (FAKE) and production (REAL) modes.
//
// Basic usage:
//
//	import "github.com/opendlt/accu-did/sdks/go/accdid"
//
//	// Create a resolver client
//	resolver, err := accdid.NewResolverClient(accdid.ClientOptions{
//	    BaseURL: "http://localhost:8080",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Resolve a DID
//	result, err := resolver.Resolve(context.Background(), "did:acc:alice")
//	if err != nil {
//	    // Handle deactivated DIDs (410 Gone)
//	    if errors.Is(err, accdid.ErrGoneDeactivated) {
//	        fmt.Println("DID has been deactivated")
//	        return
//	    }
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("DID Document: %+v\n", result.DIDDocument)
//
// The SDK automatically handles:
//   - Exponential backoff with jitter for retryable errors (429, 502, 503, 504)
//   - Request tracking with X-Request-Id headers
//   - Optional API key authentication
//   - Idempotency keys for safe retries
//   - Deterministic resolution ordering
//   - Canonical deactivation tombstone handling (410 responses)
//
// For more examples and documentation, see the README.
package accdid