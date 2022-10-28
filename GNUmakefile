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
operator-ui-autoinstall: | operator-ui ## Autoinstall frontend UI.

.PHONY: gomod
gomod: ## Ensure chainlink's go dependencies are installed.
	@if [ -z "`which gencodec`" ]; then \
		go install github.com/smartcontractkit/gencodec@latest; \
	fi || true
	go mod download

.PHONY: gomodtidy
gomodtidy: ## Run go mod tidy on all modules.
	go mod tidy
	cd ./integration-tests && go mod tidy

.PHONY: install-chainlink
install-chainlink: chainlink ## Install the chainlink binary.
	mkdir -p $(GOBIN)
	rm -f $(GOBIN)/chainlink
	cp $< $(GOBIN)/chainlink

chainlink: operator-ui ## Build the chainlink binary.
	go build $(GOFLAGS) -o $@ ./core/

.PHONY: docker ## Build the chainlink docker image
docker: 
	docker buildx build \
	--build-arg COMMIT_SHA=$(COMMIT_SHA) \
	-f core/chainlink.Dockerfile .

.PHONY: chainlink-build
chainlink-build: operator-ui ## Build & install the chainlink binary.
	go build $(GOFLAGS) -o chainlink ./core/
	rm -f $(GOBIN)/chainlink
	cp chainlink $(GOBIN)/chainlink

.PHONY: operator-ui
operator-ui: ## Fetch the frontend
	./operator_ui/install.sh	

.PHONY: abigen
abigen: ## Build & install abigen.
	./tools/bin/build_abigen

.PHONY: go-solidity-wrappers
go-solidity-wrappers: abigen ## Recompiles solidity contracts and their go wrappers.
	./contracts/scripts/native_solc_compile_all
	go generate ./core/gethwrappers

.PHONY: go-solidity-wrappers-ocr2vrf
go-solidity-wrappers-ocr2vrf: abigen ## Recompiles solidity contracts and their go wrappers.
	./contracts/scripts/native_solc_compile_all_ocr2vrf
	go generate ./core/gethwrappers/ocr2vrf

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
	go install github.com/vektra/mockery/v2@v2.14.0

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

.PHONY: test_need_operator_assets
test_need_operator_assets: ## Add blank file in web assets if operator ui has not been built
	[ -f "./core/web/assets/index.html" ] || mkdir ./core/web/assets && touch ./core/web/assets/index.html

.PHONY: test_smoke
test_smoke: test_need_operator_assets ## Run all integration smoke tests, using only simulated networks, default behavior
	ginkgo -v -r --junit-report=tests-smoke-report.xml \
	--keep-going --trace --randomize-all --randomize-suites \
	--progress $(args) ./integration-tests/smoke

.PHONY: test_smoke_simulated
test_smoke_simulated: test_need_operator_assets ## Run all integration smoke tests, using only simulated networks, default behavior (you can use `make test_smoke`)
	SELECTED_NETWORKS="SIMULATED" ginkgo -v -r --junit-report=tests-smoke-report.xml \
	--keep-going --trace --randomize-all --randomize-suites \
	--progress $(args) ./integration-tests/smoke

.PHONY: test_soak_ocr
test_soak_ocr: test_need_operator_assets ## Run the OCR soak test
	cd ./integration-tests && go test -v -run ^TestOCRSoak$$ ./soak -count=1 && cd ..

.PHONY: test_soak_forwarder_ocr
test_soak_forwarder_ocr: test_need_operator_assets ## Run the Forwarder OCR soak test
	cd ./integration-tests && go test -v -run ^TestForwarderOCRSoak$$ ./soak -count=1 && cd ..

.PHONY: test_soak_keeper
test_soak_keeper: test_need_operator_assets ## Run the OCR soak test
	cd ./integration-tests && go test -v -run ^TestKeeperSoak$$ ./soak -count=1 && cd ..

.PHONY: test_perf
test_perf: test_need_operator_assets ## Run core node performance tests.
	ginkgo -v -r --junit-report=tests-perf-report.xml \
	--keep-going --trace --randomize-all --randomize-suites \
	--progress $(args) ./integration-tests/performance

.PHONY: test_chaos
test_chaos: test_need_operator_assets ## Run core node chaos tests.
	ginkgo -r --focus @chaos --nodes 5 ./integration-tests/chaos

.PHONY: config-docs
config-docs: ## Generate core node configuration documentation
	go run ./core/config/v2/docs/cmd/generate/main.go -o ./docs/

.PHONY: golangci-lint
golangci-lint: ## Run golangci-lint for all issues.
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 > golangci-lint-output.txt

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
