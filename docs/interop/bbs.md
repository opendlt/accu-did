# BBS+ Selective Disclosure

## Overview
Planned integration of BBS+ signatures for advanced selective disclosure capabilities with `did:acc` method.

## Selective Disclosure Path

### Phase 1: Foundation (Q2 2024)
- [ ] BBS+ signature suite evaluation
- [ ] Verification method type definition
- [ ] Basic proof generation/verification

### Phase 2: Integration (Q3 2024)
- [ ] Accumulate key page mapping to BBS+ keys
- [ ] Multi-message signature support
- [ ] Selective disclosure protocols

### Phase 3: Advanced Features (Q4 2024)
- [ ] Zero-knowledge proof integration
- [ ] Predicate proofs (age verification, etc.)
- [ ] Anonymous credentials support

## Verification Method Example

```json
{
  "verificationMethod": [
    {
      "id": "did:acc:alice#bbs-key-1",
      "type": "Bls12381G2Key2020",
      "controller": "did:acc:alice",
      "publicKeyBase58": "nEP2DEdbRaQ2r5Azeatui9MG6cj7JUHa8GD7khub4egHJREEuvj4Y8YG8w671LnZ"
    }
  ],
  "assertionMethod": ["#bbs-key-1"]
}
```

## Use Cases Targeted

### Anonymous Credentials
- **Age Verification**: Prove over 18 without revealing exact age
- **Location Proof**: Prove residency in region without exact address
- **Qualification Proof**: Prove skill level without revealing institution

### Privacy-Preserving Identity
- **Selective Attribute Disclosure**: Share only required identity attributes
- **Unlinkable Presentations**: Multiple verifications can't be correlated
- **Revocation Privacy**: Check credential status without revealing identity

## Technical Considerations

### Key Management
- **BLS12-381 Key Generation**: Derive from Accumulate key pages
- **Key Rotation**: Handle BBS+ key updates via DID document updates
- **Threshold Signatures**: Explore multi-party BBS+ signing

### Performance
- **Proof Generation Time**: Optimize for mobile devices
- **Verification Speed**: Ensure practical verification times
- **Proof Size**: Minimize proof overhead for network transmission

## Library Candidates

### Go Implementations
- [ ] **[bbs-signatures-go](https://github.com/mattrglobal/bbs-signatures-go)**
  - Mattr Global implementation
  - Active development
  - Commercial support available

- [ ] **Custom Implementation**
  - Build on existing crypto libraries
  - Integrate with Accumulate cryptography
  - Full control over features

### Evaluation Criteria
- Standards compliance (W3C BBS+ Signature Suite)
- Performance benchmarks
- Security audit status
- Integration complexity with Accumulate

## Implementation TODOs

- [ ] **Research Phase**
  - [ ] BBS+ signature suite specification review
  - [ ] Library evaluation and benchmarking
  - [ ] Security model analysis

- [ ] **Proof of Concept**
  - [ ] Basic BBS+ signing with Accumulate keys
  - [ ] Simple selective disclosure demo
  - [ ] Performance baseline establishment

- [ ] **Production Planning**
  - [ ] Key derivation standards
  - [ ] Verification method registration
  - [ ] Revocation mechanism design

## Specification References

- See [spec/did-acc-method.md](../spec/did-acc-method.md) Section 3.1 (Verification Methods)
- See [spec/did-acc-method.md](../spec/did-acc-method.md) Section 6.3 (Advanced Features)

## Standards References

- [BBS+ Signature Scheme](https://eprint.iacr.org/2016/663.pdf)
- [W3C BBS+ Signature Suite (Draft)](https://w3c-ccg.github.io/ldp-bbs2020/)
- [JSON-LD BBS+ Signatures](https://w3c-ccg.github.io/lds-bbs2020/)

## Timeline

| Quarter | Milestone | Deliverable |
|---------|-----------|------------|
| Q2 2024 | Research Complete | Library selection, architecture design |
| Q3 2024 | MVP Implementation | Basic BBS+ signing and verification |
| Q4 2024 | Production Ready | Full selective disclosure support |
| Q1 2025 | Advanced Features | Anonymous credentials, predicate proofs |