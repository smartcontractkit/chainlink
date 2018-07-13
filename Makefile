.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := build
.PHONY: godep yarndep build install gui docker dockerpush

ENVIRONMENT ?= release

REPO := smartcontract/chainlink
COMMIT_SHA ?= $(shell git rev-parse HEAD)
VERSION = $(shell cat VERSION)
GO_LDFLAGS := $(shell internal/bin/ldflags)
GOFLAGS := -ldflags "$(GO_LDFLAGS)"

# SGX is disabled by default, but turned on when building from Docker
SGX_ENABLED ?= no
SGX_SIMULATION ?= yes
SGX_TARGET := ./sgx/target/$(ENVIRONMENT)/
SGX_ENCLAVE := enclave.signed.so

ifeq ($(SGX_ENABLED),yes)
	GO_LDFLAGS += -L $(SGX_TARGET)
	GOFLAGS += -tags=sgx_enclave
	SGX_BUILD_ENCLAVE := $(SGX_ENCLAVE)
else
	SGX_BUILD_ENCLAVE :=
endif

godep: ## Ensure chainlink's go dependencies are installed.
	dep ensure -vendor-only

yarndep: ## Ensure the frontend's dependencies are installed.
	yarn install

build: yarndep godep chainlink ## Build chainlink.

install: godep gui ## Install chainlink
	go install $(GOFLAGS)

gui: ## Install GUI
	yarn build
	go generate ./...

docker: ## Build the docker image.
	docker build \
		--build-arg ENVIRONMENT \
		--build-arg COMMIT_SHA \
		--build-arg SGX_ENABLED \
		--build-arg SGX_SIMULATION \
		-t $(REPO) .

dockerpush: ## Push the docker image to dockerhub
	docker push $(REPO)

chainlink: gui $(SGX_BUILD_ENCLAVE)
	go build $(GOFLAGS) -o chainlink

.PHONY: $(SGX_ENCLAVE)
$(SGX_ENCLAVE):
	@make -C sgx/

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
