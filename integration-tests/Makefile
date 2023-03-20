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
	asdf install
endif
# Now install the helm charts that are needed (should be os agnostic)
	helm repo add chainlink-qa https://raw.githubusercontent.com/smartcontractkit/qa-charts/gh-pages/
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update

.PHONY: install_gotestfmt
install_gotestfmt:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail

# Builds the test image
# tag: the tag for the test image being built, example: tag=tate
# base_tag: the tag for the base-test-image to use, example: base_tag=latest
# suite: the test suites to build into the image, example: suite="chaos soak smoke reorg migration performance"
# push: set to true if you want the image pushed or leave blank if not, example: push=true
.PHONY: build_test_image
build_test_image:
	./scripts/buildTestImage $(tag) $(base_tag) "$(suite)" $(push)

## All commands will use 16 threads to run tests in parallel. To change this, use -test.parallel n

# Smoke
.PHONY: test_smoke_raw # Run smoke tests without any gotestfmt or default args
test_smoke_raw:
	go test $(args) ./smoke

.PHONY: test_smoke
test_smoke: install_gotestfmt ## Run all smoke tests
	TEST_LOG_LEVEL="disabled" \
	go test -timeout 24h -count=1 -json $(args) ./smoke 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test_smoke_simulated
test_smoke_simulated: install_gotestfmt ## Run all smoke tests on simulated blockchain
	TEST_LOG_LEVEL="disabled" \
	SELECTED_NETWORKS="SIMULATED,SIMULATED_1,SIMULATED_2" \
	go test -timeout 1h -count=1 -json $(args) ./smoke 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test_smoke_verbose
test_smoke_verbose: ## Run all smoke tests with verbose logging
	go test -timeout 24h -count=1 -v $(args) ./smoke

.PHONY: test_smoke_simulated_verbose
test_smoke_simulated_verbose: ## Run all smoke tests with verbose logging
	SELECTED_NETWORKS="SIMULATED,SIMULATED_1,SIMULATED_2" \
	go test -timeout 24h -count=1 -v $(args) ./smoke

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
test_chaos_pods: install_gotestfmt ## Run all smoke tests
	TEST_LOG_LEVEL="disabled" \
	go test -timeout 2h -count=1 -json $(args) -run 'Test/.*pod-chaos' ./chaos 2>&1 | tee ./gotest.log | gotestfmt

.PHONY: test_chaos_network
test_chaos_network: install_gotestfmt ## Run all smoke tests
	TEST_LOG_LEVEL="disabled" \
	go test -timeout 2h -count=1 -json $(args) -run 'Test/.*network-chaos' ./chaos 2>&1 | tee ./gotest.log | gotestfmt

.PHONY: test_chaos_verbose
test_chaos_verbose: ## Run all smoke tests with verbose logging
	go test -timeout 24h -count=1 -v $(args) ./chaos

# Performance
.PHONY: test_perf
test_perf: test_need_operator_assets ## Run core node performance tests.
	TEST_LOG_LEVEL="disabled" \
	SELECTED_NETWORKS="SIMULATED,SIMULATED_1,SIMULATED_2" \
	go test -timeout 1h -count=1 -json $(args) ./performance 2>&1 | tee /tmp/gotest.log | gotestfmt

# Soak
.PHONY: test_soak_ocr
test_soak_ocr:
	go test -v -count=1 -run TestOCRSoak ./soak

.PHONY: test_soak_ocr_simulated
test_soak_ocr_simulated:
	SELECTED_NETWORKS="SIMULATED" go test -v -count=1 -run TestOCRSoak ./soak

.PHONY: test_soak_forwarder_ocr
test_soak_forwarder_ocr:
	go test -v -count=1 -run TestForwarderOCRSoak ./soak

.PHONY: test_soak_forwarder_ocr_simulated
test_soak_forwarder_ocr_simulated:
	SELECTED_NETWORKS="SIMULATED" go test -v -count=1 -run TestForwarderOCRSoak ./soak

.PHONY: test_soak_automation
test_soak_automation:
	go test -v -run ^TestAutomationBenchmark$$ ./benchmark -count=1

.PHONY: test_soak_automation_simulated
test_soak_automation_simulated:
	SELECTED_NETWORKS="SIMULATED" \
	TEST_INPUTS="TEST_TYPE=SOAK,NUMBEROFCONTRACTS=50,BLOCKRANGE=1000,BLOCKINTERVAL=50,GRAFANA_DASHBOARD_URL=https://chainlinklabs.grafana.net/d/Q8n6m1unz/chainlink-keepers-qa?orgId=1" \
	go test -v -run ^TestAutomationBenchmark$$ ./benchmark -count=1

.PHONY: test_benchmark_automation
test_benchmark_automation: ## Run the automation benchmark tests
	go test -timeout 30m -v -run ^TestAutomationBenchmark$$ ./benchmark -count=1
