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

#Build a chainlink docker image for local testing and push to k3d registry
.PHONY: build_push_docker_image
build_push_docker_image: build_docker_image
	docker tag $(image):$(tag) 127.0.0.1:5000/chainlink:develop
	docker push 127.0.0.1:5000/chainlink:develop

# Spins up containers needed to collect traces for local testing
.PHONY: run_tracing
run_tracing:
	cd ../.github/tracing
	docker compose -f ../.github/tracing/local-smoke-docker-compose.yaml up

## All commands will use 16 threads to run tests in parallel. To change this, use -test.parallel n
## Remember to set selected_networks and CL image in the TOML file (e.g. overrides.toml)

.PHONY: check_github_token
check_github_token:
	@if [ -z "$(GITHUB_TOKEN)" ]; then \
		echo "Error: GITHUB_TOKEN environment variable to pull private dependencies. Try to set it via the gh cli: GITHUB_TOKEN=\$$(gh auth token)"; \
		echo "See: https://github.com/smartcontractkit/chainlink?tab=readme-ov-file#build-plugins"; \
		exit 1; \
	fi
# image: the name for the chainlink image being built, example: image=chainlink
# tag: the tag for the chainlink image being built, example: tag=latest
# example usage: make build_docker_image image=chainlink tag=latest
.PHONY: build_docker_image
build_docker_image: check_github_token
	docker build -f ../core/chainlink.Dockerfile \
		--build-arg COMMIT_SHA=$(git rev-parse HEAD) \
		--build-arg VERSION_TAG=$(git describe --always) \
		--build-arg CL_AUTO_DOCKER_TAG=$(tag) \
		--build-arg CHAINLINK_USER=chainlink \
		--secret id=GIT_AUTH_TOKEN,env=GITHUB_TOKEN \
		-t $(image):$(tag) ../

# removes all occurrences of .run.id file in current folder and it's subdirectories
# before making any changes lists all file locations and awaits user confirmation
remove_test_execution_artefacts:
	./scripts/search_and_delete.sh .run.id


# base-image-registry: AWS registry (only secure-sdlc account works)
# check AWS Portal and search for
# <image_registry> | secure-sdlc@smartcontract.com
.PHONY: push-ccip-test-image
push-ccip-test-image:
	@if [ -z "$(base-image-registry)" ]; then echo "Error: base-image-registry variable is not set. Please provide it with 'make push base-image-registry=...'. See README.md explaining where to get it."; exit 1; fi
	aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin $(base-image-registry).dkr.ecr.us-west-2.amazonaws.com
	docker tag chainlink-tests:latest $(base-image-registry).dkr.ecr.us-west-2.amazonaws.com/chainlink-tests:latest
	docker push $(base-image-registry).dkr.ecr.us-west-2.amazonaws.com/chainlink-tests:latest
