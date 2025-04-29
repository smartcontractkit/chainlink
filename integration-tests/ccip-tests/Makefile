## To Override the default config:
# example usage: make set_config override_toml=../config/config.toml
.PHONY: set_config
set_config:
	if [ -s "$(override_toml)" ]; then \
  			echo "Overriding config with $(override_toml)"; \
			echo "export BASE64_CONFIG_OVERRIDE=$$(base64 -i $(override_toml))" > ./testconfig/override/.env; \
        	echo "export TEST_BASE64_CONFIG_OVERRIDE=$$(base64 -i $(override_toml))" >> ./testconfig/override/.env; \
			echo "BASE64_CONFIG_OVERRIDE=$$(base64 -i $(override_toml))" > ./testconfig/override/debug.env; \
			echo "TEST_BASE64_CONFIG_OVERRIDE=$$(base64 -i $(override_toml))" >> ./testconfig/override/debug.env; \
    else \
      		echo "No override config found, using default config"; \
        	echo  > ./testconfig/override/.env; \
	fi

	@echo "Checking for test secrets file in ~/.testsecrets...";
	@if [ ! -f ~/.testsecrets ]; then \
		echo "WARNING: ~/.testsecrets file not found. No test secrets will be set. To set test secrets, refer to ./testconfig/examples/.testsecrets.example for the list of secrets and instruction how to set them up."; \
	else \
		echo "Test secrets will be loaded from ~/.testsecrets file."; \
	fi	


# example usage: make test_load_ccip testimage=chainlink-ccip-tests:latest testname=TestLoadCCIPStableRPS override_toml=./testconfig/override/config.toml network_config_toml=../config/network.toml
.PHONY: test_load_ccip
test_load_ccip: set_config
	source ./testconfig/override/.env && \
	DATABASE_URL=postgresql://postgres:node@localhost:5432/chainlink_test?sslmode=disable	\
    ENV_JOB_IMAGE=$(testimage)  \
    TEST_SUITE=ccip-load \
    TEST_ARGS="-test.timeout 900h" \
    DETACH_RUNNER=true \
    RR_MEM=45Gi \
    RR_CPU=6 \
	go test -timeout 24h -count=1 -v -run ^$(testname)$$ ./load


# example usage: make test_smoke_ccip testimage=chainlink-ccip-tests:latest testname=TestSmokeCCIPForBidirectionalLane override_toml=../testconfig/override/config.toml network_config_toml=../config/network.toml
.PHONY: test_smoke_ccip
test_smoke_ccip: set_config
	source ./testconfig/override/.env && \
	DATABASE_URL=postgresql://postgres:node@localhost:5432/chainlink_test?sslmode=disable	\
    ENV_JOB_IMAGE=$(testimage)  \
    TEST_SUITE=smoke \
    TEST_ARGS="-test.timeout 900h" \
    DETACH_RUNNER=true \
	go test -timeout 24h -count=1 -v -run ^$(testname)$$ ./smoke

# run ccip smoke tests with default config; explicitly sets the override config to empty
# example usage: make test_smoke_ccip_default testname=TestSmokeCCIPForBidirectionalLane
.PHONY: test_smoke_ccip_default
test_smoke_ccip_default: set_config
	source ./testconfig/override/.env && \
	DATABASE_URL=postgresql://postgres:node@localhost:5432/chainlink_test?sslmode=disable	\
	BASE64_CONFIG_OVERRIDE=""	\
	TEST_BASE64_CONFIG_OVERRIDE=""	\
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

# image: the name for the chainlink image being built, example: image=chainlink
# tag: the tag for the chainlink image being built, example: tag=latest
# example usage: make build_ccip_image image=chainlink-ccip tag=latest
.PHONY: build_ccip_debug_image
build_ccip_debug_image:
	docker build -f ../../core/chainlink.debug.Dockerfile --build-arg COMMIT_SHA=$(git rev-parse HEAD) --build-arg CHAINLINK_USER=chainlink -t $(image):$(tag) ../../
