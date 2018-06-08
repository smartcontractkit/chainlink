.DEFAULT_GOAL := build
.PHONY: dep build install gui docker dockerpush

REPO=smartcontract/chainlink
LDFLAGS=-ldflags "-X github.com/smartcontractkit/chainlink/store.Sha=`git rev-parse HEAD`"

dep: ## Ensure chainlink's go dependencies are installed.
	@dep ensure

build: dep ## Build chainlink.
	@go build $(LDFLAGS) -o chainlink

install: dep ## Install chainlink
	@go install $(LDFLAGS)

gui: ## Install GUI 
	@cd gui
	@yarn install
	@cd ..
	@yarn build

docker: ## Build the docker image.
	@docker build . -t $(REPO)

dockerpush: ## Push the docker image to dockerhub
	@docker push $(REPO)

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
