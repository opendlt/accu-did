# Cross-platform docs builder. Prefers Docker; falls back to npx.
SHELL := /bin/sh

DOCS_DIR := docs
SPEC_DIR := $(DOCS_DIR)/spec/openapi
SITE_DIR := $(DOCS_DIR)/site

RESOLVER_YAML := $(SPEC_DIR)/resolver.yaml
REGISTRAR_YAML := $(SPEC_DIR)/registrar.yaml
RESOLVER_HTML := $(SITE_DIR)/resolver.html
REGISTRAR_HTML := $(SITE_DIR)/registrar.html
INDEX_HTML := $(SITE_DIR)/index.html

REDOC_DOCKER := docker run --rm -v "$(PWD)":/work -w /work redocly/redoc
REDOC_NPX := npx --yes redoc-cli@0.15.1

.PHONY: docs clean-docs open-docs test docker-build all clean-all ci

docs: $(RESOLVER_HTML) $(REGISTRAR_HTML) $(INDEX_HTML)
	@echo "✅ Docs built under $(SITE_DIR)"

$(RESOLVER_HTML): $(RESOLVER_YAML)
	@mkdir -p $(SITE_DIR)
	@if command -v docker >/dev/null 2>&1; then \
		$(REDOC_DOCKER) build -o $(RESOLVER_HTML) $(RESOLVER_YAML); \
	else \
		$(REDOC_NPX) build $(RESOLVER_YAML) -o $(RESOLVER_HTML); \
	fi

$(REGISTRAR_HTML): $(REGISTRAR_YAML)
	@mkdir -p $(SITE_DIR)
	@if command -v docker >/dev/null 2>&1; then \
		$(REDOC_DOCKER) build -o $(REGISTRAR_HTML) $(REGISTRAR_YAML); \
	else \
		$(REDOC_NPX) build $(REGISTRAR_YAML) -o $(REGISTRAR_HTML); \
	fi

$(INDEX_HTML):
	@mkdir -p $(SITE_DIR)
	@cp $(DOCS_DIR)/spec/diagrams/*.mmd $(SITE_DIR) 2>/dev/null || true
	@cp $(DOCS_DIR)/spec/method.md $(SITE_DIR) 2>/dev/null || true
	@python3 - <<'PY' || node -e "require('fs').writeFileSync('$(INDEX_HTML)', require('fs').readFileSync('$(DOCS_DIR)/site/index.template.html','utf8'))"
import pathlib, shutil
p = pathlib.Path("$(DOCS_DIR)/site")
tpl = (p / "index.template.html").read_text(encoding="utf-8")
(p / "index.html").write_text(tpl, encoding="utf-8")
PY

clean-docs:
	@rm -rf $(SITE_DIR)
	@mkdir -p $(SITE_DIR)
	@echo "🧹 Cleaned $(SITE_DIR)"

open-docs:
	@echo "Open these in your browser:"
	@echo "  - $(RESOLVER_HTML)"
	@echo "  - $(REGISTRAR_HTML)"
	@echo "  - $(INDEX_HTML)"

# Test all services
.PHONY: test
test:
	@echo "🧪 Running tests for all services..."
	@cd resolver-go && go test ./... -v
	@cd registrar-go && go test ./... -v

# Build Docker images
.PHONY: docker-build
docker-build:
	@echo "🐳 Building Docker images..."
	docker build -t accu-did/resolver:latest -f drivers/resolver/Dockerfile .
	docker build -t accu-did/registrar:latest -f drivers/registrar/Dockerfile .

# Combined target: test + docs + docker
.PHONY: all
all: test docs docker-build
	@echo "✅ All tasks completed"

# Clean all build artifacts
.PHONY: clean-all
clean-all: clean-docs
	@echo "🧹 Cleaning all build artifacts..."
	@cd resolver-go && go clean -cache -testcache
	@cd registrar-go && go clean -cache -testcache
	@docker rmi accu-did/resolver:latest accu-did/registrar:latest 2>/dev/null || true

# Local CI (runs all checks)
.PHONY: ci
ci:
	@if [ -x scripts/local-ci.sh ]; then \
		bash scripts/local-ci.sh; \
	else \
		$(MAKE) test && $(MAKE) docs && $(MAKE) docker-build; \
	fi

# ========================================================================
# Container-based development targets (RECOMMENDED)
# ========================================================================

.PHONY: dev-shell test-all ci-local dev-up dev-down check-imports conformance perf help lint test-race vet qa dist-clean binaries-local images-local sbom-local scan-local docs-archive release-local sdk-test sdk-merge-spec example-sdk todo-scan

dev-shell:
	@echo "🐳 Starting development container shell..."
	@docker compose -f docker-compose.dev.yml run --rm dev 'bash -lc "exec bash -i"'

test-all:
	@echo "🧪 Running all tests in container..."
	@docker compose -f docker-compose.dev.yml run --rm dev 'bash -lc "go work sync 2>/dev/null || true; go test ./... -count=1"'

docs-container:
	@echo "📚 Building documentation in container..."
	@docker compose -f docker-compose.dev.yml run --rm dev 'bash -lc "chmod +x scripts/build-docs.sh 2>/dev/null || true; ./scripts/build-docs.sh npx"'

ci-local:
	@echo "🚀 Running local CI in container..."
	@docker compose -f docker-compose.dev.yml run --rm dev 'bash -lc "chmod +x scripts/local-ci.sh 2>/dev/null || true; ./scripts/local-ci.sh"'

dev-up:
	@echo "🚀 Starting development services..."
	@docker compose -f docker-compose.dev.yml up -d --build resolver registrar

dev-down:
	@echo "🛑 Stopping development services..."
	@docker compose -f docker-compose.dev.yml down

check-imports:
	@echo "🔍 Checking imports in container..."
	@docker compose -f docker-compose.dev.yml run --rm dev 'bash -lc "chmod +x scripts/check-imports.sh 2>/dev/null || true; ./scripts/check-imports.sh"'

conformance:
	@echo "🔍 Running conformance tests in container..."
	@docker compose -f docker-compose.dev.yml run --rm dev 'bash -lc "chmod +x scripts/conformance.sh 2>/dev/null || true; ./scripts/conformance.sh"'

perf:
	@echo "🚀 Running performance tests in container..."
	@docker compose -f docker-compose.dev.yml run --rm dev 'bash -lc "chmod +x scripts/perf.sh 2>/dev/null || true; ./scripts/perf.sh"'

docs:
	@echo "📚 Building documentation in container..."
	@docker compose -f docker-compose.dev.yml run --rm dev 'bash -lc "chmod +x scripts/build-docs.sh 2>/dev/null || true; ./scripts/build-docs.sh npx"'

# Static analysis targets
lint:
	@echo "🔍 Running golangci-lint..."
	@golangci-lint run ./...

test-race:
	@echo "🏁 Running race tests..."
	@cd resolver-go && go test -race ./...
	@cd registrar-go && go test -race ./...

vet:
	@echo "🔬 Running go vet..."
	@cd resolver-go && go vet ./...
	@cd registrar-go && go vet ./...

qa:
	@echo "🔍 Running QA checks..."
	@chmod +x scripts/check-imports.sh 2>/dev/null || true && ./scripts/check-imports.sh
	$(MAKE) vet
	$(MAKE) lint
	$(MAKE) test-race

# ========================================================================
# Release workflow targets (LOCAL BUILDS)
# ========================================================================

dist-clean:
	@echo "🧹 Cleaning dist directory..."
	@rm -rf dist/*
	@mkdir -p dist/bin dist/images dist/docs dist/sbom dist/scan
	@touch dist/.gitkeep

binaries-local:
	@echo "🔨 Building cross-platform binaries..."
	@chmod +x scripts/build-binaries.sh 2>/dev/null || true
	@./scripts/build-binaries.sh

images-local:
	@echo "🐳 Building multi-arch Docker images..."
	@chmod +x scripts/build-images.sh 2>/dev/null || true
	@./scripts/build-images.sh

sbom-local:
	@echo "📋 Generating Software Bill of Materials..."
	@chmod +x scripts/sbom.local.sh 2>/dev/null || true
	@./scripts/sbom.local.sh

scan-local:
	@echo "🔐 Running vulnerability scans..."
	@chmod +x scripts/scan.local.sh 2>/dev/null || true
	@./scripts/scan.local.sh

docs-archive:
	@echo "📚 Packaging documentation..."
	@chmod +x scripts/docs-archive.sh 2>/dev/null || true
	@./scripts/docs-archive.sh

release-local:
	@echo "🚀 Creating complete local release..."
	$(MAKE) qa
	$(MAKE) binaries-local
	$(MAKE) images-local
	$(MAKE) docs-archive
	$(MAKE) sbom-local
	$(MAKE) scan-local
	@chmod +x scripts/release.local.sh 2>/dev/null || true
	@./scripts/release.local.sh

# ========================================================================
# SDK targets
# ========================================================================

sdk-test:
	@echo "🧪 Running SDK tests..."
	@docker compose -f docker-compose.dev.yml run --rm dev 'bash -lc "cd sdks/go/accdid && go mod init github.com/opendlt/accu-did/sdks/go/accdid 2>/dev/null || true && go test ./... -count=1 -v"'

sdk-merge-spec:
	@echo "🔄 Merging OpenAPI specifications..."
	@chmod +x scripts/sdk-openapi-merge.sh 2>/dev/null || true
	@./scripts/sdk-openapi-merge.sh

example-sdk:
	@echo "🚀 Running SDK example..."
	@docker compose -f docker-compose.dev.yml run --rm dev 'bash -lc "cd sdks/go/accdid/examples/basic && go mod init example 2>/dev/null || true && go mod edit -replace github.com/opendlt/accu-did/sdks/go/accdid=../.. && go mod tidy && go run main.go"'

# ========================================================================
# Devnet Management (Local Accumulate Development Network)
# ========================================================================

devnet-up:
	@echo "🚀 Starting local Accumulate devnet..."
	@pwsh -File scripts/devnet.ps1 up

devnet-down:
	@echo "🛑 Stopping local Accumulate devnet..."
	@pwsh -File scripts/devnet.ps1 down

devnet-status:
	@echo "📊 Checking devnet status..."
	@pwsh -File scripts/devnet.ps1 status

services-up: devnet-up
	@echo "🚀 Starting DID services against devnet..."
	@echo "Waiting for devnet to be ready..."
	@sleep 3
	@echo "Starting resolver on :8080..."
	@cd resolver-go && ACC_NODE_URL=http://127.0.0.1:26656 go run cmd/resolver/main.go --addr :8080 --real &
	@echo "Starting registrar on :8081..."
	@cd registrar-go && ACC_NODE_URL=http://127.0.0.1:26656 go run cmd/registrar/main.go --addr :8081 --real &
	@echo "✅ Services started. Use 'make services-down' to stop."
	@echo "   Resolver: http://localhost:8080/healthz"
	@echo "   Registrar: http://localhost:8081/healthz"
	@echo "   Devnet RPC: http://127.0.0.1:26656"

services-down:
	@echo "🛑 Stopping DID services..."
	@pkill -f "resolver.*--real" || true
	@pkill -f "registrar.*--real" || true
	@echo "✅ Services stopped."

sdk-example:
	@echo "🚀 Running SDK example against live devnet..."
	@echo "Prerequisites: Run 'make devnet-up' first"
	@cd examples/hello_accu && go run main.go

# ========================================================================
# TODO Scanner
# ========================================================================

todo-scan:
	@echo "🔍 Scanning repository for TODO markers..."
	@mkdir -p reports
	@if command -v docker >/dev/null 2>&1 && test -f docker-compose.dev.yml; then \
		docker compose -f docker-compose.dev.yml run --rm dev 'bash -lc "go run tools/todoscan/main.go . && echo \"\" && echo \"📊 Quick Summary:\" && test -f reports/todo-report.json && jq -r \".totalCount,(.summary.countsByTag | to_entries[] | \\\"  - \\(.key): \\(.value)\\\")\" reports/todo-report.json || echo \"Reports generated in ./reports/\""'; \
	elif command -v go >/dev/null 2>&1 && test -f tools/todoscan/main.go; then \
		go run tools/todoscan/main.go .; \
		echo ""; \
		echo "📊 Quick Summary:"; \
		if command -v jq >/dev/null 2>&1 && test -f reports/todo-report.json; then \
			jq -r '.totalCount,(.summary.countsByTag | to_entries[] | "  - \(.key): \(.value)")' reports/todo-report.json; \
		else \
			echo "Reports generated in ./reports/"; \
		fi; \
	else \
		echo "❌ Neither Docker+docker-compose.dev.yml nor Go+todoscan found"; \
		echo "Please ensure either:"; \
		echo "  1. Docker is available with docker-compose.dev.yml, OR"; \
		echo "  2. Go is installed locally with tools/todoscan/main.go"; \
		exit 1; \
	fi

help:
	@echo "📖 Accumulate DID Development Guide"
	@echo ""
	@echo "🐳 Container-first development (RECOMMENDED):"
	@echo "  dev-shell        - Launch interactive development container"
	@echo "  test-all         - Run all tests in container"
	@echo "  docs-container   - Build documentation in container"
	@echo "  ci-local         - Run complete local CI in container"
	@echo "  dev-up          - Start resolver and registrar services"
	@echo "  dev-down        - Stop development services"
	@echo "  check-imports   - Verify no forbidden imports"
	@echo "  conformance     - Run conformance tests"
	@echo "  perf           - Run performance tests"
	@echo ""
	@echo "🚀 Release workflow (LOCAL BUILDS):"
	@echo "  release-local   - Complete local release (QA + build + package + tag)"
	@echo "  binaries-local  - Cross-compile binaries for all platforms"
	@echo "  images-local    - Build multi-arch Docker images"
	@echo "  docs-archive    - Package documentation for distribution"
	@echo "  sbom-local      - Generate Software Bill of Materials"
	@echo "  scan-local      - Run vulnerability scans"
	@echo "  dist-clean      - Clean distribution directory"
	@echo ""
	@echo "📦 SDK development:"
	@echo "  sdk-test        - Run Go SDK tests in container"
	@echo "  sdk-merge-spec  - Merge OpenAPI specs for SDK generation"
	@echo "  example-sdk     - Run SDK example in container"
	@echo ""
	@echo "🌐 Local devnet (Live Accumulate blockchain):"
	@echo "  devnet-up       - Start local Accumulate devnet"
	@echo "  devnet-down     - Stop local Accumulate devnet"
	@echo "  devnet-status   - Check devnet health and endpoints"
	@echo "  services-up     - Start devnet + resolver + registrar in REAL mode"
	@echo "  services-down   - Stop all services"
	@echo "  sdk-example     - Run SDK example against live devnet"
	@echo ""
	@echo "🔍 Code analysis:"
	@echo "  todo-scan       - Scan repository for TODO/FIXME/XXX markers"
	@echo ""
	@echo "🔍 Static analysis (requires local tools):"
	@echo "  lint           - Run golangci-lint"
	@echo "  test-race      - Run race condition tests"
	@echo "  vet            - Run go vet"
	@echo "  qa             - Run all QA checks (imports, vet, lint, race)"
	@echo ""
	@echo "🖥️  Legacy targets (requires local tools):"
	@echo "  docs           - Build documentation (requires Node.js)"
	@echo "  test           - Run tests (requires Go)"
	@echo "  ci             - Run local CI script"
	@echo "  docker-build   - Build Docker images"
	@echo "  clean-docs     - Remove generated documentation"
	@echo "  clean-all      - Clean all build artifacts"
	@echo ""
	@echo "💡 Quick start: 'make dev-shell' for interactive development"
	@echo "💡 CI/CD ready: 'make ci-local' for complete validation"
	@echo "💡 Release ready: 'make release-local' for complete local release"