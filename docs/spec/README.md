# Spec & Docs Index

- OpenAPI (Resolver): `docs/spec/openapi/resolver.yaml`
- OpenAPI (Registrar): `docs/spec/openapi/registrar.yaml`
- DID Method Draft: `docs/spec/method.md`
- Diagrams (Mermaid): `docs/spec/diagrams/*.mmd`

## Building a static site
- **PowerShell (Windows):** `.\scripts\build-docs.ps1`
- **Unix/mac (with make):** `make docs`
- Output goes to `docs/site/`:
  - `resolver.html` (Redoc)
  - `registrar.html` (Redoc)
  - `index.html` (links + Mermaid previews)

**Builder preference:** Uses Docker image `redocly/redoc` when available, else `npx redoc-cli`.

## Quick Preview Tips

### Mermaid Diagrams
- **Online**: Copy `.mmd` content to [https://mermaid.live](https://mermaid.live) for instant preview
- **VS Code**: Install Mermaid Preview extension
- **GitHub**: All `.mmd` files render automatically in the repository

### OpenAPI Specifications
- **Swagger Editor**:
  - Web: [https://editor.swagger.io](https://editor.swagger.io) - paste YAML content
  - Local: `docker run -p 8082:8080 swaggerapi/swagger-editor`
- **Redoc**:
  - Online: [https://redocly.github.io/redoc/](https://redocly.github.io/redoc/)
  - CLI: `npm install -g redoc-cli && redoc-cli serve resolver.yaml`

## Method Specification

- **[method.md](method.md)** - Complete `did:acc` DID Method specification
  - ABNF grammar and syntax rules
  - Operation mappings to Accumulate transactions
  - Security considerations and examples
  - W3C DID Core compliance

## API Specifications

### OpenAPI 3.1.0 Documents
- **[openapi/resolver.yaml](openapi/resolver.yaml)** - DID Resolution service API
  - Health checks and DID resolution endpoints
  - W3C DID Core resolution result format
  - Universal Resolver 1.0 compatibility
  - FAKE vs REAL mode configuration

- **[openapi/registrar.yaml](openapi/registrar.yaml)** - DID Registration service API
  - Native create/update/deactivate endpoints
  - Universal Registrar 1.0 compatible endpoints
  - Patch-based updates and service management
  - Complete request/response schemas

- **[openapi/README.md](openapi/README.md)** - Usage examples and configuration
  - Swagger/Redoc preview instructions
  - Complete curl examples for all endpoints
  - Service configuration and troubleshooting

### API Documentation
- **[../api/CHANGELOG.md](../api/CHANGELOG.md)** - API-specific changes and version history
- **[../api/VERSIONING.md](../api/VERSIONING.md)** - Semantic versioning policy and stability guarantees

## Architecture Diagrams

### System Architecture
- **[diagrams/architecture.mmd](diagrams/architecture.mmd)** - High-level system overview
  - Client → Services → Accumulate Node flow
  - Universal Driver integration points
  - Network topology and port configuration

### Sequence Diagrams
- **[diagrams/create-sequence.mmd](diagrams/create-sequence.mmd)** - DID creation flow
  - CreateIdentity → CreateDataAccount → WriteData
  - Transaction submission and consensus
  - Complete end-to-end timing

- **[diagrams/resolve-sequence.mmd](diagrams/resolve-sequence.mmd)** - DID resolution flow
  - DID parsing and URL mapping
  - Data account queries and response handling
  - Error scenarios (404, 410)

- **[diagrams/update-sequence.mmd](diagrams/update-sequence.mmd)** - DID update flow
  - Current document resolution
  - Patch application and merging
  - Updated document submission

- **[diagrams/deactivate-sequence.mmd](diagrams/deactivate-sequence.mmd)** - DID deactivation flow
  - Tombstone document creation
  - Deactivation transaction submission
  - Resolution behavior after deactivation

## Testing Resources

- **[../../postman/](../../postman/)** - Postman collection and environment
  - Complete API test suite
  - Local development environment
  - Example request/response flows

- **[../../examples/hello_accu/](../../examples/hello_accu/)** - Working code example
  - Complete DID lifecycle demonstration
  - Real Accumulate API integration
  - Smoke testing scripts

## Implementation Status

| Component | Status | Description |
|-----------|--------|-------------|
| **Resolver** | ✅ Complete | DID resolution with FAKE/REAL modes |
| **Registrar** | ✅ Complete | DID registration with native + Universal APIs |
| **Universal Drivers** | 🚧 In Progress | Docker containers for ecosystem integration |
| **Method Spec** | ✅ Draft | W3C-style specification document |
| **OpenAPI Specs** | ✅ Complete | Interactive API documentation |
| **Examples** | ✅ Complete | Working code demonstrations |

## Related Resources

- **[W3C DID Core](https://www.w3.org/TR/did-core/)** - DID specification standard
- **[Universal Resolver](https://github.com/decentralized-identity/universal-resolver)** - DIF resolution infrastructure
- **[Universal Registrar](https://github.com/decentralized-identity/universal-registrar)** - DIF registration infrastructure
- **[Accumulate Protocol](https://accumulatenetwork.io)** - Blockchain platform documentation

## Quick Start

1. **Read the method spec**: [method.md](method.md)
2. **Preview diagrams**: Copy [diagrams/architecture.mmd](diagrams/architecture.mmd) to [mermaid.live](https://mermaid.live)
3. **Explore APIs**: Copy [openapi/resolver.yaml](openapi/resolver.yaml) to [editor.swagger.io](https://editor.swagger.io)
4. **Import Postman collection**: [../../postman/](../../postman/)
5. **Try the hello_accu example**: [../../examples/hello_accu/](../../examples/hello_accu/)