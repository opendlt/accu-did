# API Versioning Policy

This document defines the versioning strategy and stability guarantees for the Accumulate DID REST APIs.

## Semantic Versioning

All API versions follow [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html):

- **MAJOR**: Breaking changes to REST paths, request/response schemas, or status codes
- **MINOR**: Backward-compatible additions (new endpoints, optional fields, additional status codes)
- **PATCH**: Bug fixes and clarifications that don't affect the API contract

## Version Source of Truth

The canonical version is defined in the OpenAPI specification's `info.version` field:

- **Resolver API**: `docs/spec/openapi/resolver.yaml`
- **Registrar API**: `docs/spec/openapi/registrar.yaml`

Both services MUST maintain identical version numbers for synchronized releases.

## API Stability Rules

### Breaking Changes (MAJOR)

The following constitute breaking changes requiring a major version increment:

1. **Path Changes**:
   - Removing endpoints
   - Changing HTTP methods
   - Modifying path parameters or structure

2. **Request Schema Changes**:
   - Removing required fields
   - Changing field types or constraints
   - Adding required fields without defaults

3. **Response Schema Changes**:
   - Removing fields from successful responses
   - Changing field types or meaning
   - Modifying error response structure

4. **Status Code Changes**:
   - Changing which status codes are returned for specific conditions
   - Removing previously documented status codes

### Non-Breaking Changes (MINOR/PATCH)

The following are considered backward-compatible:

1. **Additions (MINOR)**:
   - New endpoints
   - Optional request fields
   - Additional response fields
   - New status codes for new error conditions
   - Enhanced error details

2. **Clarifications (PATCH)**:
   - Documentation improvements
   - Example updates
   - Schema description enhancements

## Deprecation Policy

### Marking Deprecated Features

Use the OpenAPI `x-deprecated` extension to mark deprecated features:

```yaml
paths:
  /legacy-endpoint:
    get:
      x-deprecated: true
      x-deprecation-notice: "Use /new-endpoint instead. Removal planned in v2.0.0"
      summary: Legacy endpoint (deprecated)
```

### Deprecation Timeline

1. **MINOR Release**: Feature marked as deprecated, continues to function
2. **Next MINOR**: Deprecation warnings in logs/responses (if applicable)
3. **Next MAJOR**: Deprecated feature removed

**Minimum Support Period**: Deprecated features MUST be supported for at least one MINOR release cycle.

## Error Response Contract

All error responses MUST follow this canonical envelope:

```json
{
  "code": 400,
  "error": "invalidRequest",
  "message": "Human-readable error description",
  "details": {
    "field": "additionalErrorContext"
  }
}
```

### Required Fields

- **`code`**: HTTP status code (integer)
- **`error`**: Machine-readable error identifier (string)
- **`message`**: Human-readable description (string)

### Optional Fields

- **`details`**: Additional error context (object, only when helpful)

### Error Identifier Stability

Error identifiers in the `error` field are part of the API contract:

- **MAJOR**: Can remove or change meaning of error identifiers
- **MINOR**: Can add new error identifiers
- **PATCH**: Can improve error messages without changing identifiers

## API Freeze Process

### Freeze Markers

APIs ready for production stability are marked with:

```yaml
openapi: 3.1.0
info:
  version: "0.9.0"
x-api-freeze: true
```

The `x-api-freeze: true` extension indicates:

- API shape is frozen for v1.0 readiness
- Only bug fixes (PATCH) until v1.0 release
- No new features or breaking changes

### Pre-v1.0 Disclaimer

For versions < 1.0.0, breaking changes MAY occur in MINOR releases, but:

- Breaking changes MUST be documented in CHANGELOG
- Migration guides SHOULD be provided for significant changes
- Advance notice SHOULD be given when possible

### Post-v1.0 Guarantees

Starting with v1.0.0:

- MAJOR version changes are the only way to introduce breaking changes
- APIs maintain backward compatibility within MAJOR versions
- Deprecation policy is strictly enforced

## Content Type Versioning

### Standard Media Types

- **DID Documents**: `application/did+json` (W3C standard)
- **API Responses**: `application/json`
- **Health Checks**: `application/json`

### No Custom Versioned Media Types

We do NOT use custom media types for versioning (e.g., `application/vnd.accu-did.v1+json`).

Version information is conveyed through:

1. OpenAPI specification version
2. Service deployment version
3. Optional response headers (if needed)

## Implementation Guidelines

### Service Implementation

1. **Version Validation**: Services SHOULD validate requests against their OpenAPI schema
2. **Error Consistency**: All error responses MUST use the canonical error envelope
3. **Deprecation Warnings**: Services MAY include deprecation warnings in response headers

### Client Implementation

1. **Schema Validation**: Clients SHOULD validate responses against expected schemas
2. **Forward Compatibility**: Clients SHOULD ignore unknown response fields
3. **Error Handling**: Clients MUST handle the canonical error response format

### Documentation

1. **Change Documentation**: All API changes MUST be documented in `docs/api/CHANGELOG.md`
2. **Migration Guides**: Breaking changes SHOULD include migration examples
3. **Deprecation Notices**: Deprecated features MUST be clearly marked in documentation

## Validation and Tooling

### Automated Verification

The repository includes automated OpenAPI validation:

- **`scripts/api-verify.ps1`** - PowerShell validation script
- **`scripts/api-verify.sh`** - Bash validation script

These scripts use `swagger-cli validate` to validate:

- OpenAPI schema compliance
- Consistent error response schemas
- Required field presence

### Continuous Integration

API verification SHOULD run in CI/CD pipelines to ensure:

- OpenAPI specifications remain valid
- Schema changes follow versioning rules
- Error responses maintain canonical format

Run verification locally:

```bash
# PowerShell (Windows)
powershell -ExecutionPolicy Bypass -File .\scripts\api-verify.ps1

# Bash (Unix/macOS)
./scripts/api-verify.sh
```

## References

- [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html)
- [OpenAPI 3.1.0 Specification](https://spec.openapis.org/oas/v3.1.0)
- [W3C DID Core](https://www.w3.org/TR/did-core/)
- [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)