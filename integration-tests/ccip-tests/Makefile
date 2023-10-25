## To Override the default config,
# example usage: make override_config overrideconfig="" override_toml=../config/config.toml env=../.env
.PHONY: override_config
override_config:
	cd ./testconfig/override && \
	go run . --overridewith=$(overrideconfig) --path=$(override_toml) --output=$(env) && \
	cd ../..

# example usage: make test_load_ccip_simulated_k8 image=chainlink-ccip tag=latest testimage=chainlink-ccip-tests:latest testname=TestLoadCCIPStableRequestTriggeringWithNetworkChaos
.PHONY: test_load_ccip_simulated_k8
test_load_ccip_simulated_k8:override_config
	source ./testconfig/override/$(env) && \
    CHAINLINK_IMAGE=$(image)  \
    CHAINLINK_VERSION=$(tag)  \
    ENV_JOB_IMAGE=$(testimage)  \
	go test -timeout 24h -count=1 -v -run ^$(testname)$$ ./load


# example usage: make test_load_ccip_simulated_k8 image=chainlink-ccip tag=latest testimage=chainlink-ccip-tests:latest
.PHONY: test_load_ccip_testnet_k8
include test-config.env
overrideconfig := "CCIP.Env.Networks = ['SEPOLIA', 'OPTIMISM_GOERLI', 'ARBITRUM_GOERLI', 'AVALANCHE_FUJI', 'BSC_TESTNET', 'MUMBAI', 'BASE_GOERLI]"
override overrideconfig += "CCIP.Groups.load.NodeFunding = 10"
test_load_ccip_testnet_k8: override_config
	source ./testconfig/override/$(env) && \
    CHAINLINK_IMAGE=$(image)  \
    CHAINLINK_VERSION=$(tag)  \
    ENV_JOB_IMAGE=$(testimage)  \
	go test -timeout 24h -count=1 -v -run ^TestLoadCCIPStableRPS$$ ./load

# example usage: make test_smoke_ccip_simulated_local image=chainlink-ccip tag=latest testname=TestSmokeCCIPForBidirectionalLane override_toml=../config/config.toml env=test-config.env
.PHONY: test_smoke_ccip_simulated_local
include test-config.env
test_smoke_ccip_simulated_local: override_config
	source ./testconfig/override/$(env) && \
    CHAINLINK_IMAGE=$(image)  \
    CHAINLINK_VERSION=$(tag)  \
    ENV_JOB_IMAGE=""  \
	go test -timeout 24h -count=1 -v -run ^$(testname)$$ ./smoke

# example usage: make test_smoke_ccip_simulated_local_besu image=chainlink-ccip tag=latest override_toml=../config/config.toml env=test-config.env
.PHONY: test_smoke_ccip_simulated_local_besu
include test-config.env
overrideconfig := "CCIP.Env.Networks = ['SIMULATED_BESU_NONDEV_1', 'SIMULATED_BESU_NONDEV_2']"
test_smoke_ccip_simulated_local_besu: override_config
	source ./testconfig/override/$(env) && \
  	CHAINLINK_IMAGE=$(image)  \
  	CHAINLINK_VERSION=$(tag)  \
    ENV_JOB_IMAGE=""  \
	go test -timeout 24h -count=1 -v -run ^TestSmokeCCIPForBidirectionalLane$$ ./smoke
