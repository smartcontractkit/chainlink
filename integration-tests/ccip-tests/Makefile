## To Override the default config,
# example usage: make override_config override_toml=../config/config.toml env=../.env
.PHONY: override_config
override_config:
	cd ./testconfig/override && \
	go run . --path=$(override_toml) --output=$(env) && \
	cd ../..


# example usage: make test_load_ccip image=chainlink-ccip tag=latest testimage=chainlink-ccip-tests:latest testname=TestLoadCCIPStableRequestTriggeringWithNetworkChaos
.PHONY: test_load_ccip
include test-config.env
test_load_ccip: override_config
	source ./testconfig/override/$(env) && \
    CHAINLINK_IMAGE=$(image)  \
    CHAINLINK_VERSION=$(tag)  \
    ENV_JOB_IMAGE=$(testimage)  \
    TEST_SUITE=load \
    TEST_ARGS="-test.timeout 900h" \
    DETACH_RUNNER=true \
    RR_MEM=16Gi \
    RR_CPU=4 \
	go test -timeout 24h -count=1 -v -run ^$(testname)$$ ./load


# example usage: make test_smoke_ccip image=chainlink-ccip tag=latest testimage=chainlink-ccip-tests:latest testname=TestSmokeCCIPForBidirectionalLane override_toml=../config/config.toml env=test-config.env
.PHONY: test_smoke_ccip
include test-config.env
test_smoke_ccip: override_config
	source ./testconfig/override/$(env) && \
    CHAINLINK_IMAGE=$(image)  \
    CHAINLINK_VERSION=$(tag)  \
    ENV_JOB_IMAGE=$(testimage)  \
    TEST_SUITE=smoke \
    TEST_ARGS="-test.timeout 900h" \
    DETACH_RUNNER=true \
	go test -timeout 24h -count=1 -v -run ^$(testname)$$ ./smoke

# run ccip smoke tests with default config; explicitly sets the override config to empty
# example usage: make test_smoke_ccip_default image=chainlink-ccip tag=latest testname=TestSmokeCCIPForBidirectionalLane
.PHONY: test_smoke_ccip_default
include test-config.env
test_smoke_ccip_default:
	BASE64_TEST_CONFIG_OVERRIDE=""	\
	TEST_BASE64_TEST_CONFIG_OVERRIDE=""	\
	CHAINLINK_IMAGE=$(image)  \
	CHAINLINK_VERSION=$(tag)  \
	ENV_JOB_IMAGE=""  \
	TEST_SUITE=smoke \
	TEST_ARGS="-test.timeout 900h" \
	DETACH_RUNNER=true \
	go test -timeout 24h -count=1 -v -run ^$(testname)$$ ./smoke


# image: the name for the chainlink image being built, example: image=chainlink
# tag: the tag for the chainlink image being built, example: tag=latest
# example usage: make build_ccip_image image=chainlink-ccip tag=latest
.PHONY: build_ccip_image
build_ccip_image:
	docker build -f ../../core/chainlink.Dockerfile --build-arg COMMIT_SHA=$(git rev-parse HEAD) --build-arg CHAINLINK_USER=chainlink -t $(image):$(tag) ../../
