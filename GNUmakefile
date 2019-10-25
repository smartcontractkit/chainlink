.DEFAULT_GOAL := build

ENVIRONMENT ?= release

REPO := smartcontract/chainlink
COMMIT_SHA ?= $(shell git rev-parse HEAD)
VERSION = $(shell cat VERSION)
GOBIN ?= $(GOPATH)/bin
GO_LDFLAGS := $(shell tools/bin/ldflags)
GOFLAGS = -ldflags "$(GO_LDFLAGS)"
DOCKERFILE := Dockerfile
DOCKER_TAG ?= latest

# SGX is disabled by default, but turned on when building from Docker
SGX_ENABLED ?= no
SGX_SIMULATION ?= yes
SGX_ENCLAVE := enclave.signed.so
SGX_TARGET := ./sgx/target/$(ENVIRONMENT)/

ifneq (,$(filter yes true,$(SGX_ENABLED)))
	GOFLAGS += -tags=sgx_enclave
	SGX_BUILD_ENCLAVE := $(SGX_ENCLAVE)
	DOCKERFILE := Dockerfile-sgx
	REPO := $(REPO)-sgx
else
	SGX_BUILD_ENCLAVE :=
endif

TAGGED_REPO := $(REPO):$(DOCKER_TAG)

.PHONY: install
install: operator-ui-autoinstall install-chainlink-autoinstall ## Install chainlink and all its dependencies.

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
yarndep: ## Ensure the frontend's dependencies are installed.
	yarn install --frozen-lockfile

.PHONY: install-chainlink
install-chainlink: chainlink ## Install the chainlink binary.
	cp $< $(GOBIN)/chainlink

chainlink: $(SGX_BUILD_ENCLAVE) operator-ui ## Build the chainlink binary.
	go build $(GOFLAGS) -o $@ ./core/

.PHONY: operator-ui
operator-ui: ## Build the static frontend UI.
	CHAINLINK_VERSION="$(VERSION)@$(COMMIT_SHA)" yarn workspace @chainlink/operator-ui run build
	CGO_ENABLED=0 go run operator_ui/main.go "${CURDIR}/core/services"

.PHONY: docker
docker: ## Build the docker image.
	docker build \
		--build-arg ENVIRONMENT=$(ENVIRONMENT) \
		--build-arg COMMIT_SHA=$(COMMIT_SHA) \
		--build-arg SGX_SIMULATION=$(SGX_SIMULATION) \
		-t $(TAGGED_REPO) \
		-f $(DOCKERFILE) \
		.

.PHONY: dockerpush
dockerpush: ## Push the docker image to dockerhub
	docker push $(TAGGED_REPO)

.PHONY: $(SGX_ENCLAVE)
$(SGX_ENCLAVE):
	@ENVIRONMENT=$(ENVIRONMENT) SGX_ENABLED=$(SGX_ENABLED) SGX_SIMULATION=$(SGX_SIMULATION) make -C sgx/
	@ln -f $(SGX_TARGET)/libadapters.so sgx/target/libadapters.so

.PHONY: enclave
enclave: $(SGX_ENCLAVE)

.PHONY: lint
lint: ## Execute yarn lint for the core project
	yarn lint

.PHONY: format
format: ## Execute yarn format for the core project
	yarn format

.PHONY: setup
setup: ## Execute yarn setup for the core project
	yarn setup

.PHONY: evm/generate-typings
evm/generate-typings: ## Execute yarn generate-typings for evm subproject
	cd ./evm/ && yarn generate-typings

.PHONY: evm/build
evm/build: ## Execute yarn build for evm subproject
	cd ./evm/ && yarn build

.PHONY: evm/postbuild
evm/postbuild: ## Execute yarn postbuild for evm subproject
	cd ./evm/ && yarn postbuild

.PHONY: evm/build\:windows
evm/build\:windows: ## Execute yarn build:windows for evm subproject
	cd ./evm/ && yarn build:windows

.PHONY: evm/depcheck
evm/depcheck: ## Execute yarn depcheck for evm subproject
	cd ./evm/ && yarn depcheck

.PHONY: evm/eslint
evm/eslint: ## Execute yarn eslint for evm subproject
	cd ./evm/ && yarn eslint

.PHONY: evm/solhint
evm/solhint: ## Execute yarn solhint for evm subproject
	cd ./evm/ && yarn solhint

.PHONY: evm/lint
evm/lint: ## Execute yarn lint for evm subproject
	cd ./evm/ && yarn lint

.PHONY: evm/slither
evm/slither: ## Execute yarn slither for evm subproject
	cd ./evm/ && yarn slither

.PHONY: evm/pretest
evm/pretest: ## Execute yarn pretest for evm subproject
	cd ./evm/ && yarn pretest

.PHONY: evm/test\:v1
evm/test\:v1: ## Execute yarn test:v1 for evm subproject
	cd ./evm/ && yarn test:v1

.PHONY: evm/test\:v2
evm/test\:v2: ## Execute yarn test:v2 for evm subproject
	cd ./evm/ && yarn test:v2

.PHONY: evm/test
evm/test: ## Execute yarn test for evm subproject
	cd ./evm/ && yarn test

.PHONY: evm/format
evm/format: ## Execute yarn format for evm subproject
	cd ./evm/ && yarn format

.PHONY: evm/prepublishOnly
evm/prepublishOnly: ## Execute yarn prepublishOnly for evm subproject
	cd ./evm/ && yarn prepublishOnly

.PHONY: evm/setup
evm/setup: ## Execute yarn setup for evm subproject
	cd ./evm/ && yarn setup

.PHONY: evm/truffle\:migrate\:cldev
evm/truffle\:migrate\:cldev: ## Execute yarn truffle:migrate:cldev for evm subproject
	cd ./evm/ && yarn truffle:migrate:cldev

.PHONY: box/compile
box/compile: ## Execute yarn compile for box subproject
	cd ./evm/box/ && yarn compile

.PHONY: box/console\:dev
box/console\:dev: ## Execute yarn console:dev for box subproject
	cd ./evm/box/ && yarn console:dev

.PHONY: box/console\:live
box/console\:live: ## Execute yarn console:live for box subproject
	cd ./evm/box/ && yarn console:live

.PHONY: box/depcheck
box/depcheck: ## Execute yarn depcheck for box subproject
	cd ./evm/box/ && yarn depcheck

.PHONY: box/eslint
box/eslint: ## Execute yarn eslint for box subproject
	cd ./evm/box/ && yarn eslint

.PHONY: box/solhint
box/solhint: ## Execute yarn solhint for box subproject
	cd ./evm/box/ && yarn solhint

.PHONY: box/lint
box/lint: ## Execute yarn lint for box subproject
	cd ./evm/box/ && yarn lint

.PHONY: box/format
box/format: ## Execute yarn format for box subproject
	cd ./evm/box/ && yarn format

.PHONY: box/migrate\:dev
box/migrate\:dev: ## Execute yarn migrate:dev for box subproject
	cd ./evm/box/ && yarn migrate:dev

.PHONY: box/migrate\:live
box/migrate\:live: ## Execute yarn migrate:live for box subproject
	cd ./evm/box/ && yarn migrate:live

.PHONY: box/setup
box/setup: ## Execute yarn setup for box subproject
	cd ./evm/box/ && yarn setup

.PHONY: box/test
box/test: ## Execute yarn test for box subproject
	cd ./evm/box/ && yarn test

.PHONY: v0.5/build
v0.5/build: ## Execute yarn build for v0.5 subproject
	cd ./evm/v0.5/ && yarn build

.PHONY: v0.5/build.windows
v0.5/build.windows: ## Execute yarn build.windows for v0.5 subproject
	cd ./evm/v0.5/ && yarn build.windows

.PHONY: v0.5/depcheck
v0.5/depcheck: ## Execute yarn depcheck for v0.5 subproject
	cd ./evm/v0.5/ && yarn depcheck

.PHONY: v0.5/eslint
v0.5/eslint: ## Execute yarn eslint for v0.5 subproject
	cd ./evm/v0.5/ && yarn eslint

.PHONY: v0.5/solhint
v0.5/solhint: ## Execute yarn solhint for v0.5 subproject
	cd ./evm/v0.5/ && yarn solhint

.PHONY: v0.5/lint
v0.5/lint: ## Execute yarn lint for v0.5 subproject
	cd ./evm/v0.5/ && yarn lint

.PHONY: v0.5/format
v0.5/format: ## Execute yarn format for v0.5 subproject
	cd ./evm/v0.5/ && yarn format

.PHONY: v0.5/slither
v0.5/slither: ## Execute yarn slither for v0.5 subproject
	cd ./evm/v0.5/ && yarn slither

.PHONY: v0.5/setup
v0.5/setup: ## Execute yarn setup for v0.5 subproject
	cd ./evm/v0.5/ && yarn setup

.PHONY: v0.5/test
v0.5/test: ## Execute yarn test for v0.5 subproject
	cd ./evm/v0.5/ && yarn test

.PHONY: echo_server/depcheck
echo_server/depcheck: ## Execute yarn depcheck for echo_server subproject
	cd ./examples/echo_server/ && yarn depcheck

.PHONY: echo_server/eslint
echo_server/eslint: ## Execute yarn eslint for echo_server subproject
	cd ./examples/echo_server/ && yarn eslint

.PHONY: echo_server/solhint
echo_server/solhint: ## Execute yarn solhint for echo_server subproject
	cd ./examples/echo_server/ && yarn solhint

.PHONY: echo_server/lint
echo_server/lint: ## Execute yarn lint for echo_server subproject
	cd ./examples/echo_server/ && yarn lint

.PHONY: echo_server/format
echo_server/format: ## Execute yarn format for echo_server subproject
	cd ./examples/echo_server/ && yarn format

.PHONY: echo_server/setup
echo_server/setup: ## Execute yarn setup for echo_server subproject
	cd ./examples/echo_server/ && yarn setup

.PHONY: echo_server/test
echo_server/test: ## Execute yarn test for echo_server subproject
	cd ./examples/echo_server/ && yarn test

.PHONY: testnet/depcheck
testnet/depcheck: ## Execute yarn depcheck for testnet subproject
	cd ./examples/testnet/ && yarn depcheck

.PHONY: testnet/eslint
testnet/eslint: ## Execute yarn eslint for testnet subproject
	cd ./examples/testnet/ && yarn eslint

.PHONY: testnet/solhint
testnet/solhint: ## Execute yarn solhint for testnet subproject
	cd ./examples/testnet/ && yarn solhint

.PHONY: testnet/lint
testnet/lint: ## Execute yarn lint for testnet subproject
	cd ./examples/testnet/ && yarn lint

.PHONY: testnet/format
testnet/format: ## Execute yarn format for testnet subproject
	cd ./examples/testnet/ && yarn format

.PHONY: testnet/setup
testnet/setup: ## Execute yarn setup for testnet subproject
	cd ./examples/testnet/ && yarn setup

.PHONY: twilio_sms/depcheck
twilio_sms/depcheck: ## Execute yarn depcheck for twilio_sms subproject
	cd ./examples/twilio_sms/ && yarn depcheck

.PHONY: twilio_sms/eslint
twilio_sms/eslint: ## Execute yarn eslint for twilio_sms subproject
	cd ./examples/twilio_sms/ && yarn eslint

.PHONY: twilio_sms/solhint
twilio_sms/solhint: ## Execute yarn solhint for twilio_sms subproject
	cd ./examples/twilio_sms/ && yarn solhint

.PHONY: twilio_sms/lint
twilio_sms/lint: ## Execute yarn lint for twilio_sms subproject
	cd ./examples/twilio_sms/ && yarn lint

.PHONY: twilio_sms/format
twilio_sms/format: ## Execute yarn format for twilio_sms subproject
	cd ./examples/twilio_sms/ && yarn format

.PHONY: twilio_sms/setup
twilio_sms/setup: ## Execute yarn setup for twilio_sms subproject
	cd ./examples/twilio_sms/ && yarn setup

.PHONY: twilio_sms/test
twilio_sms/test: ## Execute yarn test for twilio_sms subproject
	cd ./examples/twilio_sms/ && yarn test

.PHONY: uptime_sla/depcheck
uptime_sla/depcheck: ## Execute yarn depcheck for uptime_sla subproject
	cd ./examples/uptime_sla/ && yarn depcheck

.PHONY: uptime_sla/eslint
uptime_sla/eslint: ## Execute yarn eslint for uptime_sla subproject
	cd ./examples/uptime_sla/ && yarn eslint

.PHONY: uptime_sla/solhint
uptime_sla/solhint: ## Execute yarn solhint for uptime_sla subproject
	cd ./examples/uptime_sla/ && yarn solhint

.PHONY: uptime_sla/lint
uptime_sla/lint: ## Execute yarn lint for uptime_sla subproject
	cd ./examples/uptime_sla/ && yarn lint

.PHONY: uptime_sla/format
uptime_sla/format: ## Execute yarn format for uptime_sla subproject
	cd ./examples/uptime_sla/ && yarn format

.PHONY: uptime_sla/setup
uptime_sla/setup: ## Execute yarn setup for uptime_sla subproject
	cd ./examples/uptime_sla/ && yarn setup

.PHONY: uptime_sla/test
uptime_sla/test: ## Execute yarn test for uptime_sla subproject
	cd ./examples/uptime_sla/ && yarn test

.PHONY: explorer/admin\:seed
explorer/admin\:seed: ## Execute yarn admin:seed for explorer subproject
	cd ./explorer/ && yarn admin:seed

.PHONY: explorer/admin\:clnodes\:add
explorer/admin\:clnodes\:add: ## Execute yarn admin:clnodes:add for explorer subproject
	cd ./explorer/ && yarn admin:clnodes:add

.PHONY: explorer/admin\:clnodes\:delete
explorer/admin\:clnodes\:delete: ## Execute yarn admin:clnodes:delete for explorer subproject
	cd ./explorer/ && yarn admin:clnodes:delete

.PHONY: explorer/depcheck
explorer/depcheck: ## Execute yarn depcheck for explorer subproject
	cd ./explorer/ && yarn depcheck

.PHONY: explorer/predev
explorer/predev: ## Execute yarn predev for explorer subproject
	cd ./explorer/ && yarn predev

.PHONY: explorer/dev
explorer/dev: ## Execute yarn dev for explorer subproject
	cd ./explorer/ && yarn dev

.PHONY: explorer/dev\:client
explorer/dev\:client: ## Execute yarn dev:client for explorer subproject
	cd ./explorer/ && yarn dev:client

.PHONY: explorer/dev\:server
explorer/dev\:server: ## Execute yarn dev:server for explorer subproject
	cd ./explorer/ && yarn dev:server

.PHONY: explorer/build
explorer/build: ## Execute yarn build for explorer subproject
	cd ./explorer/ && yarn build

.PHONY: explorer/prod
explorer/prod: ## Execute yarn prod for explorer subproject
	cd ./explorer/ && yarn prod

.PHONY: explorer/pretest
explorer/pretest: ## Execute yarn pretest for explorer subproject
	cd ./explorer/ && yarn pretest

.PHONY: explorer/test
explorer/test: ## Execute yarn test for explorer subproject
	cd ./explorer/ && yarn test

.PHONY: explorer/test-ci
explorer/test-ci: ## Execute yarn test-ci for explorer subproject
	cd ./explorer/ && yarn test-ci

.PHONY: explorer/test-ci\:e2e
explorer/test-ci\:e2e: ## Execute yarn test-ci:e2e for explorer subproject
	cd ./explorer/ && yarn test-ci:e2e

.PHONY: explorer/test-ci\:e2e\:no-build
explorer/test-ci\:e2e\:no-build: ## Execute yarn test-ci:e2e:no-build for explorer subproject
	cd ./explorer/ && yarn test-ci:e2e:no-build

.PHONY: explorer/test-ci\:silent
explorer/test-ci\:silent: ## Execute yarn test-ci:silent for explorer subproject
	cd ./explorer/ && yarn test-ci:silent

.PHONY: explorer/test-ci\:e2e\:silent
explorer/test-ci\:e2e\:silent: ## Execute yarn test-ci:e2e:silent for explorer subproject
	cd ./explorer/ && yarn test-ci:e2e:silent

.PHONY: explorer/lint
explorer/lint: ## Execute yarn lint for explorer subproject
	cd ./explorer/ && yarn lint

.PHONY: explorer/lint\:fix
explorer/lint\:fix: ## Execute yarn lint:fix for explorer subproject
	cd ./explorer/ && yarn lint:fix

.PHONY: explorer/format
explorer/format: ## Execute yarn format for explorer subproject
	cd ./explorer/ && yarn format

.PHONY: explorer/migration\:run
explorer/migration\:run: ## Execute yarn migration:run for explorer subproject
	cd ./explorer/ && yarn migration:run

.PHONY: explorer/migration\:revert
explorer/migration\:revert: ## Execute yarn migration:revert for explorer subproject
	cd ./explorer/ && yarn migration:revert

.PHONY: explorer/test\:migration\:run
explorer/test\:migration\:run: ## Execute yarn test:migration:run for explorer subproject
	cd ./explorer/ && yarn test:migration:run

.PHONY: explorer/test\:migration\:revert
explorer/test\:migration\:revert: ## Execute yarn test:migration:revert for explorer subproject
	cd ./explorer/ && yarn test:migration:revert

.PHONY: explorer/automigrate
explorer/automigrate: ## Execute yarn automigrate for explorer subproject
	cd ./explorer/ && yarn automigrate

.PHONY: explorer/setup
explorer/setup: ## Execute yarn setup for explorer subproject
	cd ./explorer/ && yarn setup

.PHONY: explorer/dockerpush
explorer/dockerpush: ## Execute make dockerpush for explorer subproject
	cd ./explorer/ && make dockerpush

.PHONY: integration/cypressJobServer
integration/cypressJobServer: ## Execute yarn cypressJobServer for integration subproject
	cd ./integration/ && yarn cypressJobServer

.PHONY: integration/depcheck
integration/depcheck: ## Execute yarn depcheck for integration subproject
	cd ./integration/ && yarn depcheck

.PHONY: integration/eslint
integration/eslint: ## Execute yarn eslint for integration subproject
	cd ./integration/ && yarn eslint

.PHONY: integration/format
integration/format: ## Execute yarn format for integration subproject
	cd ./integration/ && yarn format

.PHONY: integration/lint
integration/lint: ## Execute yarn lint for integration subproject
	cd ./integration/ && yarn lint

.PHONY: integration/setup
integration/setup: ## Execute yarn setup for integration subproject
	cd ./integration/ && yarn setup

.PHONY: integration/test
integration/test: ## Execute yarn test for integration subproject
	cd ./integration/ && yarn test

.PHONY: integration/test\:cypress
integration/test\:cypress: ## Execute yarn test:cypress for integration subproject
	cd ./integration/ && yarn test:cypress

.PHONY: integration/test\:forks
integration/test\:forks: ## Execute yarn test:forks for integration subproject
	cd ./integration/ && yarn test:forks

.PHONY: forks/build_geth_image
forks/build_geth_image: ## Execute build_geth_image for forks subproject
	cd ./integration/forks/ && make build_geth_image

.PHONY: forks/build_chainlink_image
forks/build_chainlink_image: ## Execute build_chainlink_image for forks subproject
	cd ./integration/forks/ && make build_chainlink_image

.PHONY: forks/start_network
forks/start_network: ## Execute make start_network for forks subproject
	cd ./integration/forks/ && make start_network

.PHONY: forks/tear_down
forks/tear_down: ## Execute make tear_down for forks subproject
	cd ./integration/forks/ && make tear_down

.PHONY: forks/initial_setup
forks/initial_setup: ## Execute make initial_setup for forks subproject
	cd ./integration/forks/ && make initial_setup

.PHONY: forks/create_job
forks/create_job: ## Execute make create_job for forks subproject
	cd ./integration/forks/ && make create_job

.PHONY: forks/create_curl_script
forks/create_curl_script: ## Execute make create_curl_script for forks subproject
	cd ./integration/forks/ && make create_curl_script

.PHONY: forks/run_chain_1
forks/run_chain_1: ## Execute make run_chain_1 for forks subproject
	cd ./integration/forks/ && make run_chain_1

.PHONY: forks/run_chain_2
forks/run_chain_2: ## Execute make run_chain_2 for forks subproject
	cd ./integration/forks/ && make run_chain_2

.PHONY: integration-scripts/generate-typings
integration-scripts/generate-typings: ## Execute yarn generate-typings for integration-scripts subproject
	cd ./integration-scripts/ && yarn generate-typings

.PHONY: integration-scripts/build\:contracts
integration-scripts/build\:contracts: ## Execute yarn build:contracts for integration-scripts subproject
	cd ./integration-scripts/ && yarn build:contracts

.PHONY: integration-scripts/postbuild\:contracts
integration-scripts/postbuild\:contracts: ## Execute yarn postbuild:contracts for integration-scripts subproject
	cd ./integration-scripts/ && yarn postbuild:contracts

.PHONY: integration-scripts/prebuild
integration-scripts/prebuild: ## Execute yarn prebuild for integration-scripts subproject
	cd ./integration-scripts/ && yarn prebuild

.PHONY: integration-scripts/build
integration-scripts/build: ## Execute yarn build for integration-scripts subproject
	cd ./integration-scripts/ && yarn build

.PHONY: integration-scripts/setup
integration-scripts/setup: ## Execute yarn setup for integration-scripts subproject
	cd ./integration-scripts/ && yarn setup

.PHONY: integration-scripts/format
integration-scripts/format: ## Execute yarn format for integration-scripts subproject
	cd ./integration-scripts/ && yarn format

.PHONY: integration-scripts/lint
integration-scripts/lint: ## Execute yarn lint for integration-scripts subproject
	cd ./integration-scripts/ && yarn lint

.PHONY: integration-scripts/count-transaction-events
integration-scripts/count-transaction-events: ## Execute yarn count-transaction-events for integration-scripts subproject
	cd ./integration-scripts/ && yarn count-transaction-events

.PHONY: integration-scripts/send-runlog-transaction
integration-scripts/send-runlog-transaction: ## Execute yarn send-runlog-transaction for integration-scripts subproject
	cd ./integration-scripts/ && yarn send-runlog-transaction

.PHONY: integration-scripts/send-ethlog-transaction
integration-scripts/send-ethlog-transaction: ## Execute yarn send-ethlog-transaction for integration-scripts subproject
	cd ./integration-scripts/ && yarn send-ethlog-transaction

.PHONY: integration-scripts/fund-address
integration-scripts/fund-address: ## Execute yarn fund-address for integration-scripts subproject
	cd ./integration-scripts/ && yarn fund-address

.PHONY: integration-scripts/deploy-contracts
integration-scripts/deploy-contracts: ## Execute yarn deploy-contracts for integration-scripts subproject
	cd ./integration-scripts/ && yarn deploy-contracts

.PHONY: integration-scripts/start-echo-server
integration-scripts/start-echo-server: ## Execute yarn start-echo-server for integration-scripts subproject
	cd ./integration-scripts/ && yarn start-echo-server

.PHONY: operator_ui/start
operator_ui/start: ## Execute yarn start for operator_ui subproject
	cd ./operator_ui/ && yarn start

.PHONY: operator_ui/build
operator_ui/build: ## Execute yarn build for operator_ui subproject
	cd ./operator_ui/ && yarn build

.PHONY: operator_ui/build\:tsc
operator_ui/build\:tsc: ## Execute yarn build:tsc for operator_ui subproject
	cd ./operator_ui/ && yarn build:tsc

.PHONY: operator_ui/build\:tsc\:clean
operator_ui/build\:tsc\:clean: ## Execute yarn build:tsc:clean for operator_ui subproject
	cd ./operator_ui/ && yarn build:tsc:clean

.PHONY: operator_ui/serve
operator_ui/serve: ## Execute yarn serve for operator_ui subproject
	cd ./operator_ui/ && yarn serve

.PHONY: operator_ui/prelint
operator_ui/prelint: ## Execute yarn prelint for operator_ui subproject
	cd ./operator_ui/ && yarn prelint

.PHONY: operator_ui/eslint
operator_ui/eslint: ## Execute yarn eslint for operator_ui subproject
	cd ./operator_ui/ && yarn eslint

.PHONY: operator_ui/lint
operator_ui/lint: ## Execute yarn lint for operator_ui subproject
	cd ./operator_ui/ && yarn lint

.PHONY: operator_ui/format
operator_ui/format: ## Execute yarn format for operator_ui subproject
	cd ./operator_ui/ && yarn format

.PHONY: operator_ui/pretest
operator_ui/pretest: ## Execute yarn pretest for operator_ui subproject
	cd ./operator_ui/ && yarn pretest

.PHONY: operator_ui/test
operator_ui/test: ## Execute yarn test for operator_ui subproject
	cd ./operator_ui/ && yarn test

.PHONY: operator_ui/test\:ci
operator_ui/test\:ci: ## Execute yarn test:ci for operator_ui subproject
	cd ./operator_ui/ && yarn test:ci

.PHONY: operator_ui/watch
operator_ui/watch: ## Execute yarn watch for operator_ui subproject
	cd ./operator_ui/ && yarn watch

.PHONY: operator_ui/setup
operator_ui/setup: ## Execute yarn setup for operator_ui subproject
	cd ./operator_ui/ && yarn setup

.PHONY: operator_ui/depcheck
operator_ui/depcheck: ## Execute yarn depcheck for operator_ui subproject
	cd ./operator_ui/ && yarn depcheck

.PHONY: styleguide/start
styleguide/start: ## Execute yarn start for styleguide subproject
	cd ./styleguide/ && yarn start

.PHONY: styleguide/build-storybook
styleguide/build-storybook: ## Execute yarn build-storybook for styleguide subproject
	cd ./styleguide/ && yarn build-storybook

.PHONY: styleguide/build
styleguide/build: ## Execute yarn build for styleguide subproject
	cd ./styleguide/ && yarn build

.PHONY: styleguide/eslint
styleguide/eslint: ## Execute yarn eslint for styleguide subproject
	cd ./styleguide/ && yarn eslint

.PHONY: styleguide/lint
styleguide/lint: ## Execute yarn lint for styleguide subproject
	cd ./styleguide/ && yarn lint

.PHONY: styleguide/format
styleguide/format: ## Execute yarn format for styleguide subproject
	cd ./styleguide/ && yarn format

.PHONY: styleguide/depcheck
styleguide/depcheck: ## Execute yarn depcheck for styleguide subproject
	cd ./styleguide/ && yarn depcheck

.PHONY: styleguide/setup
styleguide/setup: ## Execute yarn setup for styleguide subproject
	cd ./styleguide/ && yarn setup

.PHONY: tools/depcheck
tools/depcheck: ## Execute yarn depcheck for tools subproject
	cd ./tools/ && yarn depcheck

.PHONY: tools/format
tools/format: ## Execute yarn format for tools subproject
	cd ./tools/ && yarn format

.PHONY: tools/lint
tools/lint: ## Execute yarn lint for tools subproject
	cd ./tools/ && yarn lint

.PHONY: tools/setup
tools/setup: ## Execute yarn setup for tools subproject
	cd ./tools/ && yarn setup

# TODO: SGX

help:
	@echo ""
	@echo "         .__           .__       .__  .__        __"
	@echo "    ____ |  |__ _____  |__| ____ |  | |__| ____ |  | __"
	@echo "  _/ ___\|  |  \\\\\\__  \ |  |/    \|  | |  |/    \|  |/ /"
	@echo "  \  \___|   Y  \/ __ \|  |   |  \  |_|  |   |  \    <"
	@echo "   \___  >___|  (____  /__|___|  /____/__|___|  /__|_ \\"
	@echo "       \/     \/     \/        \/             \/     \/"
	@echo ""
	@grep -E '^[a-zA-Z_-/]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
