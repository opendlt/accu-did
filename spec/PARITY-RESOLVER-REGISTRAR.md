# Parity Checklist: Resolver vs Registrar Consistency

## Overview
This checklist ensures consistency between the resolver and registrar implementations, preventing mismatches in data format, validation rules, and behavior.

## Data Format Consistency

### DID Document Structure

| Element | Resolver Behavior | Registrar Behavior | Consistency Status |
|---------|------------------|-------------------|-------------------|
| **@context** | Must validate and return | Must validate on input | ❌ TODO |
| **id** | Must match resolved DID | Must match registration DID | ❌ TODO |
| **controller** | Return as stored | Validate controller format | ❌ TODO |
| **verificationMethod** | Return with full details | Validate VM structure | ❌ TODO |
| **authentication** | Return reference/embed | Validate references exist | ❌ TODO |
| **assertionMethod** | Return reference/embed | Validate references exist | ❌ TODO |
| **service** | Return service endpoints | Validate service format | ❌ TODO |

### AccumulateKeyPage Format

| Property | Resolver Output | Registrar Input | Consistency Status |
|----------|----------------|-----------------|-------------------|
| **type** | "AccumulateKeyPage" | Must be "AccumulateKeyPage" | ❌ TODO |
| **keyPageUrl** | Return as "acc://..." | Validate "acc://..." format | ❌ TODO |
| **threshold** | Return as number | Validate positive integer | ❌ TODO |
| **controller** | Return DID reference | Validate matches DID | ❌ TODO |

### Metadata Structure

| Field | Resolver Output | Registrar Generation | Consistency Status |
|-------|----------------|---------------------|-------------------|
| **versionId** | From stored metadata | Generate on create/update | ❌ TODO |
| **created** | First version timestamp | Set on create | ❌ TODO |
| **updated** | Latest version timestamp | Set on update | ❌ TODO |
| **deactivated** | Boolean from document | Set on deactivate | ❌ TODO |

## Canonical JSON Consistency

### Encoding Rules

| Rule | Resolver Implementation | Registrar Implementation | Status |
|------|------------------------|-------------------------|--------|
| **Key Ordering** | Lexicographic sort | Lexicographic sort | ❌ TODO |
| **Whitespace** | No extra whitespace | No extra whitespace | ❌ TODO |
| **Number Format** | No trailing zeros | No trailing zeros | ❌ TODO |
| **String Escaping** | Minimal escaping | Minimal escaping | ❌ TODO |
| **Duplicate Keys** | Reject on parse | Reject on validation | ❌ TODO |

### Hash Computation

| Aspect | Resolver | Registrar | Consistency Status |
|--------|----------|-----------|-------------------|
| **Algorithm** | SHA-256 | SHA-256 | ❌ TODO |
| **Input Format** | Canonical JSON | Canonical JSON | ❌ TODO |
| **Output Format** | "sha256:..." | "sha256:..." | ❌ TODO |
| **Verification** | Compare stored hash | Generate content hash | ❌ TODO |

## URL Normalization Consistency

### DID URL Handling

| Case | Resolver Normalization | Registrar Normalization | Status |
|------|----------------------|------------------------|--------|
| **Case Sensitivity** | Convert to lowercase | Convert to lowercase | ❌ TODO |
| **Trailing Dots** | Remove trailing dots | Remove trailing dots | ❌ TODO |
| **Query Parameters** | Preserve order | N/A (not used) | ❌ TODO |
| **Fragments** | Preserve as-is | Validate if present | ❌ TODO |
| **Path Components** | Support dereferencing | Validate but don't use | ❌ TODO |

### ADI Name Validation

| Validation Rule | Resolver | Registrar | Status |
|----------------|----------|-----------|--------|
| **Character Set** | [a-zA-Z0-9.-_] | [a-zA-Z0-9.-_] | ❌ TODO |
| **Dot Placement** | No leading/trailing | No leading/trailing | ❌ TODO |
| **Length Limits** | Accumulate limits | Accumulate limits | ❌ TODO |
| **Reserved Names** | Check reserved list | Check reserved list | ❌ TODO |

## Error Handling Consistency

### Error Codes

| Error Scenario | Resolver Error | Registrar Error | Status |
|----------------|----------------|-----------------|--------|
| **DID Not Found** | `notFound` (404) | `notFound` (404) | ❌ TODO |
| **Invalid DID Syntax** | `invalidDid` (400) | `invalidDid` (400) | ❌ TODO |
| **Deactivated DID** | `deactivated` (410) | `conflict` (409) | ❌ TODO |
| **Unauthorized** | N/A | `unauthorized` (403) | ❌ TODO |
| **Invalid Document** | N/A | `invalidDocument` (400) | ❌ TODO |

### Error Response Format

| Field | Resolver | Registrar | Consistency Status |
|-------|----------|-----------|-------------------|
| **error** | Error code string | Error code string | ❌ TODO |
| **message** | Human readable | Human readable | ❌ TODO |
| **details** | Additional context | Additional context | ❌ TODO |
| **requestId** | Request identifier | Request identifier | ❌ TODO |
| **timestamp** | ISO 8601 | ISO 8601 | ❌ TODO |

## Validation Rules Consistency

### DID Document Validation

| Validation | Resolver (on return) | Registrar (on input) | Status |
|------------|---------------------|---------------------|--------|
| **Required Fields** | Validate structure | Validate required | ❌ TODO |
| **Field Types** | Type checking | Type checking | ❌ TODO |
| **Value Constraints** | Range/format checks | Range/format checks | ❌ TODO |
| **Cross-field Validation** | Referential integrity | Referential integrity | ❌ TODO |

### Verification Method Validation

| Check | Resolver | Registrar | Status |
|-------|----------|-----------|--------|
| **ID Format** | Must be valid URI | Must be valid URI | ❌ TODO |
| **Type Support** | Support AccumulateKeyPage | Support AccumulateKeyPage | ❌ TODO |
| **Controller Match** | Must match DID | Must match DID | ❌ TODO |
| **Required Properties** | Complete structure | Complete structure | ❌ TODO |

## Authorization Consistency

### Policy v1 Implementation

| Aspect | Resolver | Registrar | Status |
|--------|----------|-----------|--------|
| **Key Page URL** | Validate format | Enforce auth policy | ❌ TODO |
| **Expected Location** | `acc://<adi>/book/1` | `acc://<adi>/book/1` | ❌ TODO |
| **Threshold Check** | N/A (read-only) | Verify signature threshold | ❌ TODO |
| **Authorization** | N/A | Validate against policy | ❌ TODO |

### Envelope Structure

| Field | Resolver Understanding | Registrar Generation | Status |
|-------|----------------------|---------------------|--------|
| **contentType** | Parse if present | Set to "application/did+ld+json" | ❌ TODO |
| **document** | Extract DID document | Wrap DID document | ❌ TODO |
| **meta.versionId** | Use for metadata | Generate unique ID | ❌ TODO |
| **meta.timestamp** | Use for updated field | Set current time | ❌ TODO |
| **meta.authorKeyPage** | Validate authority | Set from auth policy | ❌ TODO |
| **meta.proof** | Verify integrity | Generate proof data | ❌ TODO |

## Version Management Consistency

### Version ID Generation

| Component | Resolver | Registrar | Status |
|-----------|----------|-----------|--------|
| **Format** | Parse timestamp-hash | Generate timestamp-hash | ❌ TODO |
| **Timestamp** | Extract Unix timestamp | Use current Unix timestamp | ❌ TODO |
| **Hash Prefix** | Extract first 8 chars | Use first 8 chars of hash | ❌ TODO |
| **Uniqueness** | Assume unique | Ensure uniqueness | ❌ TODO |

### Version History

| Aspect | Resolver | Registrar | Status |
|--------|----------|-----------|--------|
| **Storage Model** | Read append-only | Write append-only | ❌ TODO |
| **Version Time** | Support ?versionTime query | N/A | ❌ TODO |
| **Latest Version** | Default to latest | Create new version | ❌ TODO |
| **Previous Version** | Link to previous | Set previousVersionId | ❌ TODO |

## Content Type Handling

### Supported Formats

| Format | Resolver | Registrar | Status |
|--------|----------|-----------|--------|
| **application/did+ld+json** | Default output | Default input | ❌ TODO |
| **application/ld+json** | Alternative output | Alternative input | ❌ TODO |
| **application/json** | Fallback output | Fallback input | ❌ TODO |

### Content Negotiation

| Header | Resolver Behavior | Registrar Behavior | Status |
|--------|------------------|-------------------|--------|
| **Accept** | Respect client preference | N/A | ❌ TODO |
| **Content-Type** | Set appropriate type | Validate input type | ❌ TODO |

## Test Vector Alignment

### Shared Test Data

| Test Category | Resolver Tests | Registrar Tests | Status |
|---------------|----------------|-----------------|--------|
| **Valid Documents** | Use for resolution | Use for creation | ❌ TODO |
| **Invalid Documents** | Return errors | Reject creation | ❌ TODO |
| **Edge Cases** | Handle gracefully | Validate properly | ❌ TODO |
| **Canonical JSON** | Parse correctly | Generate correctly | ❌ TODO |

### Round-trip Testing

| Test | Description | Status |
|------|-------------|--------|
| **Create → Resolve** | Register then resolve same document | ❌ TODO |
| **Update → Resolve** | Update then resolve latest version | ❌ TODO |
| **Deactivate → Resolve** | Deactivate then resolve shows deactivated | ❌ TODO |
| **Version History** | Create multiple versions, resolve each | ❌ TODO |

## Configuration Consistency

### Accumulate Client

| Setting | Resolver | Registrar | Status |
|---------|----------|-----------|--------|
| **API URL** | Same endpoint | Same endpoint | ❌ TODO |
| **Timeout** | Read timeout | Write timeout | ❌ TODO |
| **Retry Policy** | Read retries | Write retries | ❌ TODO |
| **Authentication** | API credentials | API credentials | ❌ TODO |

### Validation Rules

| Setting | Resolver | Registrar | Status |
|---------|----------|-----------|--------|
| **Max Document Size** | Same limit | Same limit | ❌ TODO |
| **Max Array Length** | Same limit | Same limit | ❌ TODO |
| **Allowed VM Types** | Same types | Same types | ❌ TODO |
| **Service Endpoint Limits** | Same limits | Same limits | ❌ TODO |

## Monitoring and Metrics

### Shared Metrics

| Metric | Resolver | Registrar | Status |
|--------|----------|-----------|--------|
| **Request Count** | Track resolutions | Track operations | ❌ TODO |
| **Error Rate** | Track resolution errors | Track operation errors | ❌ TODO |
| **Latency** | Resolution time | Operation time | ❌ TODO |
| **Accumulate Calls** | API call count | API call count | ❌ TODO |

### Health Checks

| Check | Resolver | Registrar | Status |
|-------|----------|-----------|--------|
| **Service Health** | Return 200 OK | Return 200 OK | ✅ DONE |
| **Accumulate Connectivity** | Test API connection | Test API connection | ❌ TODO |
| **Database Connectivity** | N/A (stateless) | N/A (stateless) | ✅ N/A |

## Integration Testing

### Cross-Service Tests

| Test Scenario | Description | Status |
|---------------|-------------|--------|
| **Create → Resolve** | Registrar creates, resolver resolves | ❌ TODO |
| **Update → Resolve** | Registrar updates, resolver gets latest | ❌ TODO |
| **Deactivate → Resolve** | Registrar deactivates, resolver shows status | ❌ TODO |
| **Error Consistency** | Both services return same errors for same inputs | ❌ TODO |

### Data Consistency Tests

| Test | Description | Status |
|------|-------------|--------|
| **Canonical Equivalence** | Same document canonicalizes identically | ❌ TODO |
| **Hash Verification** | Registrar hash matches resolver verification | ❌ TODO |
| **Metadata Alignment** | Generated metadata matches resolved metadata | ❌ TODO |

## Documentation Alignment

### API Documentation

| Section | Resolver Docs | Registrar Docs | Status |
|---------|---------------|----------------|--------|
| **Error Codes** | List all errors | List all errors | ❌ TODO |
| **Request Format** | N/A | Document structure | ❌ TODO |
| **Response Format** | Document structure | Document structure | ❌ TODO |
| **Examples** | Valid requests/responses | Valid requests/responses | ❌ TODO |

### Code Examples

| Example Type | Resolver | Registrar | Status |
|--------------|----------|-----------|--------|
| **Basic Usage** | Resolution example | Creation example | ❌ TODO |
| **Error Handling** | Error response example | Error response example | ❌ TODO |
| **Advanced Features** | Version time example | Update example | ❌ TODO |

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
- Implemented: 2 (2%)
- In Progress: 0 (0%)
- Remaining: 87 (98%)

*This checklist should be reviewed after each service implementation milestone.*