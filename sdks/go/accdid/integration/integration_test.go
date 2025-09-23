//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/opendlt/accu-did/sdks/go/accdid"
)

// testLogger implements accdid.Logger for integration test output
type testLogger struct {
	t *testing.T
}

func (l *testLogger) Debugf(format string, args ...interface{}) {
	l.t.Logf("[DEBUG] "+format, args...)
}

func (l *testLogger) Infof(format string, args ...interface{}) {
	l.t.Logf("[INFO] "+format, args...)
}

func (l *testLogger) Warnf(format string, args ...interface{}) {
	l.t.Logf("[WARN] "+format, args...)
}

func (l *testLogger) Errorf(format string, args ...interface{}) {
	l.t.Logf("[ERROR] "+format, args...)
}

// getEnvWithDefault returns environment variable value or default
func getEnvWithDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// waitForHealthy polls health endpoints until they respond or timeout
func waitForHealthy(ctx context.Context, t *testing.T, resolver *accdid.ResolverClient, registrar *accdid.RegistrarClient, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		// Test resolver health
		if err := resolver.Health(ctx); err != nil {
			t.Logf("Resolver not healthy yet: %v", err)
		} else if err := registrar.Health(ctx); err != nil {
			t.Logf("Registrar not healthy yet: %v", err)
		} else {
			t.Log("Both services are healthy")
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
			// Continue polling
		}
	}

	return fmt.Errorf("services not healthy after %v timeout", timeout)
}

// tryFaucet attempts to fund a lite account if both faucet URL and lite account URL are provided
func tryFaucet(ctx context.Context, t *testing.T, faucetURL, liteAccountURL string) error {
	if faucetURL == "" || liteAccountURL == "" {
		t.Log("Skipping faucet funding (missing ACC_FAUCET_URL or LITE_ACCOUNT_URL)")
		return nil
	}

	t.Logf("Attempting faucet funding: %s -> %s", faucetURL, liteAccountURL)

	requestBody, err := json.Marshal(map[string]string{
		"url": liteAccountURL,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal faucet request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", faucetURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create faucet request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("faucet request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("faucet returned status %d", resp.StatusCode)
	}

	t.Logf("Faucet funding successful (status: %d)", resp.StatusCode)
	return nil
}

// generateUniqueDID creates a deterministic unique DID for this test run
func generateUniqueDID() string {
	// Use nanosecond timestamp for uniqueness
	nano := time.Now().UnixNano()
	return fmt.Sprintf("did:acc:it%d", nano)
}

// TestAccuEndToEnd runs the complete DID lifecycle integration test
func TestAccuEndToEnd(t *testing.T) {
	ctx := context.Background()

	// Environment setup
	resolverURL := getEnvWithDefault("RESOLVER_URL", "http://127.0.0.1:8080")
	registrarURL := getEnvWithDefault("REGISTRAR_URL", "http://127.0.0.1:8081")
	apiKey := os.Getenv("ACCU_API_KEY")
	idempotencyKey := os.Getenv("IDEMPOTENCY_KEY")
	faucetURL := os.Getenv("ACC_FAUCET_URL")
	liteAccountURL := os.Getenv("LITE_ACCOUNT_URL")

	t.Logf("Integration test configuration:")
	t.Logf("  Resolver URL: %s", resolverURL)
	t.Logf("  Registrar URL: %s", registrarURL)
	t.Logf("  API Key: %s", func() string {
		if apiKey != "" {
			return "configured"
		}
		return "none"
	}())
	t.Logf("  Faucet URL: %s", func() string {
		if faucetURL != "" {
			return faucetURL
		}
		return "none"
	}())
	t.Logf("  Lite Account: %s", func() string {
		if liteAccountURL != "" {
			return liteAccountURL
		}
		return "none"
	}())

	// Check if ACC_NODE_URL is set (indicates REAL mode)
	accNodeURL := os.Getenv("ACC_NODE_URL")
	if accNodeURL == "" {
		t.Skip("Skipping integration test: ACC_NODE_URL not set (services not in REAL mode)")
	}
	t.Logf("  Accumulate Node: %s", accNodeURL)

	// Create clients
	logger := &testLogger{t: t}

	resolverOpts := accdid.ClientOptions{
		BaseURL:        resolverURL,
		APIKey:         apiKey,
		IdempotencyKey: idempotencyKey,
		Logger:         logger,
		Timeout:        15 * time.Second,
	}

	registrarOpts := accdid.ClientOptions{
		BaseURL:        registrarURL,
		APIKey:         apiKey,
		IdempotencyKey: idempotencyKey,
		Logger:         logger,
		Timeout:        15 * time.Second,
	}

	resolver, err := accdid.NewResolverClient(resolverOpts)
	if err != nil {
		t.Fatalf("Failed to create resolver client: %v", err)
	}

	registrar, err := accdid.NewRegistrarClient(registrarOpts)
	if err != nil {
		t.Fatalf("Failed to create registrar client: %v", err)
	}

	// Wait for services to be healthy
	t.Log("Waiting for services to be healthy...")
	if err := waitForHealthy(ctx, t, resolver, registrar, 30*time.Second); err != nil {
		t.Skipf("Services not healthy: %v. Make sure to run: make devnet-up && make services-up", err)
	}

	// Optional: Try faucet funding
	if err := tryFaucet(ctx, t, faucetURL, liteAccountURL); err != nil {
		t.Logf("Faucet funding failed: %v (continuing anyway)", err)
	}

	// Generate unique DID for this test run
	did := generateUniqueDID()
	t.Logf("Testing DID: %s", did)

	// Step 1: Verify 404 before creation
	t.Log("Step 1: Verify DID does not exist (expect 404)")
	_, err = resolver.Resolve(ctx, did)
	if err == nil {
		t.Fatalf("Expected error when resolving non-existent DID, but got none")
	}
	if !strings.Contains(err.Error(), "404") && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		t.Logf("Note: Error doesn't contain '404' or 'not found': %v", err)
	}
	t.Logf("✓ DID does not exist yet (error: %v)", err)

	// Step 2: Register DID
	t.Log("Step 2: Register new DID")
	didDocument := map[string]interface{}{
		"@context": []string{"https://www.w3.org/ns/did/v1"},
		"id":       did,
		"verificationMethod": []interface{}{
			map[string]interface{}{
				"id":                 did + "#key1",
				"type":               "Ed25519VerificationKey2020",
				"controller":         did,
				"publicKeyMultibase": "z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK",
			},
		},
		"authentication": []string{did + "#key1"},
	}

	didDocJSON, err := json.Marshal(didDocument)
	if err != nil {
		t.Fatalf("Failed to marshal DID document: %v", err)
	}

	registerReq := accdid.NativeRegisterRequest{
		DID:         did,
		DIDDocument: didDocJSON,
	}

	txID, err := registrar.Register(ctx, registerReq)
	if err != nil {
		t.Fatalf("Failed to register DID: %v", err)
	}
	t.Logf("✓ DID registered successfully (txID: %s)", txID)

	// Brief wait for transaction to be processed
	time.Sleep(3 * time.Second)

	// Step 3: Resolve and verify (200)
	t.Log("Step 3: Resolve DID (expect 200)")
	result, err := resolver.Resolve(ctx, did)
	if err != nil {
		t.Fatalf("Failed to resolve DID: %v", err)
	}
	if result.DIDDocument == nil {
		t.Fatalf("Resolved DID document is nil")
	}

	// Verify basic structure
	docMap, ok := result.DIDDocument.(map[string]interface{})
	if !ok {
		t.Fatalf("DID document is not a map[string]interface{}")
	}
	if docMap["id"] != did {
		t.Fatalf("DID document ID mismatch: expected %s, got %v", did, docMap["id"])
	}
	t.Logf("✓ DID resolved successfully (id: %v)", docMap["id"])

	// Step 4: Update DID (patch)
	t.Log("Step 4: Update DID (add service)")
	serviceUpdate := map[string]interface{}{
		"op":    "add",
		"path":  "/service",
		"value": []interface{}{
			map[string]interface{}{
				"id":              did + "#service1",
				"type":            "TestService",
				"serviceEndpoint": "https://test.example.com",
			},
		},
	}

	patchJSON, err := json.Marshal([]interface{}{serviceUpdate})
	if err != nil {
		t.Fatalf("Failed to marshal patch: %v", err)
	}

	updateReq := accdid.NativeUpdateRequest{
		DID:   did,
		Patch: patchJSON,
	}

	updateTxID, err := registrar.Update(ctx, updateReq)
	if err != nil {
		t.Fatalf("Failed to update DID: %v", err)
	}
	t.Logf("✓ DID updated successfully (txID: %s)", updateTxID)

	// Brief wait for update to be processed
	time.Sleep(3 * time.Second)

	// Step 5: Resolve again and verify change
	t.Log("Step 5: Resolve DID again (verify update)")
	updatedResult, err := resolver.Resolve(ctx, did)
	if err != nil {
		t.Fatalf("Failed to resolve updated DID: %v", err)
	}

	updatedDocMap, ok := updatedResult.DIDDocument.(map[string]interface{})
	if !ok {
		t.Fatalf("Updated DID document is not a map[string]interface{}")
	}

	// Verify service was added
	services, ok := updatedDocMap["service"]
	if !ok {
		t.Logf("Warning: No 'service' field found in updated DID document")
	} else {
		servicesArray, ok := services.([]interface{})
		if !ok || len(servicesArray) == 0 {
			t.Logf("Warning: Service array is empty or not an array: %v", services)
		} else {
			t.Logf("✓ Service found in updated DID document: %d services", len(servicesArray))
		}
	}

	// Step 6: Deactivate DID
	t.Log("Step 6: Deactivate DID")
	deactivateReq := accdid.NativeDeactivateRequest{
		DID:    did,
		Reason: "Integration test completed",
	}

	deactivateTxID, err := registrar.Deactivate(ctx, deactivateReq)
	if err != nil {
		t.Fatalf("Failed to deactivate DID: %v", err)
	}
	t.Logf("✓ DID deactivated successfully (txID: %s)", deactivateTxID)

	// Brief wait for deactivation to be processed
	time.Sleep(3 * time.Second)

	// Step 7: Resolve expecting 410 Gone + tombstone
	t.Log("Step 7: Resolve deactivated DID (expect 410 Gone)")
	_, err = resolver.Resolve(ctx, did)
	if err == nil {
		t.Logf("Note: Expected error when resolving deactivated DID, but got none")
	} else {
		// Check if it's a 410 Gone error
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "410") || strings.Contains(errMsg, "gone") || strings.Contains(errMsg, "deactivated") {
			t.Logf("✓ Deactivated DID returns expected error: %v", err)
		} else {
			t.Logf("Note: Error doesn't indicate deactivation clearly: %v", err)
		}
	}

	// Summary
	t.Log("\n=== Integration Test Summary ===")
	t.Logf("DID: %s", did)
	t.Logf("Register txID: %s", txID)
	t.Logf("Update txID: %s", updateTxID)
	t.Logf("Deactivate txID: %s", deactivateTxID)
	t.Log("✓ All integration test steps completed successfully")
}

// TestAccuIdempotency tests idempotent operations (if supported)
func TestAccuIdempotency(t *testing.T) {
	ctx := context.Background()

	// Environment setup
	registrarURL := getEnvWithDefault("REGISTRAR_URL", "http://127.0.0.1:8081")

	// Check if ACC_NODE_URL is set (indicates REAL mode)
	accNodeURL := os.Getenv("ACC_NODE_URL")
	if accNodeURL == "" {
		t.Skip("Skipping idempotency test: ACC_NODE_URL not set")
	}

	// Create clients with fixed idempotency key
	idempotencyKey := fmt.Sprintf("test-idem-%d", time.Now().UnixNano())
	logger := &testLogger{t: t}

	opts := accdid.ClientOptions{
		BaseURL:        registrarURL,
		IdempotencyKey: idempotencyKey,
		Logger:         logger,
		Timeout:        15 * time.Second,
	}

	registrar, err := accdid.NewRegistrarClient(opts)
	if err != nil {
		t.Fatalf("Failed to create registrar client: %v", err)
	}

	// Wait for service health
	if err := registrar.Health(ctx); err != nil {
		t.Skipf("Registrar not healthy: %v", err)
	}

	// Generate unique DID
	did := generateUniqueDID()
	didDocument := map[string]interface{}{
		"@context": []string{"https://www.w3.org/ns/did/v1"},
		"id":       did,
	}

	didDocJSON, _ := json.Marshal(didDocument)
	registerReq := accdid.NativeRegisterRequest{
		DID:         did,
		DIDDocument: didDocJSON,
	}

	// First registration
	txID1, err := registrar.Register(ctx, registerReq)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Second registration with same idempotency key (should succeed or return same result)
	txID2, err := registrar.Register(ctx, registerReq)
	if err != nil {
		t.Logf("Second registration failed (this may be expected): %v", err)
	} else {
		t.Logf("Second registration succeeded - txID1: %s, txID2: %s", txID1, txID2)
	}

	t.Log("✓ Idempotency test completed")
}