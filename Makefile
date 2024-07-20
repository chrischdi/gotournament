# Copyright 2024 Christian Schlotter.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Ensure Make is run with bash shell as some syntax below is bash-specific
SHELL:=/usr/bin/env bash

.DEFAULT_GOAL:=help

#
# Directories.
#
# Full directory of where the Makefile resides
BIN_DIR := bin
TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(abspath $(TOOLS_DIR)/$(BIN_DIR))
GO_INSTALL := ./scripts/go_install.sh

export PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)

#
# Go.
#
GO_VERSION ?= 1.22.5
GO_DIRECTIVE_VERSION ?= 1.22.0
GO_CONTAINER_IMAGE ?= docker.io/library/golang:$(GO_VERSION)

GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT_VER := $(shell cat .github/workflows/pr-golangci-lint.yaml | grep [[:space:]]version: | sed 's/.*version: //')
GOLANGCI_LINT := $(abspath $(TOOLS_BIN_DIR)/$(GOLANGCI_LINT_BIN)-$(GOLANGCI_LINT_VER))
GOLANGCI_LINT_PKG := github.com/golangci/golangci-lint/cmd/golangci-lint

GOVULNCHECK_BIN := govulncheck
GOVULNCHECK_VER := v1.0.4
GOVULNCHECK := $(abspath $(TOOLS_BIN_DIR)/$(GOVULNCHECK_BIN)-$(GOVULNCHECK_VER))
GOVULNCHECK_PKG := golang.org/x/vuln/cmd/govulncheck

ARCH ?= $(shell go env GOARCH)
ALL_PLATFORM ?= darwin/amd64 linux/amd64 linux/arm linux/arm64 windows/386 windows/amd64

# Set build time variables including version details
LDFLAGS := $(shell scripts/version.sh)

all: test build

help:  # Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-50s\033[0m %s\n", $$1, $$2 } /^\$$\([0-9A-Za-z_-]+\):.*?##/ { gsub("_","-", $$1); printf "  \033[36m%-50s\033[0m %s\n", tolower(substr($$1, 3, length($$1)-7)), $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## --------------------------------------
## Generate
## --------------------------------------

##@ generate:

.PHONY: generate-modules
generate-modules: ## Run go mod tidy to ensure modules are up to date
	go mod tidy

## --------------------------------------
## Lint / Verify
## --------------------------------------

##@ lint and verify:

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint the codebase
	$(GOLANGCI_LINT) run -v $(GOLANGCI_LINT_EXTRA_ARGS)

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT) ## Lint the codebase and run auto-fixers if supported by the linter
	GOLANGCI_LINT_EXTRA_ARGS=--fix $(MAKE) lint

ALL_VERIFY_CHECKS = modules security

.PHONY: verify
verify: $(addprefix verify-,$(ALL_VERIFY_CHECKS)) lint-dockerfiles ## Run all verify-* targets

.PHONY: verify-modules
verify-modules: generate-modules  ## Verify go modules are up to date
	@if !(git diff --quiet HEAD -- go.sum go.mod); then \
		git diff; \
		echo "go module files are out of date"; exit 1; \
	fi

.PHONY: verify-govulncheck
verify-govulncheck: $(GOVULNCHECK) ## Verify code for vulnerabilities
	$(GOVULNCHECK) ./... && R1=$$? || R1=$$?; \
	if [ "$$R1" -ne "0" ]; then \
		exit 1; \
	fi

.PHONY: verify-security
verify-security: ## Verify code and images for vulnerabilities
	$(MAKE) verify-govulncheck && R2=$$? || R2=$$?; \
	if [ "$$R1" -ne "0" ] || [ "$$R2" -ne "0" ]; then \
	  echo "Check for vulnerabilities failed! There are vulnerabilities to be fixed"; \
		exit 1; \
	fi

## --------------------------------------
## Binaries
## --------------------------------------

##@ build:

.PHONY: build
build: ## Build the gotournament binary
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/gotournament ./cmd

## --------------------------------------
## Testing
## --------------------------------------

##@ test:

.PHONY: test
test: ## Run unit and integration tests
	go test ./... $(TEST_ARGS)

.PHONY: test-verbose
test-verbose: ## Run unit and integration tests with verbose flag
	$(MAKE) test TEST_ARGS="$(TEST_ARGS) -v"

.PHONY: test-cover
test-cover: ## Run unit and integration tests and generate a coverage report
	$(MAKE) test TEST_ARGS="$(TEST_ARGS) -coverprofile=out/coverage.out"
	go tool cover -func=out/coverage.out -o out/coverage.txt
	go tool cover -html=out/coverage.out -o out/coverage.html

## --------------------------------------
## Release
## --------------------------------------

##@ release:

## latest git tag for the commit, e.g., v0.3.10
RELEASE_TAG ?= $(shell git describe --abbrev=0 2>/dev/null)
ifneq (,$(findstring -,$(RELEASE_TAG)))
    PRE_RELEASE=true
endif
RELEASE_DIR := out

.PHONY: release
release: clean-release ## Build and push container images using the latest git tag for the commit
	@if [ -z "${RELEASE_TAG}" ]; then echo "RELEASE_TAG is not set"; exit 1; fi
	@if ! [ -z "$$(git status --porcelain)" ]; then echo "Your local git repository contains uncommitted changes, use git clean before proceeding."; exit 1; fi
	git checkout "${RELEASE_TAG}"
	# Build binaries first.
	GIT_VERSION=$(RELEASE_TAG) $(MAKE) release-binaries

# .PHONY: release-binaries
# release-binaries: ## Build the binaries to publish with a release
# 	RELEASE_BINARY=clusterctl-linux-amd64 BUILD_PATH=./cmd/clusterctl GOOS=linux GOARCH=amd64 $(MAKE) release-binary
# 	RELEASE_BINARY=clusterctl-linux-arm64 BUILD_PATH=./cmd/clusterctl GOOS=linux GOARCH=arm64 $(MAKE) release-binary
# 	RELEASE_BINARY=clusterctl-darwin-amd64 BUILD_PATH=./cmd/clusterctl GOOS=darwin GOARCH=amd64 $(MAKE) release-binary
# 	RELEASE_BINARY=clusterctl-darwin-arm64 BUILD_PATH=./cmd/clusterctl GOOS=darwin GOARCH=arm64 $(MAKE) release-binary
# 	RELEASE_BINARY=clusterctl-windows-amd64.exe BUILD_PATH=./cmd/clusterctl GOOS=windows GOARCH=amd64 $(MAKE) release-binary
# 	RELEASE_BINARY=clusterctl-linux-ppc64le BUILD_PATH=./cmd/clusterctl GOOS=linux GOARCH=ppc64le $(MAKE) release-binary

# .PHONY: release-binary
# release-binary: $(RELEASE_DIR)
# 	docker run \
# 		--rm \
# 		-e CGO_ENABLED=0 \
# 		-e GOOS=$(GOOS) \
# 		-e GOARCH=$(GOARCH) \
# 		-e GOCACHE=/tmp/ \
# 		--user $$(id -u):$$(id -g) \
# 		-v "$$(pwd):/workspace$(DOCKER_VOL_OPTS)" \
# 		-w /workspace \
# 		golang:$(GO_VERSION) \
# 		go build -a -trimpath -ldflags "$(LDFLAGS) -extldflags '-static'" \
# 		-o $(RELEASE_DIR)/$(notdir $(RELEASE_BINARY)) $(BUILD_PATH)

## --------------------------------------
## Cleanup / Verification
## --------------------------------------

##@ clean:

.PHONY: clean
clean: ## Remove generated binaries, GitBook files, Helm charts, and Tilt build files
	$(MAKE) clean-bin

.PHONY: clean-bin
clean-bin: ## Remove all generated binaries
	rm -rf $(BIN_DIR)
	rm -rf $(TOOLS_BIN_DIR)

.PHONY: clean-release
clean-release: ## Remove the release folder
	rm -rf $(RELEASE_DIR)

## --------------------------------------
## Hack / Tools
## --------------------------------------

##@ hack/tools:
.PHONY: $(GOLANGCI_LINT_BIN)
$(GOLANGCI_LINT_BIN): $(GOLANGCI_LINT) ## Build a local copy of golangci-lint.

.PHONY: $(GOVULNCHECK_BIN)
$(GOVULNCHECK_BIN): $(GOVULNCHECK) ## Build a local copy of govulncheck.

## We are forcing a rebuilt of conversion-gen via PHONY so that we're always using an up-to-date version.
## We can't use a versioned name for the binary, because that would be reflected in generated files.
$(GOLANGCI_LINT): # Build golangci-lint from tools folder.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(GOLANGCI_LINT_PKG) $(GOLANGCI_LINT_BIN) $(GOLANGCI_LINT_VER)

$(GOVULNCHECK): # Build govulncheck.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(GOVULNCHECK_PKG) $(GOVULNCHECK_BIN) $(GOVULNCHECK_VER)

## --------------------------------------
## Helpers
## --------------------------------------

##@ helpers:

go-version: ## Print the go version we use to compile our binaries and images
	@echo $(GO_VERSION)