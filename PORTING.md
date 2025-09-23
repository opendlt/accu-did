# SDK Porting Guide

This guide explains how to create Accumulate DID SDKs for other programming languages, following the hybrid strategy defined in the project.

## Overview

The Accumulate DID project uses a **hybrid SDK generation strategy** that combines:

1. **OpenAPI-generated wire layer** - HTTP client code generated from merged OpenAPI specifications
2. **Hand-written ergonomic wrapper** - Language-idiomatic interfaces implementing the behavioral contract

This approach provides both automation benefits and language-specific optimization.

## Architecture

### Generated Layer (Wire)
- HTTP request/response handling
- JSON serialization/deserialization
- Basic error handling
- Generated from `sdks/spec/openapi/accdid-sdk.yaml`

### Hand-written Layer (Ergonomic)
- Language-idiomatic error types
- Retry logic with exponential backoff
- Header management (Request-ID, API-Key, Idempotency-Key)
- DID validation
- Deterministic resolution handling
- Canonical deactivation (410) handling

## Behavioral Contract

All SDK implementations must conform to `sdks/spec/contract.yaml`:

### Required Behaviors

#### Timeouts and Retries
```yaml
defaults:
  timeoutMs: 10000
  retries:
    max: 3
    backoff: exponential
    baseMs: 250
    maxMs: 4000
    jitter: true
```

#### Error Mapping
```yaml
errors:
  retryableStatus: [429, 502, 503, 504]
  deactivatedStatus: 410
  notFoundStatus: 404
  mapping:
    404: "notFound"
    410: "goneDeactivated"
    400: "badRequest"
    # ... see contract.yaml for complete mapping
```

#### Headers
```yaml
headers:
  requestId: "X-Request-Id"
  apiKey: "X-API-Key"
  idempotencyKey: "Idempotency-Key"
```

#### DID Validation
```yaml
validation:
  method: "acc"
  requiredPrefix: "did:acc:"
  forbiddenChars: [" ", "\t", "\n", "\r"]
```

### Deterministic Resolution

When multiple entries exist for a DID, order by:
1. **Sequence number** (highest wins)
2. **Timestamp** (latest wins)
3. **SHA-256 hash** (lexicographic tiebreaker)

### Canonical Deactivation

HTTP 410 responses should return tombstone documents:

```json
{
  "@context": ["https://www.w3.org/ns/did/v1"],
  "id": "did:acc:example",
  "deactivated": true,
  "deactivatedAt": "2024-01-01T00:00:00Z"
}
```

## Implementation Steps

### 1. Generate Wire Layer

```bash
# Merge OpenAPI specifications
make sdk-merge-spec

# Generate client code for your language
# Example with openapi-generator:
openapi-generator-cli generate \
  -i sdks/spec/openapi/accdid-sdk.yaml \
  -g python \
  -o sdks/python/generated/
```

### 2. Create Ergonomic Wrapper

#### Directory Structure
```
sdks/{language}/accdid/
├── client.{ext}           # Main client classes
├── errors.{ext}           # Error types and mapping
├── retry.{ext}            # Retry logic implementation
├── validation.{ext}       # DID validation
├── types.{ext}           # Request/response types
├── examples/             # Working examples
└── tests/               # Comprehensive test suite
```

#### Core Components

**Client Classes:**
- `ResolverClient` - DID resolution with 410 handling
- `RegistrarClient` - DID registration with idempotency

**Error Types:**
- `NotFoundError` (404)
- `GoneDeactivatedError` (410)
- `BadRequestError` (4xx)
- `ServerError` (5xx)
- `TimeoutError` (timeouts)
- `NetworkError` (connectivity)

**Retry Logic:**
- Exponential backoff with jitter
- Retryable status codes: [429, 502, 503, 504]
- Context/cancellation support

### 3. Language-Specific Guidelines

#### Python
```python
# Error handling with custom exceptions
try:
    result = resolver.resolve("did:acc:alice")
except GoneDeactivatedError as e:
    print(f"DID deactivated: {e.tombstone}")
except NotFoundError:
    print("DID not found")

# Context manager support
with AccDIDClient(base_url="http://localhost:8080") as client:
    result = client.resolve("did:acc:alice")
```

#### JavaScript/TypeScript
```typescript
// Promise-based with proper error types
try {
  const result = await resolver.resolve("did:acc:alice");
} catch (error) {
  if (error instanceof GoneDeactivatedError) {
    console.log("DID deactivated:", error.tombstone);
  } else if (error instanceof NotFoundError) {
    console.log("DID not found");
  }
}

// Type definitions for all request/response objects
interface ResolutionResult {
  didDocument: any;
  metadata?: Record<string, any>;
  documentMetadata?: Record<string, any>;
}
```

#### Java
```java
// Exception hierarchy
public class AccDIDException extends Exception { }
public class NotFoundError extends AccDIDException { }
public class GoneDeactivatedError extends AccDIDException {
    public DIDDocument getTombstone() { }
}

// Builder pattern for clients
ResolverClient resolver = ResolverClient.builder()
    .baseUrl("http://localhost:8080")
    .apiKey("optional-key")
    .retryPolicy(RetryPolicy.builder()
        .maxAttempts(3)
        .baseDelay(Duration.ofMillis(250))
        .build())
    .build();
```

#### C#
```csharp
// Async/await with proper exception types
try
{
    var result = await resolver.ResolveAsync("did:acc:alice");
}
catch (GoneDeactivatedException ex)
{
    Console.WriteLine($"DID deactivated: {ex.Tombstone}");
}
catch (NotFoundExeception)
{
    Console.WriteLine("DID not found");
}

// Fluent configuration
var client = new ResolverClient()
    .WithBaseUrl("http://localhost:8080")
    .WithApiKey("optional-key")
    .WithTimeout(TimeSpan.FromSeconds(30));
```

### 4. Testing Requirements

#### Required Test Categories

**Unit Tests:**
- All public methods
- Error mapping and retry logic
- DID validation and parsing
- Header injection and request ID generation

**Integration Tests:**
- Mock HTTP services
- Timeout and retry scenarios
- Deactivation (410) handling

**Test Data:**

Valid DIDs:
- `did:acc:alice`
- `did:acc:beastmode.acme`
- `did:acc:alice/documents`

Invalid DIDs:
- `""` (empty)
- `"invalid"` (no prefix)
- `"did:web:example.com"` (wrong method)
- `"did:acc:"` (missing identifier)

Status Codes:
- Retryable: [429, 502, 503, 504]
- Non-retryable: [400, 401, 403, 404, 410, 422]

### 5. Documentation Requirements

Each SDK must include:

#### README.md
- Installation instructions
- Quick start examples
- Configuration options
- Error handling patterns
- Links to API documentation

#### Examples
- Basic resolution with error handling
- Complete DID lifecycle (register → update → deactivate)
- 410 deactivation handling
- Environment variable configuration

#### API Documentation
- Generated from code comments/docstrings
- All public methods documented
- Error types and conditions explained

## Testing Vectors

Use these test vectors for cross-language compatibility:

### Resolution Test Cases

```json
{
  "valid_dids": [
    "did:acc:alice",
    "did:acc:beastmode.acme",
    "did:acc:company.org/documents"
  ],
  "invalid_dids": [
    "",
    "invalid",
    "did:web:example.com",
    "did:acc:",
    "did:acc:/invalid",
    "did:acc:invalid/"
  ],
  "http_responses": {
    "200": {
      "didDocument": {
        "@context": ["https://www.w3.org/ns/did/v1"],
        "id": "did:acc:alice"
      }
    },
    "404": {
      "code": "notFound",
      "message": "DID not found"
    },
    "410": {
      "didDocument": {
        "@context": ["https://www.w3.org/ns/did/v1"],
        "id": "did:acc:deactivated",
        "deactivated": true,
        "deactivatedAt": "2024-01-01T00:00:00Z"
      }
    }
  }
}
```

### Retry Test Scenarios

1. **Service Unavailable** - Return 503 twice, then 200
2. **Rate Limited** - Return 429 with Retry-After header
3. **Network Error** - Connection refused, then success
4. **Timeout** - Request exceeds configured timeout

## Maintenance

### Keeping SDKs in Sync

1. **OpenAPI Changes** - Regenerate wire layer when specs change
2. **Contract Updates** - Update ergonomic layer when contract.yaml changes
3. **Version Management** - All SDKs should use same version from `VERSION` file
4. **Cross-language Testing** - Run test vectors across all SDK implementations

### Release Process

1. Update `VERSION` file
2. Regenerate all OpenAPI wire layers
3. Run cross-language test vectors
4. Update SDK documentation
5. Tag release with version
6. Publish to language-specific package managers

## Reference Implementation

The Go SDK at `sdks/go/accdid/` serves as the reference implementation demonstrating:

- Clean separation between wire and ergonomic layers
- Comprehensive error handling
- Robust retry logic with exponential backoff
- Proper header management
- Complete test coverage
- Production-ready documentation

Study this implementation when creating SDKs for other languages.