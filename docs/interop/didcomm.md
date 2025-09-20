# DIDComm v2 Integration

## Overview
Integration plan for DIDComm v2 messaging protocol with `did:acc` method.

## Service Entry Example

```json
{
  "id": "did:acc:alice",
  "service": [
    {
      "id": "did:acc:alice#didcomm-messaging",
      "type": "DIDCommMessaging",
      "serviceEndpoint": {
        "uri": "https://alice.accumulate.network/didcomm",
        "accept": ["didcomm/v2"],
        "routingKeys": []
      }
    }
  ]
}
```

## Key Agreement Verification Method

### X25519 Example
```json
{
  "verificationMethod": [
    {
      "id": "did:acc:alice#x25519-key-1",
      "type": "X25519KeyAgreementKey2020",
      "controller": "did:acc:alice",
      "publicKeyMultibase": "z6LSbysY2xFMRpGMhb7tFTLMpeuPRaqaWM1yECx2AtzE3KCc"
    }
  ],
  "keyAgreement": ["#x25519-key-1"]
}
```

## Implementation TODOs

- [ ] **Evaluate DIDComm Libraries**
  - [ ] [didcomm-rust](https://github.com/sicpa-dlab/didcomm-rust) with Go bindings
  - [ ] [didcomm-go](https://github.com/hyperledger/aries-framework-go/tree/main/pkg/didcomm) (Aries)
  - [ ] Native Go implementation assessment

- [ ] **Key Management Integration**
  - [ ] Map Accumulate key pages to DIDComm key agreement
  - [ ] X25519 key derivation from Accumulate keys
  - [ ] Key rotation handling for DIDComm contexts

- [ ] **Service Endpoint Registration**
  - [ ] Accumulate data account for service metadata
  - [ ] Service endpoint versioning and updates
  - [ ] Multi-protocol service endpoint support

## Specification References

- See [spec/did-acc-method.md](../spec/did-acc-method.md) Section 4.2 (Service Entries)
- See [spec/did-acc-method.md](../spec/did-acc-method.md) Section 3.3 (Key Agreement Methods)

## Standards Compliance

- [DIDComm v2 Specification](https://identity.foundation/didcomm-messaging/spec/)
- [DID Core - Service Endpoints](https://www.w3.org/TR/did-core/#services)
- [X25519KeyAgreementKey2020](https://w3c-ccg.github.io/lds-x25519-2020/)