# SD-JWT Integration

## Overview
Selective Disclosure JSON Web Tokens (SD-JWT) integration with `did:acc` method for privacy-preserving credential verification.

## Quick Implementation Plan

### Issuance Flow
```
1. Credential Issuer (did:acc:issuer)
   ├── Create claims structure
   ├── Generate disclosure digests
   ├── Sign JWT with Accumulate key page
   └── Return SD-JWT + disclosures

2. Holder (did:acc:holder)
   ├── Receive SD-JWT from issuer
   ├── Store disclosures securely
   └── Selectively reveal claims to verifier
```

### Verification Flow
```
1. Verifier receives presentation
   ├── Parse SD-JWT structure
   ├── Resolve issuer DID (did:acc:issuer)
   ├── Verify JWT signature against Accumulate key page
   └── Validate disclosed claims
```

## Claims Mapping Example

```json
{
  "iss": "did:acc:university.credentials",
  "sub": "did:acc:alice",
  "iat": 1704067200,
  "exp": 1735689600,
  "_sd": [
    "CrQe7S0le5JAHkYWJkYfHZrsP-7mNiNZMAhqJWKwC0U",  // name
    "JzYjH4svliH0R3PyEMfeZu6Jt69u5qehZo_F7ep6SqE",  // birthdate
    "PorFbpKuVu6xymJagvkFsFXAbRoc2JGlAUA2BA4o7cI"   // student_id
  ],
  "degree": {
    "type": "BachelorOfScience",
    "field": "Computer Science"
  },
  "_sd_alg": "sha-256"
}
```

## Library Candidates (Go)

### Primary Options
- [ ] **[go-sd-jwt](https://github.com/MichaelFraser99/go-sd-jwt)**
  - Native Go implementation
  - Active development
  - Standard compliance

- [ ] **[jose](https://github.com/go-jose/go-jose) + Custom SD-JWT**
  - Mature JWT/JWS library
  - Build SD-JWT layer on top
  - Full control over implementation

### Secondary Options
- [ ] **[golang-jwt](https://github.com/golang-jwt/jwt) + Extensions**
  - Popular JWT library
  - Need custom SD extensions
  - Good performance

## Implementation TODOs

- [ ] **SD-JWT Library Evaluation**
  - [ ] Feature comparison and compliance testing
  - [ ] Performance benchmarks
  - [ ] Security audit assessment

- [ ] **Accumulate Integration**
  - [ ] Map verification methods to JWT signatures
  - [ ] Handle key rotation in long-lived credentials
  - [ ] Threshold signature support evaluation

- [ ] **Privacy Features**
  - [ ] Disclosure management UX
  - [ ] Zero-knowledge proof integration planning
  - [ ] Unlinkability assessment

- [ ] **Credential Types**
  - [ ] University credentials schema
  - [ ] Employment verification schema
  - [ ] Government ID schema

## Specification References

- See [spec/did-acc-method.md](../spec/did-acc-method.md) Section 3.1 (Verification Methods)
- See [spec/did-acc-method.md](../spec/did-acc-method.md) Section 5.1 (Signature Verification)

## Standards Compliance

- [SD-JWT Specification](https://datatracker.ietf.org/doc/draft-ietf-oauth-selective-disclosure-jwt/)
- [JSON Web Token (JWT)](https://datatracker.ietf.org/doc/html/rfc7519)
- [JSON Web Signature (JWS)](https://datatracker.ietf.org/doc/html/rfc7515)