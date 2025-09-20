# DID Method Specification: `did:acc` v0.1

## Abstract
This specification defines the `did:acc` DID method for the Accumulate blockchain. The method provides decentralized identifier functionality leveraging Accumulate's native ADI (Accumulate Digital Identifier) system, data accounts, and key management infrastructure. DID documents are stored as append-only entries in ADI data accounts, with operations authorized by Accumulate Key Pages and validated through blockchain transactions.

## Status
- **Version**: 0.1
- **Status**: Draft
- **Last Updated**: 2024
- **Authors**: Accumulate DID Working Group
- **Specification URI**: https://docs.accumulate.defi/did/spec/v0.1

## DID Method Name
The name string that shall identify this DID method is: `acc`

A DID that uses this method MUST begin with the prefix `did:acc:`. Per the DID Core specification, this prefix MUST be in lowercase.

## DID Syntax

### Method Identifier
The DID method name is `acc`, representing the Accumulate blockchain network.

### DID Format
```abnf
did-acc         = "did:acc:" acc-specific-id
acc-specific-id = adi-name
adi-name        = 1*idchar *("." 1*idchar)
idchar          = ALPHA / DIGIT / "-" / "_"
```

### DID URL Format
```abnf
did-acc-url = did-acc [did-params] [did-path] [did-query] [did-fragment]
did-params  = ";" param *(";" param)
did-path    = "/" path-absolute
did-query   = "?" query
did-fragment = "#" fragment
param       = param-name "=" param-value
```

### Examples

#### Basic DID
```
did:acc:alice
did:acc:beastmode.acme
did:acc:org-name.division.team
```

#### Case-Folding of ADI
ADI names are case-insensitive and normalized to lowercase:
```
did:acc:ALICE        → did:acc:alice
did:acc:BeastMode.ACME → did:acc:beastmode.acme
```

#### DID URL with Fragment (Verification Method)
```
did:acc:alice#key-1                    # Primary key
did:acc:beastmode.acme#signing-key     # Signing key
did:acc:alice#authentication          # Authentication method
```

#### DID URL with versionTime Parameter
```
did:acc:alice?versionTime=2024-01-01T00:00:00Z
did:acc:alice?versionTime=1704067200
did:acc:alice?versionTime=2024-01-01T12:30:45.123Z
```

#### Complex DID URLs
```
did:acc:alice;service=messaging/inbox?format=json#key-1
did:acc:alice/documents/credential-123?versionTime=2024-01-01T00:00:00Z
```

## State Model

### DID Document Storage Location
DID documents are stored in Accumulate data accounts at a canonical location:
```
acc://<adi>/data/did
```

Where `<adi>` corresponds to the ADI name in the DID. For example:
- DID: `did:acc:alice` → Storage: `acc://alice/data/did`
- DID: `did:acc:beastmode.acme` → Storage: `acc://beastmode.acme/data/did`

### Append-Only Entry Model
DID documents are stored as append-only entries in the data account using Accumulate's `writeData` transaction. Each entry represents a version of the DID document, with the latest valid entry (authorized by the appropriate Key Page at write-time) being the canonical version.

### Entry Envelope Structure
Each DID document entry is wrapped in an envelope with the following structure:

```json
{
  "contentType": "application/did+json",
  "document": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:alice",
    "verificationMethod": [...]
  },
  "meta": {
    "versionId": "1704067200-8b4c4f7b",
    "previousVersionId": "1704067100-7a9b8c6d",
    "timestamp": "2024-01-01T00:00:00Z",
    "authorKeyPage": "acc://alice/book/1",
    "proof": {
      "txid": "0x1234567890abcdef...",
      "contentHash": "sha256:abc123..."
    }
  }
}
```

#### Envelope Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `contentType` | string | ✅ | MIME type, always "application/did+json" |
| `document` | object | ✅ | The DID document content |
| `meta.versionId` | string | ✅ | Unique version identifier (timestamp-hash) |
| `meta.previousVersionId` | string | ❌ | Previous version identifier (for updates) |
| `meta.timestamp` | string | ✅ | ISO 8601 timestamp of entry creation |
| `meta.authorKeyPage` | string | ✅ | Key Page URL that authorized this entry |
| `meta.proof.txid` | string | ✅ | Accumulate transaction ID |
| `meta.proof.contentHash` | string | ✅ | SHA-256 hash of canonical document |

### Version Selection Rules
1. **Latest Valid Entry**: By default, return the most recent entry that was properly authorized
2. **Version Time Query**: If `versionTime` parameter is provided, return the latest entry with `timestamp ≤ versionTime`
3. **Deactivation**: If the latest entry has `deactivated: true`, the DID is considered deactivated
4. **Authorization Validation**: Only entries authorized by `acc://<adi>/book/1` at write-time are considered valid

## Operations

All DID operations result in writing new entries to the data account via Accumulate's `writeData` transaction. The operations are:

### Create Operation
Creates the initial DID document entry.

**Process:**
1. **Prepare DID Document** with initial verification methods and services
2. **Create Envelope** with `contentType: "application/did+json"`
3. **Execute writeData** to `acc://<adi>/data/did`
4. **Authorization** via ADI's Key Page at `acc://<adi>/book/1`

**Example Envelope:**
```json
{
  "contentType": "application/did+json",
  "document": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:alice",
    "verificationMethod": [{
      "id": "did:acc:alice#key-1",
      "type": "AccumulateKeyPage",
      "controller": "did:acc:alice",
      "keyPageUrl": "acc://alice/book/1",
      "threshold": 1
    }],
    "authentication": ["did:acc:alice#key-1"]
  },
  "meta": {
    "versionId": "1704067200-8b4c4f7b",
    "timestamp": "2024-01-01T00:00:00Z",
    "authorKeyPage": "acc://alice/book/1"
  }
}
```

### Update Operation
Modifies the DID document by appending a new entry.

**Process:**
1. **Prepare Updated Document** with changes (services, verification methods, etc.)
2. **Create Envelope** with `previousVersionId` referencing the current version
3. **Execute writeData** to append new entry
4. **Authorization** via the same Key Page

**Typical Updates:**
- Add/remove service endpoints
- Add/remove verification methods (non-key-page)
- Update service endpoint URLs
- Modify verification method properties

### Deactivate Operation
Marks the DID as deactivated by writing a terminal entry.

**Process:**
1. **Create Deactivation Document** with `deactivated: true`
2. **Execute writeData** to append final entry
3. **Authorization** via Key Page (last authorized operation)

**Deactivation Envelope:**
```json
{
  "contentType": "application/did+json",
  "document": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:alice",
    "deactivated": true
  },
  "meta": {
    "versionId": "1704326400-final",
    "previousVersionId": "1704153600-3f4e5d6c",
    "timestamp": "2024-01-04T00:00:00Z",
    "authorKeyPage": "acc://alice/book/1"
  }
}
```

### Key Rotation and Delegation
Key management is handled separately through Accumulate's Key Page system:

**Key Rotation:**
- Execute `updateKeyPage` transaction on `acc://<adi>/book/1`
- DID document often remains unchanged
- New keys automatically apply to subsequent DID operations
- No DID document update required unless verification method properties change

**Key Delegation:**
- Create new Key Pages (e.g., `acc://<adi>/book/2`)
- Add verification methods in DID document referencing new Key Pages
- Maintain authorization hierarchy through Key Page management

## Resolution Algorithm

The DID resolution process follows these steps to retrieve and validate a DID document:

### Step 1: DID Normalization
1. **Extract ADI**: Remove `did:acc:` prefix to get ADI name
2. **Case Normalization**: Convert ADI to lowercase
3. **Trailing Dot Removal**: Remove any trailing dots
4. **Validation**: Ensure ADI name contains valid characters

**Example:**
```
did:acc:BeastMode.ACME. → adi: beastmode.acme
```

### Step 2: Locate Data Account
1. **Construct Data Account URL**: `acc://<normalized-adi>/data/did`
2. **Query Accumulate Network**: Retrieve all entries from the data account
3. **Handle Not Found**: Return `notFound` error if data account doesn't exist

### Step 3: Select Entry Version
1. **Parse versionTime Parameter**: If provided, convert to timestamp
2. **Filter Valid Entries**: Only consider entries authorized by `acc://<adi>/book/1`
3. **Apply Version Logic**:
   - **Default**: Latest valid entry
   - **With versionTime**: Latest valid entry where `timestamp ≤ versionTime`
4. **Handle No Valid Entries**: Return `notFound` if no entries meet criteria

### Step 4: Validate Entry
For the selected entry, perform comprehensive validation:

1. **Transaction Verification**:
   - Verify `meta.proof.txid` exists on Accumulate blockchain
   - Confirm transaction was properly committed
   - Check transaction was authorized by `meta.authorKeyPage`

2. **Content Hash Verification**:
   - Compute canonical JSON of `document` field
   - Calculate SHA-256 hash of canonical representation
   - Verify `meta.proof.contentHash` matches computed hash

3. **Authorization Verification**:
   - Confirm `meta.authorKeyPage` equals `acc://<adi>/book/1`
   - Verify Key Page had proper authority at transaction time
   - Check signature threshold was met for the transaction

### Step 5: Return Resolution Result
Structure the response according to W3C DID Core specification:

```json
{
  "didDocument": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:alice",
    "verificationMethod": [...]
  },
  "didDocumentMetadata": {
    "versionId": "1704067200-8b4c4f7b",
    "created": "2024-01-01T00:00:00Z",
    "updated": "2024-01-02T12:30:00Z",
    "deactivated": false,
    "nextUpdate": null,
    "nextVersionId": null
  },
  "didResolutionMetadata": {
    "contentType": "application/did+ld+json",
    "retrieved": "2024-01-05T10:15:30Z",
    "pattern": "^did:acc:",
    "driverUrl": "https://resolver.accumulate.defi",
    "duration": 150
  }
}
```

### Error Conditions

| Error | HTTP Status | Condition |
|-------|-------------|-----------|
| `notFound` | 404 | DID does not exist or no valid entries |
| `deactivated` | 410 | Latest valid entry has `deactivated: true` |
| `invalidDid` | 400 | DID syntax is malformed |
| `versionNotFound` | 404 | No entry exists for specified `versionTime` |
| `authorizationFailed` | 500 | Entry authorization could not be verified |
| `contentHashMismatch` | 500 | Content hash validation failed |
| `internalError` | 500 | Accumulate network or other system error |

## Verification Methods

### AccumulateKeyPage

A custom verification method type that references Accumulate Key Pages:

```json
{
  "id": "did:acc:alice#key-1",
  "type": "AccumulateKeyPage",
  "controller": "did:acc:alice",
  "keyPageUrl": "acc://alice/book/1",
  "threshold": 1
}
```

**Properties:**
- `keyPageUrl`: URL of the Accumulate Key Page
- `threshold`: Required signature threshold for operations
- `controller`: MUST match the DID being resolved

**Validation Rules:**
- Key Page URL MUST use `acc://` protocol
- Key Page MUST exist on the Accumulate network
- Threshold MUST be a positive integer
- Controller MUST be the DID being resolved or another DID controlled by the same ADI

## Security Considerations

### Authorization
- All DID operations require signatures from the ADI's authorized Key Page
- Threshold requirements enforced by Accumulate protocol
- Key Page updates don't require DID document changes

### Content Integrity
- Each DID document version includes SHA-256 content hash
- Hashes are verified during resolution
- Accumulate's blockchain provides tamper-evidence

### Privacy
- DID documents are publicly readable
- Avoid storing PII in DID documents
- Use service endpoints for private data references

## Normalization

### DID Normalization
- ADI names are case-insensitive
- Normalize to lowercase for comparison
- Remove trailing dots from ADI names

### URL Normalization
- Preserve query parameters and fragments
- URL-decode encoded characters
- Maintain parameter order

## Error Handling

### Resolution Errors

| Error Code | Description |
|------------|-------------|
| `notFound` | DID does not exist |
| `deactivated` | DID has been deactivated |
| `invalidDid` | DID syntax is invalid |
| `versionNotFound` | Requested version doesn't exist |

### Registration Errors

| Error Code | Description |
|------------|-------------|
| `unauthorized` | Insufficient authorization |
| `invalidDocument` | DID document validation failed |
| `conflict` | DID already exists |
| `thresholdNotMet` | Key page threshold not met |

## Implementation Requirements

### Resolver Requirements
1. Support DID URL dereferencing
2. Implement versionTime parameter
3. Return proper metadata
4. Handle deactivated DIDs
5. Validate content hashes

### Registrar Requirements
1. Validate DID document schema
2. Enforce authorization policies
3. Generate version identifiers
4. Compute content hashes
5. Support batch operations

## Test Vectors

See `spec/vectors/` for comprehensive test cases including:
- URL normalization tests
- Document validation tests
- Authorization scenarios
- Error conditions

## References

- [W3C DID Core Specification](https://www.w3.org/TR/did-core/)
- [Accumulate Protocol Documentation](https://docs.accumulate.defi/)
- [DID Resolution Specification](https://w3c-ccg.github.io/did-resolution/)
- [DID Registration](https://identity.foundation/did-registration/)

## Appendix A: Example DID Document

```json
{
  "@context": [
    "https://www.w3.org/ns/did/v1",
    "https://w3id.org/security/suites/jws-2020/v1"
  ],
  "id": "did:acc:beastmode.acme",
  "controller": "did:acc:beastmode.acme",
  "verificationMethod": [
    {
      "id": "did:acc:beastmode.acme#key-1",
      "type": "AccumulateKeyPage",
      "controller": "did:acc:beastmode.acme",
      "keyPageUrl": "acc://beastmode.acme/book/1",
      "threshold": 2
    }
  ],
  "authentication": ["did:acc:beastmode.acme#key-1"],
  "assertionMethod": ["did:acc:beastmode.acme#key-1"],
  "service": [
    {
      "id": "did:acc:beastmode.acme#resolver",
      "type": "DIDResolver",
      "serviceEndpoint": "https://resolver.accumulate.defi"
    }
  ]
}
```

## Appendix B: Registration Request Example

```json
{
  "did": "did:acc:beastmode.acme",
  "options": {
    "clientSecretMode": false
  },
  "secret": {},
  "didDocument": {
    "@context": ["https://www.w3.org/ns/did/v1"],
    "id": "did:acc:beastmode.acme",
    "verificationMethod": [
      {
        "id": "did:acc:beastmode.acme#key-1",
        "type": "AccumulateKeyPage",
        "controller": "did:acc:beastmode.acme",
        "keyPageUrl": "acc://beastmode.acme/book/1",
        "threshold": 2
      }
    ]
  }
}
```

## Change Log

### Version 0.1 (2024)
- Initial draft specification
- Define DID syntax and operations
- Specify AccumulateKeyPage verification method
- Add normalization and error handling

---

*This specification is subject to change as the implementation evolves.*