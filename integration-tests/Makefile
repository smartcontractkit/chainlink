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

## All commands will use 16 threads to run tests in parallel. To change this, use -test.parallel n
.PHONY: test_smoke
test_smoke: ## Run all smoke tests
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	TEST_LOG_LEVEL="disabled" \
	go test -timeout 24h -json $(args) ./smoke 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test_smoke_simulated
test_smoke_simulated: ## Run all smoke tests on simulated blockchain
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	TEST_LOG_LEVEL="disabled" \
	SELECTED_NETWORKS="SIMULATED,SIMULATED_1,SIMULATED_2" \
	go test -timeout 1h -json $(args) ./smoke 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test_smoke_verbose
test_smoke_verbose: ## Run all smoke tests with verbose logging
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	go test -timeout 24h -v -json $(args) ./smoke 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test_smoke_product
test_smoke_product: ## Run smoke tests for specific product ex: make test_smoke_product product="cron" args="--focus @cron -p"
	ARGS="$(args)" PRODUCT=$(product) ./scripts/run_product_tests
