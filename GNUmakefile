.DEFAULT_GOAL := chainlink

COMMIT_SHA ?= $(shell git rev-parse HEAD)
VERSION = $(shell cat VERSION)
GO_LDFLAGS := $(shell tools/bin/ldflags)
GOFLAGS = -ldflags "$(GO_LDFLAGS)"

.PHONY: install
install: install-chainlink-autoinstall ## Install chainlink and all its dependencies.

.PHONY: install-git-hooks
install-git-hooks: ## Install git hooks.
	git config core.hooksPath .githooks

.PHONY: install-chainlink-autoinstall
install-chainlink-autoinstall: | pnpmdep gomod install-chainlink ## Autoinstall chainlink.

.PHONY: pnpmdep
pnpmdep: ## Install solidity contract dependencies through pnpm
	(cd contracts && pnpm i)

.PHONY: gomod
gomod: ## Ensure chainlink's go dependencies are installed.
	@if [ -z "`which gencodec`" ]; then \
		go install github.com/smartcontractkit/gencodec@latest; \
	fi || true
	go mod download

.PHONY: gomodtidy
gomodtidy: ## Run go mod tidy on all modules.
	go mod tidy
	cd ./core/scripts && go mod tidy
	cd ./integration-tests && go mod tidy

.PHONY: godoc
godoc: ## Install and run godoc
	go install golang.org/x/tools/cmd/godoc@latest
	# http://localhost:6060/pkg/github.com/smartcontractkit/chainlink/v2/
	godoc -http=:6060

.PHONY: install-chainlink
install-chainlink: operator-ui ## Install the chainlink binary.
	go install $(GOFLAGS) .

chainlink: ## Build the chainlink binary.
	go build $(GOFLAGS) .

chainlink-dev: ## Build a dev build of chainlink binary.
	go build -tags dev $(GOFLAGS) .

chainlink-test: ## Build a test build of chainlink binary.
	go build $(GOFLAGS) .

.PHONY: chainlink-local-start
chainlink-local-start:
	./chainlink -c /etc/node-secrets-volume/default.toml -c /etc/node-secrets-volume/overrides.toml -secrets /etc/node-secrets-volume/secrets.toml node start -d -p /etc/node-secrets-volume/node-password -a /etc/node-secrets-volume/apicredentials --vrfpassword=/etc/node-secrets-volume/apicredentials

.PHONY: install-medianpoc
install-medianpoc: ## Build & install the chainlink-medianpoc binary.
	go install $(GOFLAGS) ./plugins/cmd/chainlink-medianpoc

.PHONY: docker ## Build the chainlink docker image
docker:
	docker buildx build \
	--build-arg COMMIT_SHA=$(COMMIT_SHA) \
	-f core/chainlink.Dockerfile .

.PHONY: docker-plugins ## Build the chainlink-plugins docker image
docker-plugins:
	docker buildx build \
	--build-arg COMMIT_SHA=$(COMMIT_SHA) \
	-f plugins/chainlink.Dockerfile .

.PHONY: operator-ui
operator-ui: ## Fetch the frontend
	go generate ./core/web

.PHONY: abigen
abigen: ## Build & install abigen.
	./tools/bin/build_abigen

.PHONY: generate
generate: abigen codecgen mockery ## Execute all go:generate commands.
	go generate -x ./...

.PHONY: testscripts
testscripts: chainlink-test ## Install and run testscript against testdata/scripts/* files.
	go install github.com/rogpeppe/go-internal/cmd/testscript@latest
	go run ./tools/txtar/cmd/lstxtardirs -recurse=true | PATH="$(CURDIR):${PATH}" xargs -I % \
		sh -c 'testscript -e COMMIT_SHA=$(COMMIT_SHA) -e HOME="$(TMPDIR)/home" -e VERSION=$(VERSION) $(TS_FLAGS) %/*.txtar'

.PHONY: testscripts-update
testscripts-update: ## Update testdata/scripts/* files via testscript.
	make testscripts TS_FLAGS="-u"

.PHONY: testdb
testdb: ## Prepares the test database.
	go run . local db preparetest

.PHONY: testdb
testdb-user-only: ## Prepares the test database with user only.
	go run . local db preparetest --user-only

# Format for CI
.PHONY: presubmit
presubmit: ## Format go files and imports.
	goimports -w .
	gofmt -w .
	go mod tidy

.PHONY: mockery
mockery: $(mockery) ## Install mockery.
	go install github.com/vektra/mockery/v2@v2.38.0

.PHONY: codecgen
codecgen: $(codecgen) ## Install codecgen
	go install github.com/ugorji/go/codec/codecgen@v1.2.10

.PHONY: telemetry-protobuf
telemetry-protobuf: $(telemetry-protobuf) ## Generate telemetry protocol buffers.
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-wsrpc_out=. \
	--go-wsrpc_opt=paths=source_relative \
	./core/services/synchronization/telem/*.proto

.PHONY: config-docs
config-docs: ## Generate core node configuration documentation
	go run ./core/config/docs/cmd/generate -o ./docs/

.PHONY: golangci-lint
golangci-lint: ## Run golangci-lint for all issues.
	[ -d "./golangci-lint" ] || mkdir ./golangci-lint && \
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.55.2 golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 > ./golangci-lint/$(shell date +%Y-%m-%d_%H:%M:%S).txt


GORELEASER_CONFIG ?= .goreleaser.yaml

.PHONY: goreleaser-dev-build
goreleaser-dev-build: ## Run goreleaser snapshot build
	./tools/bin/goreleaser_wrapper build --snapshot --rm-dist --config ${GORELEASER_CONFIG}

.PHONY: goreleaser-dev-release
goreleaser-dev-release: ## run goreleaser snapshot release
	./tools/bin/goreleaser_wrapper release --snapshot --rm-dist --config ${GORELEASER_CONFIG}

.PHONY: modgraph
modgraph:
	./tools/bin/modgraph > go.md

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
