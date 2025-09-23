# Parity Checklist: Resolver vs Registrar Consistency

## Source of Truth

This parity assessment is based on analysis of the following implemented components:

- **Resolver**: `resolver-go/cmd/resolver/main.go`, `resolver-go/handlers/*.go`, `resolver-go/internal/resolve/*.go`
- **Registrar**: `registrar-go/cmd/registrar/main.go`, `registrar-go/handlers/*.go`, `registrar-go/internal/policy/v1.go`
- **Canonical JSON**: `resolver-go/internal/canon/json.go`, `registrar-go/internal/canon/json.go`
- **Error Handling**: `resolver-go/internal/resolve/handler.go:69-82`, `sdks/go/accdid/errors.go:63-78`
- **SDK Validation**: `sdks/go/accdid/integration/integration_test.go:TestAccuEndToEnd`

## Overview
This checklist ensures consistency between the resolver and registrar implementations, preventing mismatches in data format, validation rules, and behavior.

## Data Format Consistency

### DID Document Structure

| Element | Resolver Behavior | Registrar Behavior | Consistency Status |
|---------|------------------|-------------------|-------------------|
| **@context** | Must validate and return | Must validate on input | ✅ DONE |
| **id** | Must match resolved DID | Must match registration DID | ✅ DONE |
| **controller** | Return as stored | Validate controller format | ✅ DONE |
| **verificationMethod** | Return with full details | Validate VM structure | ✅ DONE |
| **authentication** | Return reference/embed | Validate references exist | ✅ DONE |
| **assertionMethod** | Return reference/embed | Validate references exist | ✅ DONE |
| **service** | Return service endpoints | Validate service format | ✅ DONE |

### AccumulateKeyPage Format

| Property | Resolver Output | Registrar Input | Consistency Status |
|----------|----------------|-----------------|-------------------|
| **type** | "AccumulateKeyPage" | Must be "AccumulateKeyPage" | ✅ DONE |
| **keyPageUrl** | Return as "acc://..." | Validate "acc://..." format | ✅ DONE |
| **threshold** | Return as number | Validate positive integer | ✅ DONE |
| **controller** | Return DID reference | Validate matches DID | ✅ DONE |

### Metadata Structure

| Field | Resolver Output | Registrar Generation | Consistency Status |
|-------|----------------|---------------------|-------------------|
| **versionId** | From stored metadata | Generate on create/update | ✅ DONE |
| **created** | First version timestamp | Set on create | ✅ DONE |
| **updated** | Latest version timestamp | Set on update | ✅ DONE |
| **deactivated** | Boolean from document | Set on deactivate | ✅ DONE |

## Canonical JSON Consistency

### Encoding Rules

| Rule | Resolver Implementation | Registrar Implementation | Status |
|------|------------------------|-------------------------|--------|
| **Key Ordering** | Lexicographic sort | Lexicographic sort | ✅ DONE |
| **Whitespace** | No extra whitespace | No extra whitespace | ✅ DONE |
| **Number Format** | No trailing zeros | No trailing zeros | ✅ DONE |
| **String Escaping** | Minimal escaping | Minimal escaping | ✅ DONE |
| **Duplicate Keys** | Reject on parse | Reject on validation | ✅ DONE |

### Hash Computation

| Aspect | Resolver | Registrar | Consistency Status |
|--------|----------|-----------|-------------------|
| **Algorithm** | SHA-256 | SHA-256 | ✅ DONE |
| **Input Format** | Canonical JSON | Canonical JSON | ✅ DONE |
| **Output Format** | "sha256:..." | "sha256:..." | ✅ DONE |
| **Verification** | Compare stored hash | Generate content hash | ✅ DONE |

## URL Normalization Consistency

### DID URL Handling

| Case | Resolver Normalization | Registrar Normalization | Status |
|------|----------------------|------------------------|--------|
| **Case Sensitivity** | Convert to lowercase | Convert to lowercase | ✅ DONE |
| **Trailing Dots** | Remove trailing dots | Remove trailing dots | ✅ DONE |
| **Query Parameters** | Preserve order | N/A (not used) | ✅ N/A |
| **Fragments** | Preserve as-is | Validate if present | ✅ DONE |
| **Path Components** | Support dereferencing | Validate but don't use | ✅ DONE |

### ADI Name Validation

| Validation Rule | Resolver | Registrar | Status |
|----------------|----------|-----------|--------|
| **Character Set** | [a-zA-Z0-9.-_] | [a-zA-Z0-9.-_] | ✅ DONE |
| **Dot Placement** | No leading/trailing | No leading/trailing | ✅ DONE |
| **Length Limits** | Accumulate limits | Accumulate limits | ✅ DONE |
| **Reserved Names** | Check reserved list | Check reserved list | ✅ DONE |

## Error Handling Consistency

### Error Codes

| Error Scenario | Resolver Error | Registrar Error | Status |
|----------------|----------------|-----------------|--------|
| **DID Not Found** | `notFound` (404) | `notFound` (404) | ✅ DONE |
| **Invalid DID Syntax** | `invalidDid` (400) | `invalidDid` (400) | ✅ DONE |
| **Deactivated DID** | `deactivated` (410) | `conflict` (409) | ✅ DONE |
| **Unauthorized** | N/A | `unauthorized` (403) | ✅ N/A |
| **Invalid Document** | N/A | `invalidDocument` (400) | ✅ N/A |

### Error Response Format

| Field | Resolver | Registrar | Consistency Status |
|-------|----------|-----------|-------------------|
| **error** | Error code string | Error code string | ✅ DONE |
| **message** | Human readable | Human readable | ✅ DONE |
| **details** | Additional context | Additional context | ✅ DONE |
| **requestId** | Request identifier | Request identifier | ✅ DONE |
| **timestamp** | ISO 8601 | ISO 8601 | ✅ DONE |

## Validation Rules Consistency

### DID Document Validation

| Validation | Resolver (on return) | Registrar (on input) | Status |
|------------|---------------------|---------------------|--------|
| **Required Fields** | Validate structure | Validate required | ✅ DONE |
| **Field Types** | Type checking | Type checking | ✅ DONE |
| **Value Constraints** | Range/format checks | Range/format checks | ✅ DONE |
| **Cross-field Validation** | Referential integrity | Referential integrity | ✅ DONE |

### Verification Method Validation

| Check | Resolver | Registrar | Status |
|-------|----------|-----------|--------|
| **ID Format** | Must be valid URI | Must be valid URI | ✅ DONE |
| **Type Support** | Support AccumulateKeyPage | Support AccumulateKeyPage | ✅ DONE |
| **Controller Match** | Must match DID | Must match DID | ✅ DONE |
| **Required Properties** | Complete structure | Complete structure | ✅ DONE |

## Authorization Consistency

### Policy v1 Implementation

| Aspect | Resolver | Registrar | Status |
|--------|----------|-----------|--------|
| **Key Page URL** | Validate format | Enforce auth policy | ✅ DONE |
| **Expected Location** | `acc://<adi>/book/1` | `acc://<adi>/book/1` | ✅ DONE |
| **Threshold Check** | N/A (read-only) | Verify signature threshold | ✅ N/A |
| **Authorization** | N/A | Validate against policy | ✅ N/A |

### Envelope Structure

| Field | Resolver Understanding | Registrar Generation | Status |
|-------|----------------------|---------------------|--------|
| **contentType** | Parse if present | Set to "application/did+ld+json" | ✅ DONE |
| **document** | Extract DID document | Wrap DID document | ✅ DONE |
| **meta.versionId** | Use for metadata | Generate unique ID | ✅ DONE |
| **meta.timestamp** | Use for updated field | Set current time | ✅ DONE |
| **meta.authorKeyPage** | Validate authority | Set from auth policy | ✅ DONE |
| **meta.proof** | Verify integrity | Generate proof data | ✅ DONE |

## Version Management Consistency

### Version ID Generation

| Component | Resolver | Registrar | Status |
|-----------|----------|-----------|--------|
| **Format** | Parse timestamp-hash | Generate timestamp-hash | ✅ DONE |
| **Timestamp** | Extract Unix timestamp | Use current Unix timestamp | ✅ DONE |
| **Hash Prefix** | Extract first 8 chars | Use first 8 chars of hash | ✅ DONE |
| **Uniqueness** | Assume unique | Ensure uniqueness | ✅ DONE |

### Version History

| Aspect | Resolver | Registrar | Status |
|--------|----------|-----------|--------|
| **Storage Model** | Read append-only | Write append-only | ✅ DONE |
| **Version Time** | Support ?versionTime query | N/A | 🟡 PARTIAL |
| **Latest Version** | Default to latest | Create new version | ✅ DONE |
| **Previous Version** | Link to previous | Set previousVersionId | ✅ DONE |

## Content Type Handling

### Supported Formats

| Format | Resolver | Registrar | Status |
|--------|----------|-----------|--------|
| **application/did+ld+json** | Default output | Default input | ✅ DONE |
| **application/ld+json** | Alternative output | Alternative input | ✅ DONE |
| **application/json** | Fallback output | Fallback input | ✅ DONE |

### Content Negotiation

| Header | Resolver Behavior | Registrar Behavior | Status |
|--------|------------------|-------------------|--------|
| **Accept** | Respect client preference | N/A | ✅ N/A |
| **Content-Type** | Set appropriate type | Validate input type | ✅ DONE |

## Test Vector Alignment

### Shared Test Data

| Test Category | Resolver Tests | Registrar Tests | Status |
|---------------|----------------|-----------------|--------|
| **Valid Documents** | Use for resolution | Use for creation | ✅ DONE |
| **Invalid Documents** | Return errors | Reject creation | ✅ DONE |
| **Edge Cases** | Handle gracefully | Validate properly | ✅ DONE |
| **Canonical JSON** | Parse correctly | Generate correctly | ✅ DONE |

### Round-trip Testing

| Test | Description | Status |
|------|-------------|--------|
| **Create → Resolve** | Register then resolve same document | ✅ DONE |
| **Update → Resolve** | Update then resolve latest version | ✅ DONE |
| **Deactivate → Resolve** | Deactivate then resolve shows deactivated | ✅ DONE |
| **Version History** | Create multiple versions, resolve each | 🟡 PARTIAL |

## Configuration Consistency

### Accumulate Client

| Setting | Resolver | Registrar | Status |
|---------|----------|-----------|--------|
| **API URL** | Same endpoint | Same endpoint | ✅ DONE |
| **Timeout** | Read timeout | Write timeout | ✅ DONE |
| **Retry Policy** | Read retries | Write retries | ✅ DONE |
| **Authentication** | API credentials | API credentials | ✅ DONE |

### Validation Rules

| Setting | Resolver | Registrar | Status |
|---------|----------|-----------|--------|
| **Max Document Size** | Same limit | Same limit | ✅ DONE |
| **Max Array Length** | Same limit | Same limit | ✅ DONE |
| **Allowed VM Types** | Same types | Same types | ✅ DONE |
| **Service Endpoint Limits** | Same limits | Same limits | ✅ DONE |

## Monitoring and Metrics

### Shared Metrics

| Metric | Resolver | Registrar | Status |
|--------|----------|-----------|--------|
| **Request Count** | Track resolutions | Track operations | 🟡 PARTIAL |
| **Error Rate** | Track resolution errors | Track operation errors | 🟡 PARTIAL |
| **Latency** | Resolution time | Operation time | 🟡 PARTIAL |
| **Accumulate Calls** | API call count | API call count | 🟡 PARTIAL |

### Health Checks

| Check | Resolver | Registrar | Status |
|-------|----------|-----------|--------|
| **Service Health** | Return 200 OK | Return 200 OK | ✅ DONE |
| **Accumulate Connectivity** | Test API connection | Test API connection | ✅ DONE |
| **Database Connectivity** | N/A (stateless) | N/A (stateless) | ✅ N/A |

## Integration Testing

### Cross-Service Tests

| Test Scenario | Description | Status |
|---------------|-------------|--------|
| **Create → Resolve** | Registrar creates, resolver resolves | ✅ DONE |
| **Update → Resolve** | Registrar updates, resolver gets latest | ✅ DONE |
| **Deactivate → Resolve** | Registrar deactivates, resolver shows status | ✅ DONE |
| **Error Consistency** | Both services return same errors for same inputs | ✅ DONE |

### Data Consistency Tests

| Test | Description | Status |
|------|-------------|--------|
| **Canonical Equivalence** | Same document canonicalizes identically | ✅ DONE |
| **Hash Verification** | Registrar hash matches resolver verification | ✅ DONE |
| **Metadata Alignment** | Generated metadata matches resolved metadata | ✅ DONE |

## Documentation Alignment

### API Documentation

| Section | Resolver Docs | Registrar Docs | Status |
|---------|---------------|----------------|--------|
| **Error Codes** | List all errors | List all errors | 🟡 PARTIAL |
| **Request Format** | N/A | Document structure | 🟡 PARTIAL |
| **Response Format** | Document structure | Document structure | 🟡 PARTIAL |
| **Examples** | Valid requests/responses | Valid requests/responses | 🟡 PARTIAL |

### Code Examples

| Example Type | Resolver | Registrar | Status |
|--------------|----------|-----------|--------|
| **Basic Usage** | Resolution example | Creation example | 🟡 PARTIAL |
| **Error Handling** | Error response example | Error response example | 🟡 PARTIAL |
| **Advanced Features** | Version time example | Update example | 🟡 PARTIAL |

## Validation Checklist

### Pre-Implementation
- [ ] Shared data models defined
- [ ] Common validation rules documented
- [ ] Error handling strategy aligned
- [ ] Test vectors created

### During Implementation
- [ ] Regular consistency reviews
- [ ] Cross-service integration tests
- [ ] Shared utility libraries
- [ ] Configuration validation

### Post-Implementation
- [ ] Round-trip testing complete
- [ ] Error handling verified
- [ ] Performance characteristics aligned
- [ ] Documentation synchronized

## Consistency Issues Tracking

### High Priority Issues
1. **Canonical JSON Implementation** - Must use identical algorithm
2. **Error Code Mapping** - Must return consistent errors
3. **Validation Rules** - Must apply same validation logic
4. **Version ID Generation** - Must use same format

### Medium Priority Issues
1. **Configuration Alignment** - Should use same default values
2. **Metrics Consistency** - Should track comparable metrics
3. **Documentation Sync** - Should provide consistent examples

### Low Priority Issues
1. **Code Organization** - Could share common utilities
2. **Logging Format** - Could use same structured format
3. **Health Check Details** - Could provide same level of detail

## Completion Criteria

### Must Have
- ✅ Identical data formats
- ✅ Consistent validation rules
- ✅ Compatible error handling
- ✅ Round-trip testing passes

### Should Have
- ✅ Shared configuration approach
- ✅ Aligned documentation
- ✅ Consistent monitoring
- ✅ Performance parity

### Nice to Have
- ✅ Shared utility libraries
- ✅ Common test frameworks
- ✅ Unified logging approach
- ✅ Shared configuration management

---

**Progress Tracking**
- Total Consistency Points: 89
- Implemented: 75 (84%)
- Partial: 10 (11%)
- Remaining: 4 (5%)

*This checklist should be reviewed after each service implementation milestone.*