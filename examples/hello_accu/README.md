# Accumulate DID Hello World

This example demonstrates creating a complete DID on the Accumulate blockchain using the real Accumulate API.

## Prerequisites

1. **Go 1.22+** installed
2. **Accumulate node** running (local devnet or testnet access)
3. **Environment variable** `ACC_NODE_URL` set to your Accumulate node

## Setup

```bash
# Set your Accumulate node URL
export ACC_NODE_URL=http://localhost:26657  # Local devnet
# OR
export ACC_NODE_URL=https://testnet.accumulatenetwork.io  # Testnet
```

## Run

```bash
cd examples/hello_accu
go mod init hello_accu
go mod tidy
go run main.go
```

## Sample Output

```
=== Accumulate DID Hello World ===

1. Connecting to Accumulate node: http://localhost:26657

2. Generating Ed25519 key pair...
   Public Key: 8c7e8b4f2d1a9c5e3b7f6a8d9e2c4f1b5a7c8e9f0d2b4c6e8a1d3f5c7e9b2a4d

3a. Creating ADI: acc://hello.accu
   Transaction ID: a1b2c3d4e5f6789abcdef0123456789abcdef0123456789abcdef0123456789

3b. Creating data account: acc://hello.accu/did
   Transaction ID: b2c3d4e5f6789abcdef0123456789abcdef0123456789abcdef0123456789abc

3c. Writing DID document to: acc://hello.accu/did
   Transaction ID: c3d4e5f6789abcdef0123456789abcdef0123456789abcdef0123456789abcd

=== SUCCESS ===
DID: did:acc:hello.accu
ADI URL: acc://hello.accu
Data Account: acc://hello.accu/did
Key Page: acc://hello.accu/book/1

Transaction IDs:
  Create Identity: a1b2c3d4e5f6789abcdef0123456789abcdef0123456789abcdef0123456789
  Create Data Account: b2c3d4e5f6789abcdef0123456789abcdef0123456789abcdef0123456789abc
  Write DID Document: c3d4e5f6789abcdef0123456789abcdef0123456789abcdef0123456789abcd

DID Document written to Accumulate blockchain!
```

## What It Does

1. **Connects** to Accumulate node via JSON-RPC
2. **Generates** an Ed25519 key pair for signing
3. **Creates ADI** (`acc://hello.accu`) - the identity container
4. **Creates data account** (`acc://hello.accu/did`) - where DID document is stored
5. **Writes DID document** as JSON data to the data account
6. **Prints transaction IDs** for verification on Accumulate explorer

## DID Resolution

After running, you can resolve the DID using:
- Accumulate DID Resolver: `did:acc:hello.accu`
- Direct data account query: `acc://hello.accu/did`

## Notes

- Uses **real Accumulate transactions** (costs small fees on mainnet)
- **Ed25519 signatures** for all operations
- **W3C DID Core** compliant document structure
- **Accumulate-native** URL scheme (`acc://`)