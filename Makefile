TOOLS_BIN_DIR ?= $(shell pwd)/tmp/bin

export PATH := $(TOOLS_BIN_DIR):$(PATH)

GOLANGCILINTER_BINARY=$(TOOLS_BIN_DIR)/golangci-lint
MDOX_BINARY=$(TOOLS_BIN_DIR)/mdox
MDOX_VALIDATE_CONFIG?=.mdox.validate.yaml

TOOLING=$(MDOX_BINARY) $(GOLANGCILINTER_BINARY)

MD_FILES_TO_FORMAT=$(shell ls *.md)

GOCMD=go
GOMAIN=main.go
GOBUILD=$(GOCMD) build
GOOS?=$(shell go env GOOS)
ENVVARS=GOOS=$(GOOS) CGO_ENABLED=0

# Dashboard build configuration with defaults
OUTPUT_DIR_OPERATOR ?= ./built/dashboards/operator
OUTPUT_DIR_PERSES ?= ./built/dashboards/perses
OUTPUT_FORMAT_PERSES ?= json
PROJECT ?= default
DATASOURCE ?= prometheus-datasource

.PHONY: demo
start-demo:
	@echo "Setting up demo environment"

	@cd ./examples && docker-compose up -d

.PHONY: clean-demo
clean-demo:
	@echo "Cleaning up demo environment"

	@cd ./examples && docker-compose down -v

.PHONY: build-dashboards
build-dashboards:
	@echo "Removing old dashboards"
	@rm -rf ./examples/dashboards/
	@echo "Building dashboards"
	@$(ENVVARS) $(GOCMD) run $(GOMAIN) --output-dir="./examples/dashboards/operator" --output="operator" --project="perses-dev" --datasource="prometheus-datasource"
	@$(ENVVARS) $(GOCMD) run $(GOMAIN) --output-dir="./examples/dashboards/perses" --output="yaml" --project="perses-dev" --datasource="prometheus-datasource"
	@$(ENVVARS) $(GOCMD) run $(GOMAIN) --output-dir="./examples/dashboards/json/operator/" --output="operator-json" --project="perses-dev" --datasource="prometheus-datasource"

# Adding a new target for building and testing dashboards locally with configurable flags
.PHONY: build-dashboards-local
build-dashboards-local:
	@echo "Removing old dashboards"
	@rm -rf $(OUTPUT_DIR_OPERATOR)
	@rm -rf $(OUTPUT_DIR_PERSES)
	@echo "Building dashboards for local testing"
	@$(ENVVARS) $(GOCMD) run $(GOMAIN) --output-dir=$(OUTPUT_DIR_OPERATOR) --output="operator" --project=$(PROJECT) --datasource=$(DATASOURCE)
	@$(ENVVARS) $(GOCMD) run $(GOMAIN) --output-dir=$(OUTPUT_DIR_PERSES) --output=$(OUTPUT_FORMAT_PERSES) --project=$(PROJECT) --datasource=$(DATASOURCE)

.PHONY: deps
deps:
	$(ENVVARS) $(GOCMD) mod download

.PHONY: fmt
fmt:
	$(ENVVARS) $(GOCMD) fmt -x ./...

.PHONY: vet
vet:
	$(ENVVARS) $(GOCMD) vet ./...

.PHONY: check-golang
check-golang: $(GOLANGCILINTER_BINARY)
	$(GOLANGCILINTER_BINARY) run

.PHONY: fix-golang
fix-golang: $(GOLANGCILINTER_BINARY)
	$(GOLANGCILINTER_BINARY) run --fix

.PHONY: docs
docs: $(MDOX_BINARY)
	@echo ">> formatting and local/remote link check"
	$(MDOX_BINARY) fmt --soft-wraps -l --links.validate.config-file=$(MDOX_VALIDATE_CONFIG) $(MD_FILES_TO_FORMAT)

.PHONY: check-docs
check-docs: $(MDOX_BINARY)
	@echo ">> checking formatting and local/remote links"
	$(MDOX_BINARY) fmt --soft-wraps --check -l --links.validate.config-file=$(MDOX_VALIDATE_CONFIG) $(MD_FILES_TO_FORMAT)

.PHONY: tidy
tidy:
	go mod tidy -v
	cd scripts && go mod tidy -v -modfile=go.mod -compat=1.18

all: fmt vet deps check-golang check-docs

define require_clean_work_tree
	@git update-index -q --ignore-submodules --refresh

	@if ! git diff-files --quiet --ignore-submodules --; then \
		echo >&2 "cannot $1: you have unstaged changes."; \
		git diff -r --ignore-submodules -- >&2; \
		echo >&2 "Please commit or stash them."; \
		exit 1; \
	fi

	@if ! git diff-index --cached --quiet HEAD --ignore-submodules --; then \
		echo >&2 "cannot $1: your index contains uncommitted changes."; \
		git diff --cached -r --ignore-submodules HEAD -- >&2; \
		echo >&2 "Please commit or stash them."; \
		exit 1; \
	fi

endef

.PHONY: generate
generate: tidy deps fmt build-dashboards
	$(call require_clean_work_tree,'all generated files should be committed, run make generate and commit changes.')

$(TOOLS_BIN_DIR):
	mkdir -p $(TOOLS_BIN_DIR)

$(TOOLING): $(TOOLS_BIN_DIR)
	@echo Installing tools from scripts/tools.go
	@cat scripts/tools.go | grep _ | awk -F'"' '{print $$2}' | GOBIN=$(TOOLS_BIN_DIR) xargs -tI % go install -mod=readonly -modfile=scripts/go.mod %
