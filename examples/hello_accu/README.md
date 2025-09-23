# Accumulate DID Hello World

This example demonstrates the complete zero-to-hero DID lifecycle on a real Accumulate blockchain: from funding a lite account to creating an ADI, data account, and writing a DID document.

## Prerequisites

1. **Go 1.22+** installed
2. **Local Accumulate devnet** running (use `make devnet-up` in the repo root)

## Quick Start (One Command)

From the repository root:

```bash
# Start devnet + services + run SDK example
make services-up && make sdk-example

# Clean up when done
make services-down && make devnet-down
```

## Manual Setup

### Step 1: Start Local Devnet

```bash
# From repo root
make devnet-up
```

This starts a local Accumulate devnet with faucet support on `http://127.0.0.1:26656`.

### Step 2: Set Environment

```bash
export ACC_NODE_URL=http://127.0.0.1:26656
```

### Step 3: Run Example

```bash
cd examples/hello_accu
go run main.go
```

### Step 4: Clean Up

```bash
# From repo root
make devnet-down
```

## Sample Output

```
=== Accumulate DID Hello World ===

1. Connecting to Accumulate node: http://127.0.0.1:26656

2. Generating Ed25519 key pair...
   Public Key: 8c7e8b4f2d1a9c5e3b7f6a8d9e2c4f1b5a7c8e9f0d2b4c6e8a1d3f5c7e9b2a4d

3. Creating lite account for funding...
   Lite Account: acc://1234567890abcdef1234567890abcdef12345678

4. Funding lite account from faucet...
   Trying faucet endpoint: http://127.0.0.1:26656/faucet
   Faucet request successful via http://127.0.0.1:26656/faucet
   ✅ Lite account funded successfully

5a. Creating ADI: acc://hello.accu
   Transaction ID: a1b2c3d4e5f6789abcdef0123456789abcdef0123456789abcdef0123456789

5b. Creating data account: acc://hello.accu/did
   Transaction ID: b2c3d4e5f6789abcdef0123456789abcdef0123456789abcdef0123456789abc

5c. Writing DID document to: acc://hello.accu/did
   Transaction ID: c3d4e5f6789abcdef0123456789abcdef0123456789abcdef0123456789abcd

=== SUCCESS ===
DID: did:acc:hello.accu
ADI URL: acc://hello.accu
Data Account: acc://hello.accu/did
Key Page: acc://hello.accu/book/1
Lite Account: acc://1234567890abcdef1234567890abcdef12345678

Transaction IDs:
  Create Identity: a1b2c3d4e5f6789abcdef0123456789abcdef0123456789abcdef0123456789
  Create Data Account: b2c3d4e5f6789abcdef0123456789abcdef0123456789abcdef0123456789abc
  Write DID Document: c3d4e5f6789abcdef0123456789abcdef0123456789abcdef0123456789abcd

DID Document written to Accumulate blockchain!
You can now resolve this DID using: did:acc:hello.accu
```

## What It Does

1. **Connects** to Accumulate devnet via JSON-RPC
2. **Generates** an Ed25519 key pair for signing
3. **Creates lite account** from the public key for funding
4. **Funds lite account** using devnet faucet (auto-funding for development)
5. **Creates ADI** (`acc://hello.accu`) - the identity container, funded by lite account
6. **Creates data account** (`acc://hello.accu/did`) - where DID document is stored
7. **Writes DID document** as JSON data to the data account
8. **Prints transaction IDs** for verification on Accumulate explorer

## DID Resolution

After running, you can resolve the DID using the resolver service:

```bash
# Start resolver (if not already running)
cd ../../resolver-go && ACC_NODE_URL=http://127.0.0.1:26656 go run cmd/resolver/main.go --addr :8080 --real &

# Resolve the DID
curl "http://localhost:8080/resolve?did=did:acc:hello.accu"
```

Or query the data account directly:

```bash
# Direct Accumulate query to data account
curl -X POST http://127.0.0.1:26656 -d '{"jsonrpc":"2.0","method":"query","params":{"url":"acc://hello.accu/did"},"id":1}'
```

## Features

- **Zero-to-hero workflow**: Complete DID lifecycle from funding to resolution
- **Real blockchain transactions**: Uses live Accumulate devnet (not mocks)
- **Automatic faucet funding**: No manual token acquisition needed
- **Ed25519 signatures**: Industry-standard cryptography
- **W3C DID Core compliance**: Standards-compliant DID document structure
- **Accumulate-native**: Uses native `acc://` URL scheme and protocols