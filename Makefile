######################################################
# config
######################################################
# setting SHELL to bash allows bash commands to be executed by recipes
# options are set to exit when a recipe line exits non-zero or a piped command fails
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

## location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin

# golang ci linter binary and version
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.51.2

# tparse binary and version
TPARSE ?= $(LOCALBIN)/tparse
TPARSE_VERSION ?= latest


######################################################
# misc
######################################################
.PHONY: clean
clean:
	rm -rf build
	rm -rf $(LOCALBIN)
	rm -rf pkg/version/data/*.txt

.PHONY: localbin
localbin:
	mkdir -p $(LOCALBIN)

.PHONY: golangci-lint
golangci-lint: localbin
	test -s $(GOLANGCI_LINT)-$(GOLANGCI_LINT_VERSION) || GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	test -s $(GOLANGCI_LINT)-$(GOLANGCI_LINT_VERSION) || mv $(GOLANGCI_LINT) $(GOLANGCI_LINT)-$(GOLANGCI_LINT_VERSION)

.PHONY: tparse
tparse: localbin
	test -s $(TPARSE) || GOBIN=$(LOCALBIN) go install github.com/mfridman/tparse@$(TPARSE_VERSION)
	test -s $(TPARSE)-$(TPARSE_VERSION) || mv $(TPARSE) $(TPARSE)-$(TPARSE_VERSION)

######################################################
# go
######################################################
.PHONY: tidy
tidy: ## clean up go.mod and go.sum
	go mod tidy

.PHONY: download
download: ## downloads the dependencies
	go mod download -x

.PHONY: build
build: generate ## build kontext binary.
	go build -o build/kontext cmd/main.go

.PHONY: run
run: ## run kontext
	go run ./cmd/main.go


######################################################
# lint
######################################################
.PHONY: lint
lint: golangci-lint generate ## lint all code with golangci-lint
	$(GOLANGCI_LINT)-$(GOLANGCI_LINT_VERSION) run ./... --timeout 15m0s -v


######################################################
# test
######################################################
.PHONY: test
test: generate tparse
	set -eu
	set -o pipefail
	go test ./... -cover -json | $(TPARSE)-$(TPARSE_VERSION) -all

######################################################
# generate
######################################################
.PHONY: generate
generate: clean
	go generate github.com/orbatschow/kontext/pkg/version