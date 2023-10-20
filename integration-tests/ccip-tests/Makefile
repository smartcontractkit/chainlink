# example usage: make test_load_ccip_simulated_k8 image=chainlink-ccip tag=latest testimage=chainlink-ccip-tests:latest
.PHONY: test_load_ccip_simulated_k8
test_load_ccip_simulated_k8:
	source ./load-test.env && \
    CHAINLINK_IMAGE=$(image)  \
    CHAINLINK_VERSION=$(tag)  \
    ENV_JOB_IMAGE=$(testimage)  \
	go test -timeout 24h -count=1 -v -run ^TestLoadCCIPStableRequestTriggeringWithNetworkChaos$$ ./load


# example usage: make test_load_ccip_simulated_k8 image=chainlink-ccip tag=latest testimage=chainlink-ccip-tests:latest
.PHONY: test_load_ccip_testnet_k8
test_load_ccip_testnet_k8:
	source ./load-test.env && \
	SELECTED_NETWORKS="SIMULATED,SEPOLIA,OPTIMISM_GOERLI,ARBITRUM_GOERLI,AVALANCHE_FUJI,BSC_TESTNET,MUMBAI,BASE_GOERLI"  \
	CCIP_NO_OF_NETWORKS=7  \
	CCIP_CHAINLINK_NODE_FUNDING=10  \
    CHAINLINK_IMAGE=$(image)  \
    CHAINLINK_VERSION=$(tag)  \
    ENV_JOB_IMAGE=$(testimage)  \
	go test -timeout 24h -count=1 -v -run ^TestLoadCCIPStableRPS$$ ./load

# example usage: make test_smoke_ccip_simulated_local image=chainlink-ccip tag=latest testname=TestSmokeCCIPForBidirectionalLane
.PHONY: test_smoke_ccip_simulated_local
test_smoke_ccip_simulated_local:
	DATABASE_URL="postgresql://postgres:node@localhost:5432/chainlink_test?sslmode=disable" \
    SELECTED_NETWORKS="SIMULATED,SIMULATED_1,SIMULATED_2"  \
    CCIP_DEPLOY_ON_LOCAL=True  \
    CCIP_NETWORK_PAIRS="" \
    CCIP_NO_OF_NETWORKS=2  \
    CCIP_NO_OF_LANES_PER_PAIR=1  \
    CCIP_MSG_TYPE=WithToken \
    CHAINLINK_IMAGE=$(image)  \
    CHAINLINK_VERSION=$(tag)  \
    CCIP_TESTS_ON_EXISTING_DEPLOYMENT=False  \
    ENV_JOB_IMAGE=""  \
	go test -timeout 24h -count=1 -v -run ^$(testname)$$ ./smoke

# example usage: make test_smoke_ccip_simulated_local_besu image=chainlink tag=latest
.PHONY: test_smoke_ccip_simulated_local_besu
test_smoke_ccip_simulated_local_besu:
	DATABASE_URL="postgresql://postgres:node@localhost:5432/chainlink_test?sslmode=disable" \
    SELECTED_NETWORKS="SIMULATED,SIMULATED_BESU_NONDEV_1,SIMULATED_BESU_NONDEV_2"  \
    CCIP_DEPLOY_ON_LOCAL=True  \
    CCIP_NO_OF_NETWORKS=2  \
    CCIP_MSG_TYPE=WithToken \
    CHAINLINK_IMAGE=$(image)  \
    CHAINLINK_VERSION=$(tag)  \
    CCIP_TESTS_ON_EXISTING_DEPLOYMENT=False  \
    ENV_JOB_IMAGE=""  \
	go test -timeout 24h -count=1 -v -run ^TestSmokeCCIPForBidirectionalLane$$ ./smoke
