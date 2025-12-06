BINARY ?= converter
GOCACHE ?= $(PWD)/.gocache
IMAGE ?= ebook-pdf
GOLANGCI_LINT_VERSION ?= v2.7.1
GOLANGCI_LINT_PKG ?= github.com/golangci/golangci-lint/v2/cmd/golangci-lint

.PHONY: build test lint lint-install docker clean

build:
	GOCACHE=$(GOCACHE) go build -o $(BINARY) .

test:
	GOCACHE=$(GOCACHE) go test ./...

lint:
	@$(MAKE) lint-install
	@set -e; \
	TARGET="$(GOLANGCI_LINT_VERSION)"; \
	GOPATH_BIN=$$(go env GOPATH)/bin; \
	LINTER="$$GOPATH_BIN/golangci-lint"; \
	if [ ! -x "$$LINTER" ]; then \
		LINTER=$$(command -v golangci-lint || true); \
	fi; \
	FOUND=$$($$LINTER version 2>/dev/null | head -n1 | awk '{print $$4}'); \
	if [ -n "$$FOUND" ] && [ "$$FOUND" != "$$TARGET" ]; then \
		if printf '%s\n%s\n' "$$FOUND" "$$TARGET" | sort -V | tail -n1 | grep -qx "$$FOUND"; then \
			printf 'WARNING: golangci-lint %s detected; pinned version is %s.\n' "$$FOUND" "$$TARGET"; \
		fi; \
	fi; \
	GOCACHE=$(GOCACHE) $$LINTER run; \
	GOCACHE=$(GOCACHE) go test ./...

lint-install:
	@set -e; \
	GOPATH_BIN=$$(go env GOPATH)/bin; \
	LINTER="$$GOPATH_BIN/golangci-lint"; \
	TARGET="$(GOLANGCI_LINT_VERSION)"; \
	NEED_INSTALL=1; \
	if [ -x "$$LINTER" ]; then \
		FOUND=$$($$LINTER version 2>/dev/null | head -n1 | awk '{print $$4}'); \
		if [ "$$FOUND" = "$$TARGET" ]; then \
			NEED_INSTALL=0; \
		fi; \
	fi; \
	if [ $$NEED_INSTALL -eq 1 ]; then \
		printf 'Installing golangci-lint %s...\n' "$$TARGET"; \
		go install $(GOLANGCI_LINT_PKG)@$(GOLANGCI_LINT_VERSION); \
	fi

docker:
	docker build -t $(IMAGE) .

clean:
	rm -rf $(BINARY) $(GOCACHE)
