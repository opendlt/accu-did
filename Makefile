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

.PHONY: docs clean-docs open-docs

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