BIN_DIR = bin
export GOPATH ?= $(shell go env GOPATH)
export GO111MODULE ?= on

LINUX=LINUX
OSX=OSX
WINDOWS=WIN32
OSFLAG :=
ifeq ($(OS),Windows_NT)
	OSFLAG = $(WINDOWS)
else
	UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
		OSFLAG = $(LINUX)
endif
ifeq ($(UNAME_S),Darwin)
		OSFLAG = $(OSX)
endif
endif

install_qa_tools:
ifeq ($(OSFLAG),$(WINDOWS))
	echo "If you are running windows and know how to install what is needed, please contribute by adding it here!"
	echo "You will need nodejs, golang, k3d, and helm."
	exit 1
else

# linux and mac can use asdf to install all of the dependencies
ifeq ($(shell which asdf), )

# install asdf
ifeq ($(OSFLAG),$(LINUX))
	echo "You will need to install asdf via your linux installer https://asdf-vm.com/guide/getting-started.html"
	exit 1
else
ifeq ($(OSFLAG),$(OSX))
	brew install asdf
endif
endif
endif
# install the plugins if needed and then install the dependencies
	asdf plugin-add nodejs || true
	asdf plugin-add golang || true
	asdf plugin-add k3d || true
	asdf plugin-add helm || true
	asdf plugin-add kubectl || true
	asdf install
endif
# Now install the helm charts that are needed (should be os agnostic)
	helm repo add chainlink-qa https://raw.githubusercontent.com/smartcontractkit/qa-charts/gh-pages/
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update

.PHONY: install_gotestloghelper
install_gotestloghelper:
	go install github.com/smartcontractkit/chainlink-testing-framework/tools/gotestloghelper@latest
	set -euo pipefail

lint:
	golangci-lint --color=always run ./... --fix -v

build:
	@go build ./... go test -run=^# ./...

# Builds the test image
# tag: the tag for the test image being built, example: tag=tate
# base_tag: the tag for the base-test-image to use, example: base_tag=latest
# suite: the test suites to build into the image, example: suite="chaos soak smoke reorg migration"
# push: set to true if you want the image pushed or leave blank if not, example: push=true
.PHONY: build_test_image
build_test_image:
	./scripts/buildTestImage $(tag) $(base_tag) "$(suite)" $(push)

#Build a chainlink docker image for local testing and push to k3d registry
.PHONY: build_push_docker_image
build_push_docker_image:
	docker build -f ../core/chainlink.Dockerfile --build-arg COMMIT_SHA=$(git rev-parse HEAD) --build-arg CHAINLINK_USER=chainlink -t 127.0.0.1:5000/chainlink:develop ../ ; docker push 127.0.0.1:5000/chainlink:develop

#Build a chainlink docker image in plugin mode for local testing and push to k3d registry
.PHONY: build_push_plugin_docker_image
build_push_plugin_docker_image:
	docker build -f ../plugins/chainlink.Dockerfile --build-arg COMMIT_SHA=$(git rev-parse HEAD) --build-arg CHAINLINK_USER=chainlink -t 127.0.0.1:5000/chainlink:develop ../ ; docker push 127.0.0.1:5000/chainlink:develop

# Spins up containers needed to collect traces for local testing
.PHONY: run_tracing
run_tracing:
	cd ../.github/tracing
	docker compose -f ../.github/tracing/local-smoke-docker-compose.yaml up

## Test Runner
.PHONY: run
run:
	go run .

## All commands will use 16 threads to run tests in parallel. To change this, use -test.parallel n
## Remember to set selected_networks and CL image in the TOML file (e.g. overrides.toml)

# Smoke
.PHONY: test_smoke_product
test_smoke_product: ## Run smoke tests for specific product ex: make test_smoke_product product="cron" args="--focus @cron -p"
	ARGS="$(args)" PRODUCT=$(product) ./scripts/run_product_tests

# Chaos
.PHONY: test_chaos_pods_raw
test_chaos_pods_raw:
	go test -timeout 2h -v -count=1 $(args) -p 2 -run 'Test/.*pod-chaos' ./chaos

.PHONY: test_chaos_network_raw
test_chaos_network_raw:
	go test -timeout 2h -v -count=1 $(args) -p 2 -run 'Test/.*network-chaos' ./chaos

.PHONY: test_chaos_pods
test_chaos_pods: install_gotestloghelper ## Run all chaos pod tests
	TEST_LOG_LEVEL="disabled" \
	go test -timeout 2h -count=1 -json $(args) -run 'Test/.*pod-chaos' ./chaos 2>&1 | tee ./gotest.log | gotestloghelper -json -hidepassingtests -tlogprefix -color -singlepackage

.PHONY: test_chaos_network
test_chaos_network: install_gotestloghelper ## Run all chaos network tests
	TEST_LOG_LEVEL="disabled" \
	go test -timeout 2h -count=1 -json $(args) -run 'Test/.*network-chaos' ./chaos 2>&1 | tee ./gotest.log | gotestloghelper -json -hidepassingtests -tlogprefix -color -singlepackage

.PHONY: test_chaos_verbose
test_chaos_verbose: ## Run all smoke tests with verbose logging
	go test -timeout 24h -count=1 -v $(args) ./chaos

# Migrations
.PHONY: test_node_migrations
test_node_migrations: install_gotestloghelper ## Run all node migration tests.
	TEST_LOG_LEVEL="disabled" \
	go test -timeout 1h -count=1 -json $(args) ./migration 2>&1 | tee /tmp/gotest.log | gotestloghelper -json -tlogprefix -color -singlepackage -hidepassinglogs

.PHONY: test_node_migrations_simulated
test_node_migrations_simulated: install_gotestloghelper
	TEST_LOG_LEVEL="disabled" \
	go test -timeout 1h -count=1 -json $(args) ./migration 2>&1 | tee /tmp/gotest.log | gotestloghelper -json -tlogprefix -color -singlepackage -hidepassinglogs

.PHONY: test_node_migrations_verbose
test_node_migrations_verbose:
	go test -timeout 1h -count=1 -v $(args) ./migration

.PHONY: test_node_migrations_simulated_verbose
test_node_migrations_simulated_verbose:
	go test -timeout 1h -count=1 -v $(args) ./migration

# Soak
.PHONY: test_soak_ocr1
test_soak_ocr1:
	. ./scripts/check_base64_env_var.sh
	go test -v -count=1 -run TestOCRv1Soak ./soak

.PHONY: test_soak_ocr2
test_soak_ocr2:
	. ./scripts/check_base64_env_var.sh
	go test -v -count=1 -run TestOCRv2Soak ./soak

.PHONY: test_soak_forwarder_ocr1
test_soak_forwarder_ocr1:
	. ./scripts/check_base64_env_var.sh
	go test -v -count=1 -run TestForwarderOCRv1Soak ./soak

.PHONY: test_soak_forwarder_ocr2
test_soak_forwarder_ocr2:
	. ./scripts/check_base64_env_var.sh
	go test -v -count=1 -run TestForwarderOCRv2Soak ./soak

.PHONY: test_soak_ocr_reorg_1
test_soak_ocr_reorg_1:
	. ./scripts/check_base64_env_var.sh
	go test -v -count=1 -run ^TestOCRSoak_GethReorgBelowFinality_FinalityTagDisabled$$ ./soak

.PHONY: test_soak_ocr_reorg_2
test_soak_ocr_reorg_2:
	. ./scripts/check_base64_env_var.sh
	go test -v -count=1 -run ^TestOCRSoak_GethReorgBelowFinality_FinalityTagEnabled$$ ./soak

.PHONY: test_soak_ocr_gas_spike
test_soak_ocr_gas_spike:
	. ./scripts/check_base64_env_var.sh
	go test -v -count=1 -run ^TestOCRSoak_GasSpike$$ ./soak

.PHONY: test_soak_ocr_gas_limit_change
test_soak_ocr_gas_limit_change:
	. ./scripts/check_base64_env_var.sh
	go test -v -count=1 -run ^TestOCRSoak_ChangeBlockGasLimit$$ ./soak

.PHONY: test_soak_ocr_rpc_down_all_cl_nodes
test_soak_ocr_rpc_down_all_cl_nodes:
	. ./scripts/check_base64_env_var.sh
	go test -v -count=1 -run ^TestOCRSoak_RPCDownForAllCLNodes$$ ./soak

.PHONY: test_soak_ocr_rpc_down_half_cl_nodes
test_soak_ocr_rpc_down_half_cl_nodes:
	. ./scripts/check_base64_env_var.sh
	go test -v -count=1 -run ^TestOCRSoak_RPCDownForHalfCLNodes$$ ./soak


.PHONY: test_soak_automation
test_soak_automation:
	. ./scripts/check_base64_env_var.sh
	go test -v -run ^TestAutomationBenchmark$$ ./benchmark -count=1

.PHONY: test_benchmark_automation
test_benchmark_automation: ## Run the automation benchmark tests
	. ./scripts/check_base64_env_var.sh
	go test -timeout 30m -v -run ^TestAutomationBenchmark$$ ./benchmark -count=1

.PHONY: test_reorg_automation
test_reorg_automation: ## Run the automation reorg tests
	go test -timeout 300m -v -run ^TestAutomationReorg$$ ./reorg -count=1 | tee automation_reorg_run_`date +"%Y%m%d-%H%M%S"`.log

# image: the name for the chainlink image being built, example: image=chainlink
# tag: the tag for the chainlink image being built, example: tag=latest
# example usage: make build_docker_image image=chainlink tag=latest
.PHONY: build_docker_image
build_docker_image:
	docker build -f ../core/chainlink.Dockerfile --build-arg COMMIT_SHA=$(git rev-parse HEAD) --build-arg CHAINLINK_USER=chainlink -t $(image):$(tag) ../

# image: the name for the chainlink image being built, example: image=chainlink
# tag: the tag for the chainlink image being built, example: tag=latest
# example usage: make build_docker_image image=chainlink tag=latest
.PHONY: build_plugin_docker_image
build_plugin_docker_image:
	docker build -f ../plugins/chainlink.Dockerfile --build-arg COMMIT_SHA=$(git rev-parse HEAD) --build-arg CHAINLINK_USER=chainlink -t 127.0.0.1:5000/chainlink:develop ../

# image: the name for the chainlink image being built, example: image=chainlink
# tag: the tag for the chainlink image being built, example: tag=latest
# args: the args to pass to the test runner, example: args="--focus @cron -p"
# product: the product to run tests for, example: product=cron
# example usage: make run_test_with_local_image image=chainlink tag=latest-dev product=cron
# remember to put the case CL image name and tag in the TOML config (and don't forget about selected network configuration)
.PHONY: run_test_with_local_image
run_test_with_local_image: build_docker_image
	ARGS="$(args)" \
	PRODUCT=$(product) 	\
	./scripts/run_product_tests

# removes all occurrences of .run.id file in current folder and it's subdirectories
# before making any changes lists all file locations and awaits user confirmation
remove_test_execution_artefacts:
	./scripts/search_and_delete.sh .run.id
