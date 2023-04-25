.DEFAULT_GOAL := chainlink

COMMIT_SHA ?= $(shell git rev-parse HEAD)
VERSION = $(shell cat VERSION)
GO_LDFLAGS := $(shell tools/bin/ldflags)
GOFLAGS = -ldflags "$(GO_LDFLAGS)"

.PHONY: install
install: operator-ui-autoinstall install-chainlink-autoinstall ## Install chainlink and all its dependencies.

.PHONY: install-git-hooks
install-git-hooks: ## Install git hooks.
	git config core.hooksPath .githooks

.PHONY: install-chainlink-autoinstall
install-chainlink-autoinstall: | pnpmdep gomod install-chainlink ## Autoinstall chainlink.
.PHONY: operator-ui-autoinstall
operator-ui-autoinstall: | operator-ui ## Autoinstall frontend UI.

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

chainlink: operator-ui ## Build the chainlink binary.
	go build $(GOFLAGS) .

.PHONY: docker ## Build the chainlink docker image
docker:
	docker buildx build \
	--build-arg COMMIT_SHA=$(COMMIT_SHA) \
	-f core/chainlink.Dockerfile .

.PHONY: operator-ui
operator-ui: ## Fetch the frontend
	./operator_ui/install.sh

.PHONY: abigen
abigen: ## Build & install abigen.
	./tools/bin/build_abigen

.PHONY: go-solidity-wrappers
go-solidity-wrappers: pnpmdep abigen ## Recompiles solidity contracts and their go wrappers.
	./contracts/scripts/native_solc_compile_all
	go generate ./core/gethwrappers

.PHONY: go-solidity-wrappers-transmission
go-solidity-wrappers-transmission: pnpmdep abigen ## Recompiles solidity contracts and their go wrappers.
	./contracts/scripts/transmission/native_solc_compile_all_transmission
	go generate ./core/gethwrappers/transmission

.PHONY: go-solidity-wrappers-ocr2vrf
go-solidity-wrappers-ocr2vrf: pnpmdep abigen ## Recompiles solidity contracts and their go wrappers.
	./contracts/scripts/native_solc_compile_all_ocr2vrf
	# replace the go:generate_disabled directive with the regular go:generate directive
	sed -i '' 's/go:generate_disabled/go:generate/g' core/gethwrappers/ocr2vrf/go_generate.go
	go generate ./core/gethwrappers/ocr2vrf
	go generate ./core/internal/mocks
	# put the go:generate_disabled directive back
	sed -i '' 's/go:generate/go:generate_disabled/g' core/gethwrappers/ocr2vrf/go_generate.go

.PHONY: generate
generate: abigen codecgen mockery ## Execute all go:generate commands.
	go generate -x ./...

.PHONY: testscripts
testscripts: chainlink ## Install and run testscript against testdata/scripts/* files.
	go install github.com/rogpeppe/go-internal/cmd/testscript@latest
	PATH=$(CURDIR):$(PATH) testscript -e CL_DEV=true -e COMMIT_SHA=$(COMMIT_SHA) -e VERSION=$(VERSION) $(TS_FLAGS) testdata/scripts/*

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
	go install github.com/vektra/mockery/v2@v2.22.1

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

.PHONY: test_need_operator_assets
test_need_operator_assets: ## Add blank file in web assets if operator ui has not been built
	[ -f "./core/web/assets/index.html" ] || mkdir ./core/web/assets && touch ./core/web/assets/index.html

.PHONY: config-docs
config-docs: ## Generate core node configuration documentation
	go run ./core/config/v2/docs/cmd/generate -o ./docs/

.PHONY: golangci-lint
golangci-lint: ## Run golangci-lint for all issues.
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.52.1 golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 > golangci-lint-output.txt

.PHONY: snapshot
snapshot:
	cd ./contracts && forge snapshot --match-test _gas

GORELEASER_CONFIG ?= .goreleaser.yaml

.PHONY: goreleaser-dev-build
goreleaser-dev-build: ## Run goreleaser snapshot build
	./tools/bin/goreleaser_wrapper build --snapshot --rm-dist --config ${GORELEASER_CONFIG}

.PHONY: goreleaser-dev-release
goreleaser-dev-release: ## run goreleaser snapshot release
	./tools/bin/goreleaser_wrapper release --snapshot --rm-dist --config ${GORELEASER_CONFIG}

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
