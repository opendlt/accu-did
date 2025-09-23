// Package main demonstrates basic usage of the Accumulate DID SDK
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/opendlt/accu-did/sdks/go/accdid"
)

func main() {
	// Get configuration from environment variables
	resolverURL := getEnvOrDefault("RESOLVER_URL", "http://127.0.0.1:8080")
	registrarURL := getEnvOrDefault("REGISTRAR_URL", "http://127.0.0.1:8081")
	apiKey := os.Getenv("ACCU_API_KEY")         // Optional
	idempotencyKey := os.Getenv("IDEMPOTENCY_KEY") // Optional

	fmt.Printf("🚀 Accumulate DID SDK Example\n")
	fmt.Printf("Resolver URL: %s\n", resolverURL)
	fmt.Printf("Registrar URL: %s\n", registrarURL)
	fmt.Printf("SDK Version: %s\n\n", accdid.Version)

	// Create clients
	resolverClient, err := accdid.NewResolverClient(accdid.ClientOptions{
		BaseURL: resolverURL,
		APIKey:  apiKey,
	})
	if err != nil {
		log.Fatalf("Failed to create resolver client: %v", err)
	}

	registrarClient, err := accdid.NewRegistrarClient(accdid.ClientOptions{
		BaseURL:        registrarURL,
		APIKey:         apiKey,
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		log.Fatalf("Failed to create registrar client: %v", err)
	}

	ctx := context.Background()

	// Test health endpoints
	fmt.Printf("🔍 Checking service health...\n")

	if err := resolverClient.Health(ctx); err != nil {
		log.Printf("⚠️  Resolver health check failed: %v", err)
	} else {
		fmt.Printf("✅ Resolver is healthy\n")
	}

	if err := registrarClient.Health(ctx); err != nil {
		log.Printf("⚠️  Registrar health check failed: %v", err)
	} else {
		fmt.Printf("✅ Registrar is healthy\n")
	}

	fmt.Println()

	// Example 1: Try to resolve an existing DID
	fmt.Printf("📖 Example 1: Resolving existing DID\n")
	existingDID := "did:acc:alice"

	result, err := resolverClient.Resolve(ctx, existingDID)
	if err != nil {
		if errors.Is(err, accdid.ErrNotFound) {
			fmt.Printf("❌ DID not found: %s\n", existingDID)
		} else if errors.Is(err, accdid.ErrGoneDeactivated) {
			fmt.Printf("⚠️  DID has been deactivated: %s\n", existingDID)
			// In case of deactivation, we might still get the tombstone document
			if result != nil && result.DIDDocument != nil {
				fmt.Printf("📄 Deactivation tombstone received\n")
				printJSON("Tombstone", result.DIDDocument)
			}
		} else {
			fmt.Printf("❌ Resolution failed: %v\n", err)
		}
	} else {
		fmt.Printf("✅ Successfully resolved: %s\n", existingDID)
		printJSON("DID Document", result.DIDDocument)
		if result.Metadata != nil {
			printJSON("Resolution Metadata", result.Metadata)
		}
	}

	fmt.Println()

	// Example 2: Register, Update, and Deactivate workflow (FAKE mode)
	// Note: This will only work if the services are running in FAKE mode
	// or if you have a funded lite account with sufficient credits in REAL mode
	testDID := "did:acc:sdktest"

	fmt.Printf("📝 Example 2: Complete DID lifecycle (FAKE mode or funded account required)\n")

	// Create a simple DID document
	didDocument := map[string]interface{}{
		"@context": []string{"https://www.w3.org/ns/did/v1"},
		"id":       testDID,
		"verificationMethod": []map[string]interface{}{{
			"id":                 testDID + "#key1",
			"type":               "Ed25519VerificationKey2020",
			"controller":         testDID,
			"publicKeyMultibase": "z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK",
		}},
		"authentication":  []string{testDID + "#key1"},
		"assertionMethod": []string{testDID + "#key1"},
	}

	didDocBytes, _ := json.Marshal(didDocument)

	// Step 1: Register the DID
	fmt.Printf("📝 Registering DID: %s\n", testDID)

	registerReq := accdid.NativeRegisterRequest{
		DID:         testDID,
		DIDDocument: json.RawMessage(didDocBytes),
	}

	txID, err := registrarClient.Register(ctx, registerReq)
	if err != nil {
		if errors.Is(err, accdid.ErrBadRequest) {
			fmt.Printf("⚠️  Registration failed (possibly already exists or insufficient credits): %v\n", err)
		} else {
			fmt.Printf("❌ Registration failed: %v\n", err)
		}
	} else {
		fmt.Printf("✅ DID registered successfully, Transaction ID: %s\n", txID)

		// Step 2: Resolve the newly created DID
		fmt.Printf("📖 Resolving newly created DID: %s\n", testDID)
		result, err := resolverClient.Resolve(ctx, testDID)
		if err != nil {
			fmt.Printf("❌ Failed to resolve new DID: %v\n", err)
		} else {
			fmt.Printf("✅ Successfully resolved new DID\n")
			printJSON("New DID Document", result.DIDDocument)
		}

		// Step 3: Update the DID (add a service)
		fmt.Printf("🔄 Updating DID: %s\n", testDID)

		updatePatch := map[string]interface{}{
			"addService": map[string]interface{}{
				"id":              testDID + "#website",
				"type":            "LinkedDomains",
				"serviceEndpoint": "https://example.com",
			},
		}

		patchBytes, _ := json.Marshal(updatePatch)
		updateReq := accdid.NativeUpdateRequest{
			DID:   testDID,
			Patch: json.RawMessage(patchBytes),
		}

		txID, err = registrarClient.Update(ctx, updateReq)
		if err != nil {
			fmt.Printf("❌ Update failed: %v\n", err)
		} else {
			fmt.Printf("✅ DID updated successfully, Transaction ID: %s\n", txID)
		}

		// Step 4: Deactivate the DID
		fmt.Printf("🗑️  Deactivating DID: %s\n", testDID)

		deactivateReq := accdid.NativeDeactivateRequest{
			DID:    testDID,
			Reason: "SDK example complete",
		}

		txID, err = registrarClient.Deactivate(ctx, deactivateReq)
		if err != nil {
			fmt.Printf("❌ Deactivation failed: %v\n", err)
		} else {
			fmt.Printf("✅ DID deactivated successfully, Transaction ID: %s\n", txID)

			// Step 5: Try to resolve deactivated DID (should return 410)
			fmt.Printf("📖 Attempting to resolve deactivated DID: %s\n", testDID)
			result, err := resolverClient.Resolve(ctx, testDID)
			if err != nil {
				if errors.Is(err, accdid.ErrGoneDeactivated) {
					fmt.Printf("✅ Correctly received 410 Gone for deactivated DID\n")
					if result != nil && result.DIDDocument != nil {
						printJSON("Deactivation Tombstone", result.DIDDocument)
					}
				} else {
					fmt.Printf("❌ Unexpected error resolving deactivated DID: %v\n", err)
				}
			} else {
				fmt.Printf("⚠️  DID resolved successfully (might not be deactivated yet)\n")
			}
		}
	}

	fmt.Println()
	fmt.Printf("🎉 SDK example complete!\n\n")

	// Environment variable hints
	fmt.Printf("💡 Environment Variables:\n")
	fmt.Printf("   RESOLVER_URL=%s\n", resolverURL)
	fmt.Printf("   REGISTRAR_URL=%s\n", registrarURL)
	if apiKey != "" {
		fmt.Printf("   ACCU_API_KEY=****** (set)\n")
	} else {
		fmt.Printf("   ACCU_API_KEY=(not set, optional)\n")
	}
	if idempotencyKey != "" {
		fmt.Printf("   IDEMPOTENCY_KEY=%s\n", idempotencyKey)
	} else {
		fmt.Printf("   IDEMPOTENCY_KEY=(not set, optional)\n")
	}

	fmt.Println()
	fmt.Printf("📚 For REAL mode (live Accumulate network):\n")
	fmt.Printf("   1. Ensure ACC_NODE_URL is set for the services\n")
	fmt.Printf("   2. Fund your lite account with credits (~20 credits for full lifecycle)\n")
	fmt.Printf("   3. Testnet faucet: curl -X POST https://devnet-faucet.accumulate.io/get -d '{\"url\":\"<lite-account-url>\"}'\n")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func printJSON(title string, data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("   Failed to marshal %s: %v\n", title, err)
		return
	}
	fmt.Printf("   %s:\n%s\n", title, string(jsonBytes))
}