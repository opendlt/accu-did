# API Changelog

This document tracks API-specific changes to the Accumulate DID REST services. For general project changes, see the main [CHANGELOG.md](../../CHANGELOG.md).

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) for REST API contracts.

## [Unreleased]

## [0.9.0] - 2024-09-21

### API Freeze 🔒

- **BREAKING**: API shape frozen for v1.0 readiness
- Added `x-api-freeze: true` to OpenAPI specifications
- Established canonical error response format: `{ code, error, message, details? }`
- Created comprehensive API versioning policy

### Added

#### Resolver API (`/resolve`, `/1.0/identifiers/{did}`)
- **W3C DID Core Compliance**: Complete DID resolution following W3C specifications
- **Universal Resolver 1.0**: Compatible with DIF Universal Resolver patterns
- **Historical Resolution**: Optional `versionTime` parameter for point-in-time resolution
- **Multiple Operation Modes**: FAKE (offline testing) and REAL (blockchain) modes
- **Rich Error Handling**: Structured error responses with appropriate HTTP status codes

**Endpoints**:
- `GET /healthz` - Service health monitoring
- `GET /resolve?did={did}` - Native DID resolution with full metadata
- `GET /1.0/identifiers/{did}` - Universal Resolver compatible endpoint

**Status Codes**:
- `200` - Successful resolution with DID document and metadata
- `404` - DID not found on blockchain
- `410` - DID has been deactivated
- `422` - Invalid DID syntax or format
- `502` - Accumulate node unavailable

#### Registrar API (`/register`, `/1.0/create`, `/1.0/update`, `/1.0/deactivate`)
- **Native Endpoints**: Clean, direct API for Accumulate DID operations
- **Universal Registrar 1.0**: Full compatibility with DIF Universal Registrar specification
- **Complete Lifecycle**: Create, update, and deactivate operations
- **Patch-Based Updates**: Structured service addition/removal via patch operations
- **Transaction Tracking**: Full Accumulate transaction ID tracking for all operations

**Native Endpoints**:
- `POST /register` - Create new DID with complete Accumulate transaction flow
- `POST /update` - Update existing DID with patch operations
- `POST /deactivate` - Deactivate DID with tombstone entry

**Universal Registrar Endpoints**:
- `POST /1.0/create?method=acc` - Universal Registrar create compatibility
- `POST /1.0/update?method=acc` - Universal Registrar update with patch and JSON Patch support
- `POST /1.0/deactivate?method=acc` - Universal Registrar deactivate compatibility

**Status Codes**:
- `200` - Successful operation with transaction details
- `400` - Invalid request data or malformed DID document
- `409` - DID already exists (create only)
- `404` - DID not found (update/deactivate only)
- `502` - Accumulate node unavailable

### Changed

#### Error Response Format Standardization
- **Canonical Envelope**: All errors now use `{ code, error, message, details? }` format
- **Machine-Readable Codes**: Error `code` field provides HTTP status code
- **Structured Identifiers**: Error `error` field provides machine-readable error type
- **Human-Readable Messages**: Error `message` field provides descriptive text
- **Optional Context**: Error `details` field provides additional context when helpful

**Error Identifier Mapping**:
- `invalidDid` - Malformed DID syntax or unsupported format
- `notFound` - DID does not exist on the blockchain
- `deactivated` - DID has been deactivated (resolver only)
- `invalidRequest` - Malformed request body or missing required fields
- `alreadyExists` - DID already exists (registrar create only)
- `unauthorized` - Authentication or authorization failure (when applicable)
- `methodNotSupported` - Accumulate node unavailable or unsupported operation

#### Schema Enhancements
- **W3C DID Document Schema**: Complete compliance with W3C DID Core specification
- **Verification Method Support**: Ed25519VerificationKey2020 with multibase encoding
- **Service Endpoint Management**: Full CRUD operations for service endpoints
- **Flexible Authentication**: Support for authentication and assertionMethod arrays
- **Metadata Enrichment**: Comprehensive DID document and resolution metadata

### Security

#### Input Validation and Sanitization
- **DID Syntax Validation**: Strict pattern matching for `did:acc:` format
- **Schema Validation**: Request/response validation against OpenAPI schemas
- **Multibase Key Validation**: Proper validation of cryptographic key formats
- **URL Sanitization**: Safe handling of service endpoint URLs

#### Key Management
- **Secure Key Handling**: Proper cryptographic key management for DID operations
- **No Secret Exposure**: Secrets and private keys never logged or exposed in responses
- **Transaction Security**: Secure submission to Accumulate blockchain

### Documentation

#### OpenAPI 3.1.0 Specifications
- **Complete API Coverage**: Full specification of all endpoints and schemas
- **Rich Examples**: Comprehensive request/response examples for all operations
- **Interactive Documentation**: Redoc-compatible specifications for easy exploration
- **Type Safety**: Detailed schema definitions with proper type constraints

#### API Documentation
- **Versioning Policy**: Comprehensive documentation of API versioning strategy
- **Migration Guides**: Clear guidance for API version transitions
- **Error Handling**: Complete documentation of error scenarios and responses

## [0.1.0] - 2024-09-21

### Added
- Initial API implementation
- Basic DID resolution and registration
- OpenAPI specifications foundation
- Health check endpoints

---

**Version History**:
- [Unreleased]: https://github.com/opendlt/accu-did/compare/v0.9.0...HEAD
- [0.9.0]: https://github.com/opendlt/accu-did/compare/v0.1.0...v0.9.0
- [0.1.0]: https://github.com/opendlt/accu-did/releases/tag/v0.1.0