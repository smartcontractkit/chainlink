.DEFAULT_GOAL := build
.PHONY: dep build install gui docker dockerpush

ENVIRONMENT ?= release

COMMIT_SHA ?= $(shell git rev-parse HEAD)
REPO=smartcontract/chainlink
LDFLAGS=-X github.com/smartcontractkit/chainlink/store.Sha=$(COMMIT_SHA)
GOFLAGS=-ldflags "$(LDFLAGS)"

# SGX is disabled by default, but turned on when building from Docker
SGX_ENABLED ?= no
SGX_SIMULATION ?= yes
SGX_TARGET := ./sgx/target/$(ENVIRONMENT)/
SGX_ENCLAVE := enclave.signed.so

ifeq ($(SGX_ENABLED),yes)
	LDFLAGS += -L $(SGX_TARGET)
	GOFLAGS += -tags=sgx_enclave
	SGX_BUILD_ENCLAVE := $(SGX_ENCLAVE)
else
	SGX_BUILD_ENCLAVE :=
endif

dep: ## Ensure chainlink's go dependencies are installed.
	dep ensure -vendor-only

build: dep gui $(SGX_BUILD_ENCLAVE) ## Build chainlink.
	ENVIRONMENT=$(ENVIRONMENT) go build $(GOFLAGS) -o chainlink

install: dep gui ## Install chainlink
	@go install $(GOFLAGS)

gui: ## Install GUI
	@cd gui
	yarn install
	@cd ..
	yarn build
	go generate ./...

docker: ## Build the docker image.
	docker build \
		--build-arg ENVIRONMENT=$(ENVIRONMENT) \
		--build-arg COMMIT_SHA=$(COMMIT_SHA) \
		--build-arg SGX_ENABLED=$(SGX_ENABLED) \
		--build-arg SGX_SIMULATION=$(SGX_SIMULATION) \
		-t $(REPO) .

dockerpush: ## Push the docker image to dockerhub
	docker push $(REPO)

.PHONY: $(SGX_ENCLAVE)
$(SGX_ENCLAVE):
	@SGX_SIMULATION=$(SGX_SIMULATION) make -C sgx/

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
