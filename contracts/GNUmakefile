.DEFAULT_GOAL := build

GOPATH ?= $(HOME)/go
COMMIT_SHA ?= $(shell git rev-parse HEAD)
VERSION = $(shell cat VERSION)
GOBIN ?= $(GOPATH)/bin
GO_LDFLAGS := $(shell tools/bin/ldflags)
GOFLAGS = -ldflags "$(GO_LDFLAGS)"


.PHONY: abigen
abigen: ## Build & install abigen.
	../tools/bin/build_abigen

.PHONY: wrappers
wrappers: abigen ## Recompiles solidity contracts and their go wrappers.
	./scripts/native_solc_compile_all_ccip
	go generate ../core/gethwrappers/go_generate_ccip.go

.PHONY: snapshot
snapshot:
	FOUNDRY_PROFILE= forge snapshot --match-test _gas

ccip-precommit: export FOUNDRY_PROFILE=ccip
.PHONY: ccip-precommit
ccip-precommit:
	forge test
	make snapshot
	pnpm prettier:write
	pnpm lint:ccip

ccip-lcov: export FOUNDRY_PROFILE=ccip
.PHONY: ccip-lcov
ccip-lcov:
	forge coverage --report lcov
	../tools/ci/ccip_lcov_prune ./lcov.info ./lcov.info.pruned
	genhtml -o report lcov.info.pruned --branch-coverage

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
