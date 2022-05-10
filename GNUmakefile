.DEFAULT_GOAL := build

GOPATH ?= $(HOME)/go
COMMIT_SHA ?= $(shell git rev-parse HEAD)
VERSION = $(shell cat VERSION)
GOBIN ?= $(GOPATH)/bin
GO_LDFLAGS := $(shell tools/bin/ldflags)
GOFLAGS = -ldflags "$(GO_LDFLAGS)"

.PHONY: install
install: operator-ui-autoinstall install-chainlink-autoinstall ## Install chainlink and all its dependencies.

.PHONY: install-git-hooks
install-git-hooks: ## Install git hooks.
	git config core.hooksPath .githooks

.PHONY: install-chainlink-autoinstall
install-chainlink-autoinstall: | gomod install-chainlink ## Autoinstall chainlink.
.PHONY: operator-ui-autoinstall
operator-ui-autoinstall: | yarndep operator-ui ## Autoinstall frontend UI.

.PHONY: gomod
gomod: ## Ensure chainlink's go dependencies are installed.
	@if [ -z "`which gencodec`" ]; then \
		go install github.com/smartcontractkit/gencodec@latest; \
	fi || true
	go mod download

.PHONY: yarndep
yarndep: ## Ensure all yarn dependencies are installed.
	yarn install --frozen-lockfile --prefer-offline

.PHONY: install-chainlink
install-chainlink: chainlink ## Install the chainlink binary.
	mkdir -p $(GOBIN)
	cp $< $(GOBIN)/chainlink

chainlink: operator-ui ## Build the chainlink binary.
	go build $(GOFLAGS) -o $@ ./core/

.PHONY: chainlink-build
chainlink-build: ## Build & install the chainlink binary.
	go build $(GOFLAGS) -o chainlink ./core/
	cp chainlink $(GOBIN)/chainlink

.PHONY: operator-ui
operator-ui: ## Build the static frontend UI.
	yarn setup:chainlink
	CHAINLINK_VERSION="$(VERSION)@$(COMMIT_SHA)" yarn workspace @chainlink/operator-ui build

.PHONY: contracts-operator-ui-build
contracts-operator-ui-build: # Only compiles tsc and builds contracts and operator-ui.
	yarn setup:chainlink
	CHAINLINK_VERSION="$(VERSION)@$(COMMIT_SHA)" yarn workspace @chainlink/operator-ui build

.PHONY: abigen
abigen: ## Build & install abigen.
	./tools/bin/build_abigen

.PHONY: go-solidity-wrappers
go-solidity-wrappers: tools/bin/abigen ## Recompiles solidity contracts and their go wrappers.
	./contracts/scripts/native_solc_compile_all
	go generate ./core/internal/gethwrappers

.PHONY: testdb
testdb: ## Prepares the test database.
	go run ./core/main.go local db preparetest

.PHONY: testdb
testdb-user-only: ## Prepares the test database with user only.
	go run ./core/main.go local db preparetest --user-only

# Format for CI
.PHONY: presubmit
presubmit: ## Format go files and imports.
	goimports -w ./core
	gofmt -w ./core
	go mod tidy

.PHONY: mockery
mockery: $(mockery) ## Install mockery.
	go install github.com/vektra/mockery/v2@v2.10.1

.PHONY: telemetry-protobuf
telemetry-protobuf: $(telemetry-protobuf) ## Generate telemetry protocol buffers.
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-wsrpc_out=. \
	--go-wsrpc_opt=paths=source_relative \
	./core/services/synchronization/telem/*.proto

.PHONY: test_install_ginkgo
test_install_ginkgo: ## Install ginkgo executable to run tests
	go install github.com/onsi/ginkgo/v2/ginkgo@v$(shell cat ./.tool-versions | grep ginkgo | sed -En "s/ginkgo.(.*)/\1/p")

.PHONY: test_smoke
test_smoke: # Run integration smoke tests.
	ginkgo -v -r --junit-report=tests-smoke-report.xml \
	--keep-going --trace --randomize-all --randomize-suites \
	--progress $(args) ./integration-tests/smoke


help:
	@echo ""
	@echo "         .__           .__       .__  .__        __"
	@echo "    ____ |  |__ _____  |__| ____ |  | |__| ____ |  | __"
	@echo "  _/ ___\|  |  \\\\\\__  \ |  |/    \|  | |  |/    \|  |/ /"
	@echo "  \  \___|   Y  \/ __ \|  |   |  \  |_|  |   |  \    <"
	@echo "   \___  >___|  (____  /__|___|  /____/__|___|  /__|_ \\"
	@echo "       \/     \/     \/        \/             \/     \/"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
