package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.com/accumulatenetwork/accumulate/pkg/api/v3"
	"gitlab.com/accumulatenetwork/accumulate/pkg/api/v3/jsonrpc"
	"gitlab.com/accumulatenetwork/accumulate/pkg/build"
	"gitlab.com/accumulatenetwork/accumulate/pkg/types/messaging"
	"gitlab.com/accumulatenetwork/accumulate/pkg/url"
	"gitlab.com/accumulatenetwork/accumulate/protocol"
)

func main() {
	fmt.Println("=== Accumulate DID Hello World ===\n")

	// 1. Load ACC_NODE_URL
	nodeURL := os.Getenv("ACC_NODE_URL")
	if nodeURL == "" {
		log.Fatal("ACC_NODE_URL environment variable is required")
	}
	fmt.Printf("1. Connecting to Accumulate node: %s\n", nodeURL)

	// Create JSON-RPC client
	client := jsonrpc.NewClient(nodeURL)

	// 2. Generate Ed25519 key pair
	fmt.Println("\n2. Generating Ed25519 key pair...")
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}
	fmt.Printf("   Public Key: %x\n", publicKey)

	// 3. Build and submit transactions
	ctx := context.Background()
	adiLabel := "hello.accu"
	adiURL := fmt.Sprintf("acc://%s", adiLabel)
	keyPageURL := fmt.Sprintf("acc://%s/book/1", adiLabel)
	dataAccountURL := fmt.Sprintf("acc://%s/did", adiLabel)

	// Step 3a: Create Identity (ADI)
	fmt.Printf("\n3a. Creating ADI: %s\n", adiURL)

	// Use the builder pattern with signing
	envelope1, err := build.Transaction().
		For(url.MustParse(adiURL)).
		CreateIdentity(adiURL).
		WithKeyBook(keyPageURL).
		WithKey(publicKey, protocol.SignatureTypeED25519).
		SignWith(url.MustParse("acc://ACME")).  // Use faucet for funding
		Version(1).
		Timestamp(build.UnixTimeNow()).
		PrivateKey(privateKey).
		Done()
	if err != nil {
		log.Fatalf("Failed to create identity: %v", err)
	}

	// Submit to network
	submissions1, err := client.Submit(ctx, envelope1, api.SubmitOptions{})
	if err != nil {
		log.Fatalf("Failed to submit identity transaction: %v", err)
	}
	if len(submissions1) == 0 || !submissions1[0].Success {
		log.Fatalf("Identity transaction submission failed")
	}
	fmt.Printf("   Transaction ID: %s\n", extractTxID(envelope1))

	// Wait for transaction to process
	time.Sleep(2 * time.Second)

	// Step 3b: Create Data Account
	fmt.Printf("\n3b. Creating data account: %s\n", dataAccountURL)

	envelope2, err := build.Transaction().
		For(url.MustParse(adiURL)).
		CreateDataAccount(dataAccountURL).
		SignWith(url.MustParse(keyPageURL)).
		Version(1).
		Timestamp(build.UnixTimeNow()).
		PrivateKey(privateKey).
		Done()
	if err != nil {
		log.Fatalf("Failed to create data account: %v", err)
	}

	// Submit to network
	submissions2, err := client.Submit(ctx, envelope2, api.SubmitOptions{})
	if err != nil {
		log.Fatalf("Failed to submit data account transaction: %v", err)
	}
	if len(submissions2) == 0 || !submissions2[0].Success {
		log.Fatalf("Data account transaction submission failed")
	}
	fmt.Printf("   Transaction ID: %s\n", extractTxID(envelope2))

	// Wait for transaction to process
	time.Sleep(2 * time.Second)

	// Step 3c: Write DID Document
	fmt.Printf("\n3c. Writing DID document to: %s\n", dataAccountURL)

	// Create minimal DID document
	didDocument := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/ns/did/v1",
			"https://w3id.org/security/suites/ed25519-2020/v1",
		},
		"id": fmt.Sprintf("did:acc:%s", adiLabel),
		"verificationMethod": []map[string]interface{}{
			{
				"id":           fmt.Sprintf("did:acc:%s#key1", adiLabel),
				"type":         "Ed25519VerificationKey2020",
				"controller":   fmt.Sprintf("did:acc:%s", adiLabel),
				"publicKeyMultibase": fmt.Sprintf("z%x", publicKey), // Simplified multibase encoding
			},
		},
		"authentication": []string{
			fmt.Sprintf("did:acc:%s#key1", adiLabel),
		},
	}

	didDocData, err := json.Marshal(didDocument)
	if err != nil {
		log.Fatalf("Failed to marshal DID document: %v", err)
	}

	envelope3, err := build.Transaction().
		For(url.MustParse(dataAccountURL)).
		WriteData(didDocData).
		SignWith(url.MustParse(keyPageURL)).
		Version(1).
		Timestamp(build.UnixTimeNow()).
		PrivateKey(privateKey).
		Done()
	if err != nil {
		log.Fatalf("Failed to write DID document: %v", err)
	}

	// Submit to network
	submissions3, err := client.Submit(ctx, envelope3, api.SubmitOptions{})
	if err != nil {
		log.Fatalf("Failed to submit DID document transaction: %v", err)
	}
	if len(submissions3) == 0 || !submissions3[0].Success {
		log.Fatalf("DID document transaction submission failed")
	}
	fmt.Printf("   Transaction ID: %s\n", extractTxID(envelope3))

	// 4. Print results
	fmt.Printf("\n=== SUCCESS ===\n")
	fmt.Printf("DID: did:acc:%s\n", adiLabel)
	fmt.Printf("ADI URL: %s\n", adiURL)
	fmt.Printf("Data Account: %s\n", dataAccountURL)
	fmt.Printf("Key Page: %s\n", keyPageURL)
	fmt.Printf("\nTransaction IDs:\n")
	fmt.Printf("  Create Identity: %s\n", extractTxID(envelope1))
	fmt.Printf("  Create Data Account: %s\n", extractTxID(envelope2))
	fmt.Printf("  Write DID Document: %s\n", extractTxID(envelope3))
	fmt.Printf("\nDID Document written to Accumulate blockchain!\n")
}


// extractTxID extracts transaction ID from envelope
func extractTxID(envelope *messaging.Envelope) string {
	if len(envelope.Transaction) > 0 {
		return fmt.Sprintf("%x", envelope.Transaction[0].GetHash())
	}
	return "unknown"
}