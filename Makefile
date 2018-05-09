.DEFAULT_GOAL := build
.PHONY: dep build install docker dockerpush go_test truffle_test ethereum_test

REPO := smartcontract/chainlink
LDFLAGS := -ldflags "-X github.com/smartcontractkit/chainlink/store.Sha=`git rev-parse HEAD`"

HTTP_LIB := ./adapters/http/target/release/libhttp.dylib
LIBS := $(HTTP_LIB)

dep: ## Ensure chainlink's go dependencies are installed.
	@dep ensure

build: dep $(HTTP_LIB) ## Build chainlink.
	@go build $(LDFLAGS) -o chainlink

install: dep ## Install chainlink
	@go install $(LDFLAGS)

docker: ## Build the docker image.
	@docker build . -t $(REPO)

dockerpush: ## Push the docker image to dockerhub
	@docker push $(REPO)

$(HTTP_LIB): adapters/http/Cargo.toml adapters/http/src/*
	cargo build --release --manifest-path $<

go_test: $(HTTP_LIB)
	internal/ci/go_test

truffle_test: $(HTTP_LIB)
	internal/ci/truffle_test

ethereum_test: $(HTTP_LIB)
	internal/ci/ethereum_test

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
