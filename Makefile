.PHONY: all build clean run

# Variables
BINARY_NAME_LISTENER=rayls-listener
BINARY_NAME_FLAGGER=rayls-flagger
BINARY_NAME_GOVERNANCE_API=rayls-governance-api

GOLANGCI_LINT_VERSION=v2.7.1

SRC_DIR_LISTENER=cmd/listener
SRC_DIR_FLAGGER=cmd/flagger
SRC_DIR_API=cmd/api

MAIN_FILE_LISTENER=$(SRC_DIR_LISTENER)/main.go
MAIN_FILE_FLAGGER=$(SRC_DIR_FLAGGER)/main.go
MAIN_FILE_API=$(SRC_DIR_API)/main.go
BUILD_DIR=build

# Targets
all: build

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME_LISTENER) $(MAIN_FILE_LISTENER)
	go build -o $(BUILD_DIR)/$(BINARY_NAME_FLAGGER) $(MAIN_FILE_FLAGGER)
	go build -o $(BUILD_DIR)/$(BINARY_NAME_GOVERNANCE_API) $(MAIN_FILE_API)

listener:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME_LISTENER) $(MAIN_FILE_LISTENER)

flagger:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME_FLAGGER) $(MAIN_FILE_FLAGGER)

api: swagger
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME_GOVERNANCE_API) $(MAIN_FILE_API)

swagger:
	swag init --parseDependency -q -g ./cmd/api/main.go -o ./cmd/api/docs

clean:
	rm -rf $(BUILD_DIR)

run: build
	./$(BUILD_DIR)/$(BINARY_NAME_GOVERNANCE_API)
	./$(BUILD_DIR)/$(BINARY_NAME_LISTENER)
	./$(BUILD_DIR)/$(BINARY_NAME_FLAGGER)

# Install all Go linters and code-quality tools used by the project
install-linters:
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@GOBIN=$${GOBIN:-$$(go env GOPATH)/bin}; \
	mkdir -p "$$GOBIN"; \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh \
		| sh -s -- -b "$$GOBIN" $(GOLANGCI_LINT_VERSION)
	@echo ""
	@echo "  ✓ golangci-lint installed successfully!"

# Install and activate pre-commit so every contributor runs the same checks locally
install-precommit:
	@command -v pre-commit >/dev/null 2>&1 || { \
		echo "Installing pre-commit..."; \
		if command -v pipx >/dev/null 2>&1; then \
			pipx install pre-commit; \
		elif command -v brew >/dev/null 2>&1; then \
			brew install pre-commit; \
		else \
			pip install --user pre-commit; \
		fi \
	}
	pre-commit install
	@echo ""
	@echo "  ✓ pre-commit installed successfully!"

# Full environment setup: installs linters and registers the project's Git hooks
setup-linters: install-linters install-precommit
	@echo ""
	@echo "  ✓ Setup ready to Go! :)"

# Run all linters and formatters checks (usage: make lint [target])
# Examples: make lint . | make lint ./cmd/api | make lint ./cmd/flagger/core/ports.go
lint:
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Error: golangci-lint is not installed"; \
		echo "Run: make install-tools"; \
		exit 1; \
	fi
	@TARGET="$(or $(filter-out $@,$(MAKECMDGOALS)),.)"; \
	if [ -d "$$TARGET" ]; then \
		LINT_TARGET="$$TARGET/..."; \
	else \
		LINT_TARGET="$$TARGET"; \
	fi; \
	if [ "$$TARGET" = "." ]; then \
		LINT_TARGET="./..."; \
		echo "Running golangci-lint on entire codebase..."; \
	else \
		echo "Running golangci-lint on $$TARGET..."; \
	fi; \
	echo "Running formatters..."; \
	golangci-lint fmt $$LINT_TARGET; \
	echo "Running linters..."; \
	golangci-lint run $$LINT_TARGET
listener-coverage:
	@echo "Running listener business logic coverage tests..."
	@go test -coverprofile=coverage_core.out ./cmd/listener/core > /dev/null 2>&1 || true
	@go test -coverprofile=coverage_services.out ./cmd/listener/core/services > /dev/null 2>&1 || true
	@go test -coverprofile=coverage_handlers.out ./cmd/listener/core/handlers > /dev/null 2>&1 || true
	@go test -coverprofile=coverage_combined.out ./cmd/listener/core ./cmd/listener/core/services ./cmd/listener/core/handlers > /dev/null 2>&1 || true
	@echo "Business Logic Test Coverage:"
	@core_cov=$$(go tool cover -func=coverage_core.out 2>/dev/null | grep total | awk '{print $$NF}'); \
	services_cov=$$(go tool cover -func=coverage_services.out 2>/dev/null | grep total | awk '{print $$NF}'); \
	handlers_cov=$$(go tool cover -func=coverage_handlers.out 2>/dev/null | grep total | awk '{print $$NF}'); \
	combined_cov=$$(go tool cover -func=coverage_combined.out 2>/dev/null | grep total | awk '{print $$NF}'); \
	if [ -n "$$core_cov" ]; then echo "  Core:     $$core_cov"; fi; \
	if [ -n "$$services_cov" ]; then echo "  Services: $$services_cov"; fi; \
	if [ -n "$$handlers_cov" ]; then echo "  Handlers: $$handlers_cov"; fi; \
	if [ -n "$$combined_cov" ]; then echo "  Average:  $$combined_cov"; fi
	@rm -f coverage_core.out coverage_services.out coverage_handlers.out coverage_combined.out

flagger-coverage:
	@echo "Running flagger business logic coverage tests..."
	@go test -coverprofile=coverage_flagger_core.out ./cmd/flagger/core > /dev/null 2>&1 || true
	@echo "Business Logic Test Coverage:"
	@core_cov=$$(go tool cover -func=coverage_flagger_core.out 2>/dev/null | grep total | awk '{print $$NF}'); \
	if [ -n "$$core_cov" ]; then echo "  Core: $$core_cov"; else echo "  No coverage data available"; fi
	@rm -f coverage_flagger_core.out

api-coverage:
	@echo "Running API business logic coverage tests..."
	@go test -coverprofile=coverage_api_core.out ./cmd/api/core > /dev/null 2>&1 || true
	@echo "Business Logic Test Coverage:"
	@core_cov=$$(go tool cover -func=coverage_api_core.out 2>/dev/null | grep total | awk '{print $$NF}'); \
	if [ -n "$$core_cov" ]; then echo "  Core: $$core_cov"; else echo "  No coverage data available"; fi
	@rm -f coverage_api_core.out
