package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// ConformanceConfig holds configuration for conformance testing
type ConformanceConfig struct {
	ResolverURL  string
	RegistrarURL string
	TestDID      string
	APIKey       string
}

// ConformanceResult holds the results of conformance testing
type ConformanceResult struct {
	Timestamp   time.Time              `json:"timestamp"`
	Config      ConformanceConfig      `json:"config"`
	Tests       []TestResult           `json:"tests"`
	Summary     TestSummary            `json:"summary"`
	TotalTimeMS int64                  `json:"total_time_ms"`
}

// TestResult holds the result of a single test
type TestResult struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Success     bool        `json:"success"`
	Error       string      `json:"error,omitempty"`
	StatusCode  int         `json:"status_code,omitempty"`
	TimeMS      int64       `json:"time_ms"`
	Response    interface{} `json:"response,omitempty"`
}

// TestSummary holds overall test statistics
type TestSummary struct {
	Total   int `json:"total"`
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Success bool `json:"success"`
}

func main() {
	start := time.Now()

	// Load configuration from environment
	config := ConformanceConfig{
		ResolverURL:  getEnvOrDefault("RESOLVER_URL", "http://127.0.0.1:8080"),
		RegistrarURL: getEnvOrDefault("REGISTRAR_URL", "http://127.0.0.1:8081"),
		TestDID:      getEnvOrDefault("TEST_DID", "did:acc:conformance-test"),
		APIKey:       os.Getenv("REGISTRAR_API_KEY"),
	}

	// Initialize result
	result := ConformanceResult{
		Timestamp: start,
		Config:    config,
		Tests:     []TestResult{},
	}

	// Run conformance tests
	log.Println("🔍 Starting DID conformance tests...")
	log.Printf("  Resolver:  %s", config.ResolverURL)
	log.Printf("  Registrar: %s", config.RegistrarURL)
	log.Printf("  Test DID:  %s", config.TestDID)

	// Test 1: Health checks
	result.Tests = append(result.Tests, testHealthCheck(config.ResolverURL+"/healthz", "Resolver Health Check"))
	result.Tests = append(result.Tests, testHealthCheck(config.RegistrarURL+"/healthz", "Registrar Health Check"))

	// Test 2: Create DID
	createResult := testCreateDID(config)
	result.Tests = append(result.Tests, createResult)

	if createResult.Success {
		// Test 3: Resolve newly created DID
		result.Tests = append(result.Tests, testResolveDID(config))

		// Test 4: Update DID (add service)
		updateResult := testUpdateDID(config)
		result.Tests = append(result.Tests, updateResult)

		if updateResult.Success {
			// Test 5: Resolve updated DID
			result.Tests = append(result.Tests, testResolveUpdatedDID(config))
		}

		// Test 6: Deactivate DID
		deactivateResult := testDeactivateDID(config)
		result.Tests = append(result.Tests, deactivateResult)

		if deactivateResult.Success {
			// Test 7: Resolve deactivated DID (should return 410)
			result.Tests = append(result.Tests, testResolveDeactivatedDID(config))
		}
	}

	// Calculate summary
	total := len(result.Tests)
	passed := 0
	for _, test := range result.Tests {
		if test.Success {
			passed++
		}
	}

	result.Summary = TestSummary{
		Total:   total,
		Passed:  passed,
		Failed:  total - passed,
		Success: passed == total,
	}

	result.TotalTimeMS = time.Since(start).Milliseconds()

	// Output JSON result
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}

	fmt.Println(string(jsonOutput))

	// Exit with appropriate code
	if result.Summary.Success {
		log.Printf("✅ All %d conformance tests passed", result.Summary.Total)
		os.Exit(0)
	} else {
		log.Printf("❌ %d of %d conformance tests failed", result.Summary.Failed, result.Summary.Total)
		os.Exit(1)
	}
}

func testHealthCheck(url, name string) TestResult {
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return TestResult{
			Name:        name,
			Description: "Check service health endpoint",
			Success:     false,
			Error:       err.Error(),
			TimeMS:      time.Since(start).Milliseconds(),
		}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	return TestResult{
		Name:        name,
		Description: "Check service health endpoint",
		Success:     resp.StatusCode == 200,
		StatusCode:  resp.StatusCode,
		TimeMS:      time.Since(start).Milliseconds(),
		Response:    string(body),
	}
}

func testCreateDID(config ConformanceConfig) TestResult {
	start := time.Now()

	// Create DID document
	didDoc := map[string]interface{}{
		"@context": []string{"https://www.w3.org/ns/did/v1"},
		"id":       config.TestDID,
		"verificationMethod": []map[string]interface{}{
			{
				"id":                 config.TestDID + "#key-1",
				"type":               "Ed25519VerificationKey2020",
				"controller":         config.TestDID,
				"publicKeyMultibase": "z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK",
			},
		},
		"authentication": []string{config.TestDID + "#key-1"},
	}

	createReq := map[string]interface{}{
		"didDocument": didDoc,
	}

	resp, err := makeRequest("POST", config.RegistrarURL+"/register", createReq, config.APIKey)
	if err != nil {
		return TestResult{
			Name:        "Create DID",
			Description: "Create a new DID document",
			Success:     false,
			Error:       err.Error(),
			TimeMS:      time.Since(start).Milliseconds(),
		}
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var response map[string]interface{}
	json.Unmarshal(body, &response)

	return TestResult{
		Name:        "Create DID",
		Description: "Create a new DID document",
		Success:     resp.StatusCode >= 200 && resp.StatusCode < 300,
		StatusCode:  resp.StatusCode,
		TimeMS:      time.Since(start).Milliseconds(),
		Response:    response,
	}
}

func testResolveDID(config ConformanceConfig) TestResult {
	start := time.Now()

	url := fmt.Sprintf("%s/resolve?did=%s", config.ResolverURL, config.TestDID)
	resp, err := http.Get(url)
	if err != nil {
		return TestResult{
			Name:        "Resolve DID",
			Description: "Resolve the newly created DID",
			Success:     false,
			Error:       err.Error(),
			TimeMS:      time.Since(start).Milliseconds(),
		}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)

	return TestResult{
		Name:        "Resolve DID",
		Description: "Resolve the newly created DID",
		Success:     resp.StatusCode == 200,
		StatusCode:  resp.StatusCode,
		TimeMS:      time.Since(start).Milliseconds(),
		Response:    response,
	}
}

func testUpdateDID(config ConformanceConfig) TestResult {
	start := time.Now()

	updateReq := map[string]interface{}{
		"did": config.TestDID,
		"patch": map[string]interface{}{
			"addService": []map[string]interface{}{
				{
					"id":              config.TestDID + "#test-service",
					"type":            "TestService",
					"serviceEndpoint": "https://test.example.com",
				},
			},
		},
	}

	resp, err := makeRequest("POST", config.RegistrarURL+"/native/update", updateReq, config.APIKey)
	if err != nil {
		return TestResult{
			Name:        "Update DID",
			Description: "Add a service to the DID document",
			Success:     false,
			Error:       err.Error(),
			TimeMS:      time.Since(start).Milliseconds(),
		}
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var response map[string]interface{}
	json.Unmarshal(body, &response)

	return TestResult{
		Name:        "Update DID",
		Description: "Add a service to the DID document",
		Success:     resp.StatusCode >= 200 && resp.StatusCode < 300,
		StatusCode:  resp.StatusCode,
		TimeMS:      time.Since(start).Milliseconds(),
		Response:    response,
	}
}

func testResolveUpdatedDID(config ConformanceConfig) TestResult {
	start := time.Now()

	url := fmt.Sprintf("%s/resolve?did=%s", config.ResolverURL, config.TestDID)
	resp, err := http.Get(url)
	if err != nil {
		return TestResult{
			Name:        "Resolve Updated DID",
			Description: "Resolve DID after update to verify changes",
			Success:     false,
			Error:       err.Error(),
			TimeMS:      time.Since(start).Milliseconds(),
		}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)

	// Check if service was added
	success := resp.StatusCode == 200
	if success {
		if didDoc, ok := response["didDocument"].(map[string]interface{}); ok {
			if services, ok := didDoc["service"].([]interface{}); ok && len(services) > 0 {
				success = true
			}
		}
	}

	return TestResult{
		Name:        "Resolve Updated DID",
		Description: "Resolve DID after update to verify changes",
		Success:     success,
		StatusCode:  resp.StatusCode,
		TimeMS:      time.Since(start).Milliseconds(),
		Response:    response,
	}
}

func testDeactivateDID(config ConformanceConfig) TestResult {
	start := time.Now()

	deactivateReq := map[string]interface{}{
		"did":        config.TestDID,
		"deactivate": true,
	}

	resp, err := makeRequest("POST", config.RegistrarURL+"/native/deactivate", deactivateReq, config.APIKey)
	if err != nil {
		return TestResult{
			Name:        "Deactivate DID",
			Description: "Deactivate the DID document",
			Success:     false,
			Error:       err.Error(),
			TimeMS:      time.Since(start).Milliseconds(),
		}
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var response map[string]interface{}
	json.Unmarshal(body, &response)

	return TestResult{
		Name:        "Deactivate DID",
		Description: "Deactivate the DID document",
		Success:     resp.StatusCode >= 200 && resp.StatusCode < 300,
		StatusCode:  resp.StatusCode,
		TimeMS:      time.Since(start).Milliseconds(),
		Response:    response,
	}
}

func testResolveDeactivatedDID(config ConformanceConfig) TestResult {
	start := time.Now()

	url := fmt.Sprintf("%s/resolve?did=%s", config.ResolverURL, config.TestDID)
	resp, err := http.Get(url)
	if err != nil {
		return TestResult{
			Name:        "Resolve Deactivated DID",
			Description: "Resolve deactivated DID (should return 410)",
			Success:     false,
			Error:       err.Error(),
			TimeMS:      time.Since(start).Milliseconds(),
		}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)

	return TestResult{
		Name:        "Resolve Deactivated DID",
		Description: "Resolve deactivated DID (should return 410)",
		Success:     resp.StatusCode == 410,
		StatusCode:  resp.StatusCode,
		TimeMS:      time.Since(start).Milliseconds(),
		Response:    response,
	}
}

func makeRequest(method, url string, body interface{}, apiKey string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	return client.Do(req)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}