include test-config.env
## To Override the default config,
# example usage: make override_config override_toml=../config/config.toml env=../.env
.PHONY: override_config
override_config:
	cd ./testconfig/override && \
	go run . --path=$(override_toml) --output=$(env) && \
	cd ../..


# example usage: make test_load_ccip_simulated_k8 image=chainlink-ccip tag=latest testimage=chainlink-ccip-tests:latest testname=TestLoadCCIPStableRequestTriggeringWithNetworkChaos
.PHONY: test_load_ccip
test_load_ccip:
	source ./testconfig/override/$(env) && \
    CHAINLINK_IMAGE=$(image)  \
    CHAINLINK_VERSION=$(tag)  \
    ENV_JOB_IMAGE=$(testimage)  \
    TEST_SUITE=load \
    TEST_ARGS="-test.timeout 900h" \
    DETACH_RUNNER=true \
	go test -timeout 24h -count=1 -v -run ^$(testname)$$ ./load


# example usage: make test_smoke_ccip image=chainlink-ccip tag=latest testimage=chainlink-ccip-tests:latest testname=TestSmokeCCIPForBidirectionalLane override_toml=../config/config.toml env=test-config.env
# To run the test in besu simulated network, use the following command:
# example usage: make test_smoke_ccip image=chainlink-ccip tag=latest testimage="" testname=TestSmokeCCIPForBidirectionalLane override_toml="" env=test-config.env
.PHONY: test_smoke_ccip
test_smoke_ccip: override_config
	source ./testconfig/override/$(env) && \
    CHAINLINK_IMAGE=$(image)  \
    CHAINLINK_VERSION=$(tag)  \
    ENV_JOB_IMAGE=$(testimage)  \
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
