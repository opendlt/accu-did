# Phase 1 Implementation Backlog

## Overview
Prioritized task list for Phase 1 implementation of the Accumulate DID stack. Tasks are organized by priority and dependencies.

## Task Status Legend
- 🔴 **Critical** - Blocks other work
- 🟡 **High** - Core functionality
- 🟢 **Medium** - Important but not blocking
- 🔵 **Low** - Nice to have

## Sprint 1: Foundation (Days 1-3)

### 🔴 Critical Path
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| F-001 | Create did-acc-method.md specification | spec | TODO | Mapper | Defines DID syntax, operations |
| F-002 | Define Rules.md encoding standards | spec | TODO | SDK Mentor | Canonical JSON, hashing |
| F-003 | Update CLAUDE.md with commands | docs | TODO | Docsmith | Build, test, run instructions |
| F-004 | Create example DID documents | spec | TODO | Mapper | entry.v1, update, deactivate |
| F-005 | Create test vectors JSON files | spec | TODO | SDK Mentor | URL norm, envelope, auth |

### 🟡 High Priority
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| F-006 | Create PARITY-SPEC-RESOLVER.md | spec | TODO | Mapper | Spec compliance checklist |
| F-007 | Create PARITY-RESOLVER-REGISTRAR.md | spec | TODO | Mapper | Service consistency matrix |
| F-008 | Create PARITY-UNI-DRIVERS.md | spec | TODO | Mapper | Universal driver compat |
| F-009 | Set up Go module dependencies | all | TODO | Scaffolder | gin, testify, viper |

## Sprint 2: Resolver Core (Days 4-7)

### 🔴 Critical Path
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| R-001 | Create resolver handler structure | resolver-go | TODO | Resolver Builder | handlers/resolve.go |
| R-002 | Implement URL normalization | resolver-go | TODO | Resolver Builder | Case-insensitive, fragments |
| R-003 | Implement canonical JSON | resolver-go | TODO | Resolver Builder | RFC8785 or stable alt |
| R-004 | Create Accumulate client stub | resolver-go | TODO | Resolver Builder | Mock for offline testing |
| R-005 | Implement DID resolution logic | resolver-go | TODO | Resolver Builder | Core resolve function |

### 🟡 High Priority
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| R-006 | Add versionTime support | resolver-go | TODO | Resolver Builder | Query parameter handling |
| R-007 | Generate metadata fields | resolver-go | TODO | Resolver Builder | updated, versionId, deactivated |
| R-008 | Create table-driven tests | resolver-go | TODO | Resolver Builder | Using golden files |
| R-009 | Add SHA-256 content hash | resolver-go | TODO | Resolver Builder | For verification |
| R-010 | Create Makefile | resolver-go | TODO | Scaffolder | build, test, run, lint |

### 🟢 Medium Priority
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| R-011 | Add .golangci.yml | resolver-go | TODO | Scaffolder | Linting configuration |
| R-012 | Create README.md | resolver-go | TODO | Docsmith | API docs, examples |
| R-013 | Copy goldens to testdata | resolver-go | TODO | Scaffolder | From spec/examples |
| R-014 | Add error handling | resolver-go | TODO | Resolver Builder | Proper HTTP errors |
| R-015 | Add logging | resolver-go | TODO | Resolver Builder | Structured logging |

## Sprint 3: Registrar Core (Days 8-11)

### 🔴 Critical Path
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| G-001 | Create registrar handlers | registrar-go | TODO | Registrar Builder | create, update, deactivate |
| G-002 | Implement envelope structure | registrar-go | TODO | Registrar Builder | ops/envelope.go |
| G-003 | Implement auth policy v1 | registrar-go | TODO | Registrar Builder | acc://<adi>/book/1 |
| G-004 | Create Accumulate client stub | registrar-go | TODO | Registrar Builder | Echo tx stub |
| G-005 | Implement registration logic | registrar-go | TODO | Registrar Builder | Core operations |

### 🟡 High Priority
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| G-006 | Add DID validation | registrar-go | TODO | Registrar Builder | Document structure |
| G-007 | Generate versionId | registrar-go | TODO | Registrar Builder | Unique identifiers |
| G-008 | Add content hash tracking | registrar-go | TODO | Registrar Builder | For verification |
| G-009 | Create integration tests | registrar-go | TODO | Registrar Builder | End-to-end flows |
| G-010 | Create Makefile | registrar-go | TODO | Scaffolder | build, test, run, lint |

### 🟢 Medium Priority
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| G-011 | Add .golangci.yml | registrar-go | TODO | Scaffolder | Linting configuration |
| G-012 | Create README.md | registrar-go | TODO | Docsmith | API docs, examples |
| G-013 | Copy goldens to testdata | registrar-go | TODO | Scaffolder | From spec/vectors |
| G-014 | Add request validation | registrar-go | TODO | Registrar Builder | Input sanitization |
| G-015 | Add audit logging | registrar-go | TODO | Registrar Builder | Operation tracking |

## Sprint 4: Universal Drivers (Days 12-14)

### 🟡 High Priority
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| U-001 | Create uni-resolver driver | drivers | TODO | Uni Engineer | GET /1.0/identifiers/{did} |
| U-002 | Create uni-registrar driver | drivers | TODO | Uni Engineer | POST /1.0/{operations} |
| U-003 | Write Dockerfiles | drivers | TODO | Uni Engineer | Multi-stage builds |
| U-004 | Create docker-compose.yml | drivers | TODO | Uni Engineer | Service orchestration |
| U-005 | Write smoke tests | drivers | TODO | Uni Engineer | smoke.ps1, smoke.sh |

### 🟢 Medium Priority
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| U-006 | Add driver READMEs | drivers | TODO | Docsmith | Setup instructions |
| U-007 | Create healthchecks | drivers | TODO | Uni Engineer | Docker health probes |
| U-008 | Add driver configuration | drivers | TODO | Uni Engineer | Environment variables |
| U-009 | Test Universal format | drivers | TODO | Uni Engineer | Compliance validation |

## Sprint 5: SDK & Documentation (Days 15-17)

### 🟡 High Priority
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| S-001 | Create SDK client interface | sdks/go | TODO | SDK Mentor | High-level API |
| S-002 | Implement resolver helpers | sdks/go | TODO | SDK Mentor | Resolution utilities |
| S-003 | Implement registrar helpers | sdks/go | TODO | SDK Mentor | Registration utilities |
| S-004 | Define common types | sdks/go | TODO | SDK Mentor | Shared structures |
| S-005 | Create usage examples | sdks/go | TODO | SDK Mentor | Sample applications |

### 🟢 Medium Priority
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| D-001 | Write index.md | docs | TODO | Docsmith | Getting started guide |
| D-002 | Write resolver.md | docs | TODO | Docsmith | Resolver API reference |
| D-003 | Write registrar.md | docs | TODO | Docsmith | Registrar API reference |
| D-004 | Write quickstart-go.md | docs | TODO | Docsmith | Go SDK tutorial |
| D-005 | Create mkdocs.yml | docs | TODO | Docsmith | Documentation config |

### 🔵 Low Priority
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| D-006 | Write didcomm.md stub | docs/interop | TODO | Interop Engineer | Future integration |
| D-007 | Write sd-jwt.md stub | docs/interop | TODO | Interop Engineer | Future support |
| D-008 | Write bbs.md stub | docs/interop | TODO | Interop Engineer | Future signatures |

## Sprint 6: CI/CD (Days 18-20)

### 🟡 High Priority
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| C-001 | Create resolver.yml workflow | .github | TODO | CI Keeper | Build and test |
| C-002 | Create registrar.yml workflow | .github | TODO | CI Keeper | Build and test |
| C-003 | Create drivers.yml workflow | .github | TODO | CI Keeper | Docker builds |
| C-004 | Create docs.yml workflow | .github | TODO | CI Keeper | MkDocs build |
| C-005 | Add golangci-lint action | .github | TODO | CI Keeper | Code quality |

### 🟢 Medium Priority
| ID | Task | Component | Status | Assignee | Notes |
|----|------|-----------|--------|----------|-------|
| C-006 | Set up test coverage | .github | TODO | CI Keeper | Codecov integration |
| C-007 | Add release automation | .github | TODO | CI Keeper | Tag-based releases |
| C-008 | Configure dependabot | .github | TODO | CI Keeper | Dependency updates |

## Dependencies Graph

```
Foundation Tasks (F-*)
    ↓
Resolver Core (R-001 to R-005)
    ↓
Resolver Tests (R-006 to R-015)
    ↓
Registrar Core (G-001 to G-005)
    ↓
Registrar Tests (G-006 to G-015)
    ↓
Universal Drivers (U-*)
    ↓
SDK Development (S-*)
    ↓
Documentation (D-*)
    ↓
CI/CD Setup (C-*)
```

## Risk Items

### High Risk
- Accumulate API integration complexity
- Canonical JSON implementation variations
- Key rotation and auth verification

### Mitigation
- Use stub clients for offline development
- Implement multiple canonical JSON options
- Clear separation of concerns in auth flow

## Definition of Done

### For Code Tasks
- [ ] Code implemented and compiles
- [ ] Unit tests written and passing
- [ ] Integration tests passing (if applicable)
- [ ] Code reviewed (self-review minimum)
- [ ] Documentation updated
- [ ] Linting passes

### For Documentation Tasks
- [ ] Content complete and accurate
- [ ] Examples provided
- [ ] Reviewed for clarity
- [ ] Links verified
- [ ] Formatting correct

### For Infrastructure Tasks
- [ ] Scripts tested on target platforms
- [ ] Configuration documented
- [ ] Smoke tests passing
- [ ] README updated

## Progress Tracking

### Sprint 1: Foundation
- Total Tasks: 9
- Completed: 0
- In Progress: 0
- Blocked: 0

### Sprint 2: Resolver Core
- Total Tasks: 15
- Completed: 0
- In Progress: 0
- Blocked: 0

### Sprint 3: Registrar Core
- Total Tasks: 15
- Completed: 0
- In Progress: 0
- Blocked: 0

### Sprint 4: Universal Drivers
- Total Tasks: 9
- Completed: 0
- In Progress: 0
- Blocked: 0

### Sprint 5: SDK & Documentation
- Total Tasks: 13
- Completed: 0
- In Progress: 0
- Blocked: 0

### Sprint 6: CI/CD
- Total Tasks: 8
- Completed: 0
- In Progress: 0
- Blocked: 0

## Notes

- All tasks should be completed with offline/mock testing first
- Real Accumulate integration comes after Phase 1
- Focus on clean interfaces and testability
- Maintain backwards compatibility in all changes

---

*Last Updated: Sprint 1 - Foundation*
*Total Tasks: 69*
*Estimated Completion: 20 days*