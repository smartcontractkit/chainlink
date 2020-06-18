.DEFAULT_GOAL := build

ENVIRONMENT ?= release

REPO := smartcontract/chainlink
COMMIT_SHA ?= $(shell git rev-parse HEAD)
VERSION = $(shell cat VERSION)
GOBIN ?= $(GOPATH)/bin
GO_LDFLAGS := $(shell tools/bin/ldflags)
GOFLAGS = -ldflags "$(GO_LDFLAGS)"
DOCKERFILE := core/chainlink.Dockerfile

# SGX is disabled by default, but turned on when building from Docker
SGX_ENABLED ?= no
SGX_SIMULATION ?= yes
SGX_ENCLAVE := enclave.signed.so
SGX_TARGET := ./core/sgx/target/$(ENVIRONMENT)/

ifneq (,$(filter yes true,$(SGX_ENABLED)))
	GOFLAGS += -tags=sgx_enclave
	SGX_BUILD_ENCLAVE := $(SGX_ENCLAVE)
	DOCKERFILE := core/chainlink-sgx.Dockerfile
	REPO := $(REPO)-sgx
else
	SGX_BUILD_ENCLAVE :=
endif

TAGGED_REPO := $(REPO):$(VERSION)

.PHONY: install
install: operator-ui-autoinstall install-chainlink-autoinstall ## Install chainlink and all its dependencies.

.PHONY: install-git-hooks
install-git-hooks:
	git config core.hooksPath .githooks

.PHONY: install-chainlink-autoinstall
install-chainlink-autoinstall: | gomod install-chainlink
.PHONY: operator-ui-autoinstall
operator-ui-autoinstall: | yarndep operator-ui

.PHONY: gomod
gomod: ## Ensure chainlink's go dependencies are installed.
	@if [ -z "`which gencodec`" ]; then \
		go get github.com/smartcontractkit/gencodec; \
	fi || true
	go mod download

.PHONY: yarndep
yarndep: ## Ensure all yarn dependencies are installed
	yarn install --frozen-lockfile
	./tools/bin/restore-solc-cache

.PHONY: gen-builder-cache
gen-builder-cache: gomod # generate a cache for the builder image
	yarn install --frozen-lockfile
	./tools/bin/restore-solc-cache

.PHONY: install-chainlink
install-chainlink: chainlink ## Install the chainlink binary.
	cp $< $(GOBIN)/chainlink

chainlink: $(SGX_BUILD_ENCLAVE) operator-ui ## Build the chainlink binary.
	CGO_ENABLED=0 go run packr/main.go "${CURDIR}/core/services/eth" ## embed contracts in .go file
	go build $(GOFLAGS) -o $@ ./core/

.PHONY: chainlink-build
chainlink-build:
	CGO_ENABLED=0 go run packr/main.go "${CURDIR}/core/eth" ## embed contracts in .go file
	CGO_ENABLED=0 go run packr/main.go "${CURDIR}/core/services"
	go build $(GOFLAGS) -o chainlink ./core/
	cp chainlink $(GOBIN)/chainlink

.PHONY: operator-ui
operator-ui: ## Build the static frontend UI.
	yarn setup:chainlink
	CHAINLINK_VERSION="$(VERSION)@$(COMMIT_SHA)" yarn workspace @chainlink/operator-ui build
	CGO_ENABLED=0 go run packr/main.go "${CURDIR}/core/services"

.PHONY: contracts-operator-ui-build
contracts-operator-ui-build: # only compiles tsc and builds contracts and operator-ui
	yarn setup:chainlink
	CHAINLINK_VERSION="$(VERSION)@$(COMMIT_SHA)" yarn workspace @chainlink/operator-ui build

.PHONY: abigen
abigen:
	./tools/bin/build_abigen

.PHONY: go-solidity-wrappers
go-solidity-wrappers: abigen ## Recompiles solidity contracts and their go wrappers
	yarn workspace @chainlink/contracts compile
	go generate ./core/internal/gethwrappers
	go run ./packr/main.go ./core/eth/

.PHONY: testdb
testdb: ## Prepares the test database
	go run ./core/main.go local db preparetest

.PHONY: docker
docker: ## Build the docker image.
	docker build \
		--build-arg ENVIRONMENT=$(ENVIRONMENT) \
		--build-arg COMMIT_SHA=$(COMMIT_SHA) \
		--build-arg SGX_SIMULATION=$(SGX_SIMULATION) \
		-t $(TAGGED_REPO) \
		-t $(REPO):$(COMMIT_SHA) \
		-f $(DOCKERFILE) \
		.

.PHONY: dockerpush
dockerpush: ## Push the docker image to dockerhub
	docker push $(TAGGED_REPO)

.PHONY: $(SGX_ENCLAVE)
$(SGX_ENCLAVE):
	@ENVIRONMENT=$(ENVIRONMENT) SGX_ENABLED=$(SGX_ENABLED) SGX_SIMULATION=$(SGX_SIMULATION) make -C core/sgx/
	@ln -f $(SGX_TARGET)/libadapters.so core/sgx/target/libadapters.so

.PHONY: enclave
enclave: $(SGX_ENCLAVE)

help:
	@echo ""
	@echo "         .__           .__       .__  .__        __"
	@echo "    ____ |  |__ _____  |__| ____ |  | |__| ____ |  | __"
	@echo "  _/ ___\|  |  \\\\\\__  \ |  |/    \|  | |  |/    \|  |/ /"
	@echo "  \  \___|   Y  \/ __ \|  |   |  \  |_|  |   |  \    <"
	@echo "   \___  >___|  (____  /__|___|  /____/__|___|  /__|_ \\"
	@echo "       \/     \/     \/        \/             \/     \/"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
