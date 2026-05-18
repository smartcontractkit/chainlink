# CCIP Configuration

The CCIP configuration is used to specify the test configuration for running the CCIP integration tests.
The configuration is specified in a TOML file. The configuration is used to specify the test environment, test type, test parameters, and other necessary details for running the tests.
The test config is read in following order:

- The test reads the default configuration from [ccip-default.toml](./tomls/ccip-default.toml).
- The default can be overridden by specifying the test config in a separate file.
  - The file content needs to be encoded in base64 format and set in `BASE64_CCIP_CONFIG_OVERRIDE` environment variable.
  - The config mentioned in this file will override the default config.
  - Example override file - [override.toml.example](./examples/override.toml.example)
- If there are sensitive details like private keys, credentials in test config, they can be specified in a separate dotenv file as env vars
  - The `~/.testsecrets` file in home directory is automatically loaded and should have all test secrets as env vars. Learn more about it [here](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/config/README.md#test-secrets) 
  - Example secret file - [.testsecrets.example](./examples/.testsecrets.example)

## CCIP.ContractVersions

Specifies contract versions of different contracts to be referred by test.
Supported versions are:

- **PriceRegistry**: '1.2.0', 'Latest'
- **OffRamp**: '1.2.0', 'Latest'
- **OnRamp**: '1.2.0', 'Latest'
- **TokenPool**: '1.4.0', 'Latest'
- **CommitStore**: '1.2.0', 'Latest'

Example Usage:

```toml
[CCIP.ContractVersions]
PriceRegistry = "1.2.0"
OffRamp = "1.2.0"
OnRamp = "1.2.0"
TokenPool = "1.4.0"
CommitStore = "1.2.0"
```

## CCIP.Deployments

CCIP Deployment contains all necessary contract addresses for various networks. This is mandatory if the test are to be run for [existing deployments](#ccipgroupstestgroupexistingdeployment)
The deployment data can be specified -

- Under `CCIP.Deployments.Data` field with value as stringify format of json.
- Under `CCIP.Deployments.DataFile` field with value as the path of the file containing the deployment data in json format.

The json schema is specified in https://github.com/smartcontractkit/ccip/blob/ccip-develop/integration-tests/ccip-tests/contracts/laneconfig/parse_contracts.go#L96

Example Usage:

```toml
[CCIP.Deployments]
Data = """
{
    "lane_configs": {
        "Arbitrum Mainnet": {
            "is_native_fee_token": true,
            "fee_token": "0xf97f4df75117a78c1A5a0DBb814Af92458539FB4",
            "bridge_tokens": ["0x82aF49447D8a07e3bd95BD0d56f35241523fBab1"],
            "bridge_tokens_pools": ["0x82aF49947D8a07e3bd95BD0d56f35241523fBab1"],
            "arm": "0xe06b0e8c4bd455153e8794ad7Ea8Ff5A14B64E4b",
            "router": "0x141fa059441E0ca23ce184B6A78bafD2A517DdE8",
            "price_registry": "0x13015e4E6f839E1Aa1016DF521ea458ecA20438c",
            "wrapped_native": "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
            "src_contracts": {
                "Ethereum Mainnet": {
                    "on_ramp": "0xCe11020D56e5FDbfE46D9FC3021641FfbBB5AdEE",
                    "deployed_at": 11111111
                }
            },
            "dest_contracts": {
                "Ethereum Mainnet": {
                    "off_ramp": "0x542ba1902044069330e8c5b36A84EC503863722f",
                    "commit_store": "0x060331fEdA35691e54876D957B4F9e3b8Cb47d20",
                    "receiver_dapp": "0x1A2A69e3eB1382FE34Bc579AdD5Bae39e31d4A2c"
                }
            }
        },
        "Ethereum Mainnet": {
            "is_native_fee_token": true,
            "fee_token": "0x514910771AF9Ca656af840dff83E8264EcF986CA",
            "bridge_tokens": ["0x8B63b3DE93431C0f756A493644d128134291fA1b"],
            "bridge_tokens_pools": ["0x8B63b3DE93431C0f756A493644d128134291fA1b"],
            "arm": "0x8B63b3DE93431C0f756A493644d128134291fA1b",
            "router": "0x80226fc0Ee2b096224EeAc085Bb9a8cba1146f7D",
            "price_registry": "0x8c9b2Efb7c64C394119270bfecE7f54763b958Ad",
            "wrapped_native": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
            "src_contracts": {
                "Arbitrum Mainnet": {
                    "on_ramp": "0x925228D7B82d883Dde340A55Fe8e6dA56244A22C",
                    "deployed_at": 11111111
                }
            },
            "dest_contracts": {
                "Arbitrum Mainnet": {
                    "off_ramp": "0xeFC4a18af59398FF23bfe7325F2401aD44286F4d",
                    "commit_store": "0x9B2EEd6A1e16cB50Ed4c876D2dD69468B21b7749",
                    "receiver_dapp": "0x1A2A69e3eB1382FE34Bc579AdD5Bae39e31d4A2c"
                }
            }
        }
    }
}
"""
```

Or,

```toml
[CCIP.Deployments]
DataFile = '<path/to/deployment.json>'
```

## CCIP.Env 

Specifies the environment details for the test to be run on.
Mandatory fields are:

- **Networks**: [CCIP.Env.Networks](#ccipenvnetworks)
- **NewCLCluster**: [CCIP.Env.NewCLCluster](#ccipenvnewclcluster) - This is mandatory if the test needs to deploy Chainlink nodes.
- **ExistingCLCluster**: [CCIP.Env.ExistingCLCluster](#ccipenvexistingclcluster) - This is mandatory if the test needs to run on existing Chainlink nodes to deploy ccip jobs.

Test needs network/chain details to be set through configuration. Set network urls in ~/.testsecrets [see docs](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/config/README.md#test-secrets).

### CCIP.Env.Networks

Specifies the network details for the test to be run.
The NetworkConfig is imported from https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/config/network.go#L39

#### CCIP.Env.Networks.selected_networks

It denotes the network names in which tests will be run. These networks are used to deploy ccip contracts and set up lanes between them.
If more than 2 networks are specified, then lanes will be set up between all possible pairs of networks.

For example , if `selected_networks = ['SIMULATED_1', 'SIMULATED_2', 'SIMULATED_3']`, it denotes that lanes will be set up between SIMULATED_1 and SIMULATED_2, SIMULATED_1 and SIMULATED_3, SIMULATED_2 and SIMULATED_3
This behaviour can be varied based on [NoOfNetworks](#ccipgroupstestgroupnoofnetworks), [NetworkPairs](#ccipgroupstestgroupnetworkpairs), [MaxNoOfLanes](#ccipgroupstestgroupmaxnooflanes) fields in test config.

The name of the networks are taken from [known_networks](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/networks/known_networks.go#L884) in chainlink-testing-framework
If the network is not present in known_networks, then the network details can be specified in the config file itself under the following `EVMNetworks` key.

#### CCIP.Env.Network.EVMNetworks

Specifies the network config to be used while creating blockchain EVMClient for test. 
It is a map of network name to EVMNetworks where key is network name specified under `CCIP.Env.Networks.selected_networks` and value is `EVMNetwork`. 
The EVMNetwork is imported from [EVMNetwork](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/blockchain/config.go#L43) in chainlink-testing-framework.

If `CCIP.Env.Network.EVMNetworks` config is not set for a network name specified under `CCIP.Env.Networks.selected_networks`, test will try to find the corresponding network config from defined networks in `MappedNetworks` under [known_networks.go](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/networks/known_networks.go)

#### CCIP.Env.Network.AnvilConfigs

If the test needs to run on chains created using Anvil, then the AnvilConfigs can be specified.
It is a map of network name to `AnvilConfig` where key is network name specified under `CCIP.Env.Networks.selected_networks` and value is `AnvilConfig`.
The AnvilConfig is imported from [AnvilConfig](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/config/network.go#L20) in chainlink-testing-framework.

**The following network configs are required for tests running on live networks. It can be ignored if the tests are running on simulated networks.**
Refer to [secrets.toml.example](./examples/secrets.toml.example) for details.

#### CCIP.ENV.Network.RpcHttpUrls

RpcHttpUrls is the RPC HTTP endpoints for each network, key is the network name as declared in selected_networks slice.

#### CCIP.ENV.Network.RpcWsUrls

RpcWsUrls is the RPC WS endpoints for each network, key is the network name as declared in selected_networks slice.

#### CCIP.ENV.Network.WalletKeys

WalletKeys is the private keys for each network, key is the network name as declared in selected_networks slice.

Example Usage of Network Config:

```toml
[CCIP.Env.Network]
selected_networks= ['PRIVATE-CHAIN-1', 'PRIVATE-CHAIN-2']

[CCIP.Env.Network.EVMNetworks.PRIVATE-CHAIN-1]
evm_name = 'private-chain-1'
evm_chain_id = 2337
evm_urls = ['wss://ignore-this-url.com']
evm_http_urls = ['https://ignore-this-url.com']
evm_keys = ['59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d']
evm_simulated = true
client_implementation = 'Ethereum'
evm_chainlink_transaction_limit = 5000
evm_transaction_timeout = '3m'
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 1000
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000
evm_finality_depth = 400

[CCIP.Env.Network.EVMNetworks.PRIVATE-CHAIN-2]
evm_name = 'private-chain-2'
evm_chain_id = 1337
evm_urls = ['wss://ignore-this-url.com']
evm_http_urls = ['https://ignore-this-url.com']
evm_keys = ['ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80']
evm_simulated = true
client_implementation = 'Ethereum'
evm_chainlink_transaction_limit = 5000
evm_transaction_timeout = '3m'
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 1000
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000
evm_finality_depth = 400

[CCIP.Env.Network.AnvilConfigs.PRIVATE-CHAIN-1]
block_time = 1

[CCIP.Env.Network.AnvilConfigs.PRIVATE-CHAIN-2]
block_time = 1
```

### CCIP.Env.NewCLCluster

The NewCLCluster config holds the overall deployment configuration for Chainlink nodes.

#### CCIP.Env.NewCLCluster.NoOfNodes

Specifies the number of Chainlink nodes to be deployed.

#### CCIP.Env.NewCLCluster.Common

Specifies the common configuration for all Chainlink nodes if they share the same configuration.

##### Name

Name of the node.

##### NeedsUpgrade

Indicates if the node needs an upgrade during test.

##### ChainlinkImage

Configuration for the Chainlink image.

##### ChainlinkUpgradeImage

Configuration for the Chainlink upgrade image. It is used when the node needs an upgrade.

##### BaseConfigTOML

String containing the base configuration toml content for the Chainlink node config.

##### CommonChainConfigTOML

String containing the common chain configuration toml content for all EVMNodes in chainlink node config.

##### ChainConfigTOMLByChain

String containing the chain-specific configuration toml content for individual EVMNodes in chainlink node config. This is keyed by chain ID.

##### DBImage

Database image for the Chainlink node.

##### DBTag

Database tag/version for the Chainlink node.

#### CCIP.Env.NewCLCluster.Nodes

Specifies the configuration for individual nodes if they differ from the common configuration. The fields are the same as the common configuration.

#### CCIP.Env.NewCLCluster.NodeMemory

Specifies the memory to be allocated for each Chainlink node. This is valid only if the deployment is on Kubernetes.

#### CCIP.Env.NewCLCluster.NodeCPU

Specifies the CPU to be allocated for each Chainlink node. This is valid only if the deployment is on Kubernetes.

#### CCIP.Env.NewCLCluster.DBMemory

Specifies the memory to be allocated for the database. This is valid only if the deployment is on Kubernetes.

#### CCIP.Env.NewCLCluster.DBCPU

Specifies the CPU to be allocated for the database. This is valid only if the deployment is on Kubernetes.

#### CCIP.Env.NewCLCluster.IsStateful

Specifies whether the deployment is StatefulSet on Kubernetes.

#### CCIP.Env.NewCLCluster.DBStorageClass

Specifies the storage class for the database. This is valid only if the deployment is StatefulSet on Kubernetes.

#### CCIP.Env.NewCLCluster.DBCapacity

Specifies the capacity of the database. This is valid only if the deployment is StatefulSet on Kubernetes.

#### CCIP.Env.NewCLCluster.PromPgExporter

Specifies whether to enable Prometheus PostgreSQL exporter. This is valid only if the deployment is on Kubernetes.

#### CCIP.Env.NewCLCluster.DBArgs

Specifies the arguments to be passed to the database. This is valid only if the deployment is on Kubernetes.

Example Usage:

```toml
[CCIP.Env.NewCLCluster]
NoOfNodes = 17
NodeMemory = '12Gi'
NodeCPU = '6'
DBMemory = '10Gi'
DBCPU = '2'
DBStorageClass = 'gp3'
PromPgExporter = true
DBCapacity = '50Gi'
IsStateful = true
DBArgs = ['shared_buffers=2048MB', 'effective_cache_size=4096MB', 'work_mem=64MB']

[CCIP.Env.NewCLCluster.Common]
Name = 'node1'      
DBImage = 'postgres'
DBTag = '13.12'
CommonChainConfigTOML = """
[HeadTracker]
HistoryDepth = 400

[GasEstimator]
PriceMax = '200 gwei'
LimitDefault = 6000000
FeeCapDefault = '200 gwei'
"""
```

### CCIP.Env.ExistingCLCluster

The ExistingCLCluster config holds the overall connection configuration for existing Chainlink nodes.
It is needed when the tests are to be run on Chainlink nodes already deployed on some environment.
If this is specified, test will not need to connect to k8 namespace using [CCIP.Env.EnvToConnect](#ccipenvenvtoconnect).
Test can directly connect to the existing Chainlink nodes using node credentials without knowing the k8 namespace details.

#### CCIP.Env.ExistingCLCluster.Name

Specifies the name of the existing Chainlink cluster. This is used to identify the cluster in the test.

#### CCIP.Env.ExistingCLCluster.NoOfNodes

Specifies the number of Chainlink nodes in the existing cluster.

#### CCIP.Env.ExistingCLCluster.NodeConfigs

Specifies the configuration for individual nodes in the existing cluster. Each node config contains the following fields to connect to the Chainlink node:

##### CCIP.Env.ExistingCLCluster.NodeConfigs.URL

The URL of the Chainlink node.

##### CCIP.Env.ExistingCLCluster.NodeConfigs.Email

The username/email of the Chainlink node credential.

##### CCIP.Env.ExistingCLCluster.NodeConfigs.Password

The password of the Chainlink node credential.

##### CCIP.Env.ExistingCLCluster.NodeConfigs.InternalIP

The internal IP of the Chainlink node.

Example Usage:

```toml
[CCIP.Env.ExistingCLCluster]
Name = 'crib-sample'
NoOfNodes = 5

[[CCIP.Env.ExistingCLCluster.NodeConfigs]]
URL = 'https://crib-sample-demo-node1.main.stage.cldev.sh/'
Email = 'notreal@fakeemail.ch'
Password = 'fj293fbBnlQ!f9vNs'
InternalIP = 'app-node-1'


[[CCIP.Env.ExistingCLCluster.NodeConfigs]]
URL = 'https://crib-sample-demo-node2.main.stage.cldev.sh/'
Email = 'notreal@fakeemail.ch'
Password = 'fj293fbBnlQ!f9vNs'
InternalIP = 'app-node-2'

[[CCIP.Env.ExistingCLCluster.NodeConfigs]]
URL = 'https://crib-sample-demo-node3.main.stage.cldev.sh/'
Email = 'notreal@fakeemail.ch'
Password = 'fj293fbBnlQ!f9vNs'
InternalIP = 'app-node-3'

[[CCIP.Env.ExistingCLCluster.NodeConfigs]]
URL = 'https://crib-ani-demo-node4.main.stage.cldev.sh/'
Email = 'notreal@fakeemail.ch'
Password = 'fj293fbBnlQ!f9vNs'
InternalIP = 'app-node-4'

[[CCIP.Env.ExistingCLCluster.NodeConfigs]]
URL = 'https://crib-sample-demo-node5.main.stage.cldev.sh/'
Email = 'notreal@fakeemail.ch'
Password = 'fj293fbBnlQ!f9vNs'
InternalIP = 'app-node-5'
```

### CCIP.Env.EnvToConnect

This is specified when the test needs to connect to already existing k8s namespace. User needs to have access to the k8 namespace to run the tests through specific kubeconfig file.
Example usage:

```toml
[CCIP.Env]
EnvToConnect="load-ccip-c8972"
```

### CCIP.ENV.TTL

Specifies the time to live for the k8 namespace. This is used to terminate the namespace after the tests are run. This is only valid if the tests are run on k8s.
Example usage:

```toml
[CCIP.Env]
TTL = "11h"
```

### CCIP.Env.Logging

Specifies the logging configuration for the test. Imported from [LoggingConfig](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/config/logging.go#L11) in chainlink-testing-framework.
Example usage:

```toml
[CCIP.Env.Logging]
test_log_collect = false # if set to true will save logs even if test did not fail

[CCIP.Env.Logging.LogStream]
# supported targets: file, loki, in-memory. if empty no logs will be persistet
log_targets = ["file"]
# context timeout for starting log producer and also time-frame for requesting logs
log_producer_timeout = "10s"
# number of retries before log producer gives up and stops listening to logs
log_producer_retry_limit = 10

[CCIP.Env.Logging.Loki]
tenant_id = "..."
endpoint = "https://loki...."

[CCIP.Env.Logging.Grafana]
base_url = "https://grafana..../"
dashboard_url = "/d/6vjVx-1V8/ccip-long-running-tests"
```

### CCIP.Env.Lane.LeaderLaneEnabled

Specifies whether to enable the leader lane feature. This setting is only applicable for new deployments.

## CCIP.Groups

Specifies the test config specific to each test type. Available test types are:

- **CCIP.Groups.load**
- **CCIP.Groups.smoke**
- **CCIP.Groups.chaos**

### CCIP.Groups.[testgroup].KeepEnvAlive

Specifies whether to keep the k8 namespace alive after the test is run. This is only valid if the tests are run on k8s.

### CCIP.Groups.[testgroup].BiDirectionalLane

Specifies whether to set up bi-directional lanes between networks.

### CCIP.Groups.[testgroup].CommitAndExecuteOnSameDON

Specifies whether commit and execution jobs are to be run on the same Chainlink node.

### CCIP.Groups.[testgroup].AllowOutOfOrder

Specifies whether out of order execution is allowed globally for all the chains.

### CCIP.Groups.[testgroup].NoOfCommitNodes

Specifies the number of nodes on which commit jobs are to be run. This needs to be lesser than the total number of nodes mentioned in [CCIP.Env.NewCLCluster.NoOfNodes](#ccipenvnewclclusternoofnodes) or [CCIP.Env.ExistingCLCluster.NoOfNodes](#ccipenvexistingclclusternoofnodes).
If the value of total nodes is `n`, then the max value of NoOfCommitNodes should be less than `n-1`. As the first node is used for bootstrap job.
If the NoOfCommitNodes is lesser than `n-1`, then the remaining nodes are used for execution jobs if `CCIP.Groups.[testgroup].CommitAndExecuteOnSameDON` is set to false.

### CCIP.Groups.[testgroup].TokenConfig

Specifies the token configuration for the test. The token configuration is used to set up tokens and token pools for all chains.

#### CCIP.Groups.[testgroup].TokenConfig.NoOfTokensPerChain

Specifies the number of tokens to be set up for each chain.

#### CCIP.Groups.[testgroup].TokenConfig.WithPipeline

Specifies whether to set up token pipelines in commit jobspec. If set to false, the token prices will be set with DynamicPriceGetterConfig.

#### CCIP.Groups.[testgroup].TokenConfig.TimeoutForPriceUpdate

Specifies the timeout to wait for token and gas price updates to be available in price registry for each chain.

#### CCIP.Groups.[testgroup].TokenConfig.NoOfTokensWithDynamicPrice

Specifies the number of tokens to be set up with dynamic price update. The rest of the tokens will be set up with static price. This is only valid if [WithPipeline](#ccipgroupstestgrouptokenconfigwithpipeline) is set to false.

#### CCIP.Groups.[testgroup].TokenConfig.DynamicPriceUpdateInterval

Specifies the interval for dynamic price update for tokens. This is only valid if [NoOfTokensWithDynamicPrice](#ccipgroupstestgrouptokenconfignooftokenswithdynamicprice) is set to value greater tha zero.

#### CCIP.Groups.[testgroup].TokenConfig.CCIPOwnerTokens

Specifies the tokens to be owned by the CCIP owner. If this is false, the tokens and pools will be owned by an address other than rest of CCIP contract admin addresses.
This is applicable only if the contract versions are '1.5' or higher.

Example Usage:

```toml
[CCIP.Groups.load.TokenConfig]
TimeoutForPriceUpdate = '15m'
NoOfTokensPerChain = 60
NoOfTokensWithDynamicPrice = 15
DynamicPriceUpdateInterval ='15s'
CCIPOwnerTokens = true
```

### CCIP.Groups.[testgroup].MsgDetails

Specifies the ccip message details to be sent by the test.

#### CCIP.Groups.[testgroup].MsgDetails.MsgType

Specifies the type of message to be sent. The supported message types are:

- **Token**
- **Data**
- **DataWithToken**

#### CCIP.Groups.[testgroup].MsgDetails.DestGasLimit

Specifies the gas limit for the destination chain. This is used to in `ExtraArgs` field of CCIPMessage. Change this to 0, if you are doing ccip-send to an EOA in the destination chain.

#### CCIP.Groups.[testgroup].MsgDetails.DataLength

Specifies the length of data to be sent in the message. This is only valid if [MsgType](#ccipgroupstestgroupmsgdetailsmsgtype) is set to 'Data' or 'DataWithToken'.

#### CCIP.Groups.[testgroup].MsgDetails.NoOfTokens

Specifies the number of tokens to be sent in the message. This is only valid if [MsgType](#ccipgroupstestgroupmsgdetailsmsgtype) is set to 'Token' or 'DataWithToken'.
It needs to be less than or equal to [NoOfTokensPerChain](#ccipgroupstestgrouptokenconfignooftokensperchain) specified in the test config.

#### CCIP.Groups.[testgroup].MsgDetails.TokenAmount

Specifies the amount for each token to be sent in the message. This is only valid if [MsgType](#ccipgroupstestgroupmsgdetailsmsgtype) is set to 'Token' or 'DataWithToken'.

Example Usage:

```toml
[CCIP.Groups.smoke.MsgDetails]
MsgType = 'DataWithToken' 
DestGasLimit = 100000
DataLength = 1000
NoOfTokens = 2
AmountPerToken = 1
```

### CCIP.Groups.[testgroup].MulticallInOneTx

Specifies whether to send multiple ccip messages in a single transaction.

### CCIP.Groups.[testgroup].NoOfSendsInMulticall

Specifies the number of ccip messages to be sent in a single transaction. This is only valid if [MulticallInOneTx](#ccipgroupstestgroupmulticallinonetx) is set to true.

### CCIP.Groups.[testgroup].PhaseTimeout

The test validates various events in a ccip request lifecycle, like commit, execute, etc. This field specifies the timeout for each phase in the lifecycle.
The timeout is calculated from the time the previous phase event is received.
The following contract events are validated:

- **CCIPSendRequested on OnRamp**
- **CCIPSendRequested event log to be Finalized**
- **ReportAccepted on CommitStore**
- **TaggedRootBlessed on ARM/RMN**
- **ExecutionStateChanged on OffRamp**

### CCIP.Groups.[testgroup].LocalCluster

Specifies whether the test is to be run on a local docker. If set to true, the test environment will be set up on a local docker.

### CCIP.Groups.[testgroup].ExistingDeployment

Specifies whether the test is to be run on existing deployments. If set to true, the test will use the deployment data specified in [CCIP.Deployments](#ccipdeployments) for interacting with the ccip contracts.
If the deployment data does not contain the required contract addresses, the test will fail.

### CCIP.Groups.[testgroup].ReuseContracts

Test loads contract/lane config from [contracts.json](../contracts/laneconfig/contracts.json) if no lane config is specified in [CCIP.Deployments](#ccipdeployments)
If a certain contract is present in the contracts.json, the test will use the contract address from the contracts.json.
This field specifies whether to reuse the contracts from [contracts.json](../contracts/laneconfig/contracts.json)
For example if the contracts.json contains the contract address for PriceRegistry for `Arbitrum Mainnet`, the test by default will use the contract address from contracts.json instead of redeploying the contract.
If `ReuseContracts` is set to false, the test will redeploy the contract instead of using the contract address from contracts.json.

### CCIP.Groups.[testgroup].NodeFunding

Specified the native token funding for each Chainlink node. It assumes that the native token decimals is 18.
The funding is done by the private key specified in [CCIP.Env.Networks](#ccipenvnetworks) for each network.
The funding is done only if the test is run on local docker or k8s. This is not applicable for [existing deployments](#ccipgroupstestgroupexistingdeployment) is set to true.

### CCIP.Groups.[testgroup].NetworkPairs

Specifies the network pairs for which the test is to be run. The test will set up lanes only between the specified network pairs.
If the network pairs are not specified, the test will set up lanes between all possible pairs of networks mentioned in selected_networks in [CCIP.Env.Networks](#ccipenvnetworksselectednetworks)

### CCIP.Groups.[testgroup].NoOfNetworks

Specifies the number of networks to be used for the test.
If the number of networks is greater than the total number of networks specified in [CCIP.Env.Networks.selected_networks](#ccipenvnetworksselectednetworks):

- the test will fail if the networks are live networks.
- the test will create equal number of replicas of the first network with a new chain id if the networks are simulated networks. 
  For example, if the `selected_networks` is ['SIMULATED_1','SIMULATED_2'] and `NoOfNetworks` is 3, the test will create 1 more network config by copying the network config of `SIMULATED_1` with a different chain id and use that as 3rd network.

### CCIP.Groups.[testgroup].NoOfRoutersPerPair

Specifies the number of routers to be set up for each network.

### CCIP.Groups.[testgroup].MaxNoOfLanes

Specifies the maximum number of lanes to be set up between networks. If this values is not set, the test will set up lanes between all possible pairs of networks mentioned in `selected_networks` in [CCIP.Env.Networks](#ccipenvnetworksselectednetworks).
For example, if `selected_networks = ['SIMULATED_1', 'SIMULATED_2', 'SIMULATED_3']`, and `MaxNoOfLanes` is set to 2, it denotes that the test will randomly select the 2 lanes between all possible pairs `SIMULATED_1`, `SIMULATED_2`, and `SIMULATED_3` for the test run.

### CCIP.Groups.[testgroup].DenselyConnectedNetworkChainIds

This is applicable only if [MaxNoOfLanes](#ccipgroupstestgroupmaxnooflanes) is specified.
Specifies the chain ids for networks to be densely connected. If this is provided the test will include all possible pairs of networks mentioned in `DenselyConnectedNetworkChainIds`.
The rest of the networks will be connected randomly based on the value of `MaxNoOfLanes`.

### CCIP.Groups.[testgroup].ChaosDuration

Specifies the duration for which the chaos experiment is to be run. This is only valid if the test type is 'chaos'.

### CCIP.Groups.[testgroup].USDCMockDeployment

Specifies whether to deploy USDC mock contract for the test. This is only valid if the test is not run on [existing deployments](#ccipgroupstestgroupexistingdeployment).

The following fields are used for various parameters in OCR2 commit and execution jobspecs. All of these are only valid if the test is not run on [existing deployments](#ccipgroupstestgroupexistingdeployment).

### CCIP.Groups.[testgroup].CommitOCRParams

Specifies the OCR parameters for the commit job. This is only valid if the test is not run on [existing deployments](#ccipgroupstestgroupexistingdeployment).

### CCIP.Groups.[testgroup].ExecuteOCRParams

Specifies the OCR parameters for the execute job. This is only valid if the test is not run on [existing deployments](#ccipgroupstestgroupexistingdeployment).

### CCIP.Groups.[testgroup].CommitInflightExpiry

Specifies the value for the `InflightExpiry` in commit job's offchain config. This is only valid if the test is not run on [existing deployments](#ccipgroupstestgroupexistingdeployment).

### CCIP.Groups.[testgroup].SkipRequestIfAnotherRequestTriggeredWithin

If there is CCIP Send requested event present within this duration, the test will skip sending another 
request during load run or avoid sending request in smoke test in that lane. For Example,
if `SkipRequestIfAnotherRequestTriggeredWithin` is set to `40m`, and a request is triggered at 0th second, the test will skip sending another request for another 40m.
This particular field is used to avoid sending transaction when there is traffic already in that lane.

### CCIP.Groups.[testgroup].OffRampConfig

Specifies the offramp configuration for the execution job. This is only valid if the test is not run on [existing deployments](#ccipgroupstestgroupexistingdeployment).
This is used to set values for following fields in execution jobspec's offchain and onchain config:

- **OffRampConfig.MaxDataBytes**
- **OffRampConfig.BatchGasLimit**
- **OffRampConfig.InflightExpiry**
- **OffRampConfig.RootSnooze**

Example Usage:

```toml
[CCIP.Groups.load]
CommitInflightExpiry = '5m'

[CCIP.Groups.load.CommitOCRParams]
DeltaProgress = '2m'
DeltaResend = '5s'
DeltaRound = '75s'
DeltaGrace = '5s'
MaxDurationQuery = '100ms'
MaxDurationObservation = '35s'
MaxDurationReport = '10s'
MaxDurationShouldAcceptFinalizedReport = '5s'
MaxDurationShouldTransmitAcceptedReport = '10s'

[CCIP.Groups.load.ExecOCRParams]
DeltaProgress = '2m'
DeltaResend = '5s'
DeltaRound = '75s'
DeltaGrace = '5s'
MaxDurationQuery = '100ms'
MaxDurationObservation = '35s'
MaxDurationReport = '10s'
MaxDurationShouldAcceptFinalizedReport = '5s'
MaxDurationShouldTransmitAcceptedReport = '10s'

[CCIP.Groups.load.OffRampConfig]
BatchGasLimit = 11000000
MaxDataBytes = 1000
InflightExpiry = '5m'
RootSnooze = '5m'
```

### CCIP.Groups.[testgroup].StoreLaneConfig

This is only valid if the tests are run on remote runners in k8s. If set to true, the test will store the lane config in the remote runner.

### CCIP.Groups.[testgroup].LoadProfile

Specifies the load profile for the test. Only valid if the testgroup is 'load'.

### CCIP.Groups.[testgroup].LoadProfile.LoadFrequency.[DestNetworkName]

#### CCIP.Groups.[testgroup].LoadProfile.RequestPerUnitTime

Specifies the number of requests to be sent per unit time. This is applicable to all networks if [LoadFrequency](#ccipgroupstestgrouploadprofileloadfrequencydestnetworkname) is not specified for a destination network.

#### CCIP.Groups.[testgroup].LoadProfile.TimeUnit

Specifies the unit of time for the load profile. This is applicable to all networks if [LoadFrequency](#ccipgroupstestgrouploadprofileloadfrequencydestnetworkname) is not specified for a destination network.

#### CCIP.Groups.[testgroup].LoadProfile.StepDuration

Specifies the duration for each step in the load profile. This is applicable to all networks if [LoadFrequency](#ccipgroupstestgrouploadprofileloadfrequencydestnetworkname) is not specified for a destination network.

#### CCIP.Groups.[testgroup].LoadProfile.TestDuration

Specifies the total duration for the load test.

#### CCIP.Groups.[testgroup].LoadProfile.NetworkChaosDelay

Specifies the duration network delay used for `NetworkChaos` experiment. This is only valid if the test is run on k8s and not on [existing deployments](#ccipgroupstestgroupexistingdeployment).

#### CCIP.Groups.[testgroup].LoadProfile.WaitBetweenChaosDuringLoad

If there are multiple chaos experiments, this specifies the duration to wait between each chaos experiment. This is only valid if the test is run on k8s and not on [existing deployments](#ccipgroupstestgroupexistingdeployment).

#### CCIP.Groups.[testgroup].LoadProfile.OptimizeSpace

This is used internally to optimize memory usage during load run. If set to true, after the initial lane set up is over the test will discard the lane config to save memory.
The test will only store contract addresses strictly necessary to trigger/validate ccip-send requests.
Except for following contracts all other contract addresses will be discarded after the initial lane set up -

- Router
- ARM
- CommitStore
- OffRamp
- OnRamp

#### CCIP.Groups.[testgroup].LoadProfile.FailOnFirstErrorInLoad

If set to true, the test will fail on the first error encountered during load run. If set to false, the test will continue to run even if there are errors during load run.

#### CCIP.Groups.[testgroup].LoadProfile.SendMaxDataInEveryMsgCount

Specifies the number of requests to send with maximum data in every mentioned count iteration.
For example, if `SendMaxDataInEveryMsgCount` is set to 5, the test will send ccip message with max allowable data length(as set in onRamp config) in every 5th request.

#### CCIP.Groups.[testgroup].LoadProfile.TestRunName

Specifies the name of the test run. This is used to identify the test run in CCIP test dashboard or logs. If multiple tests are run with same `TestRunName`, the test results will be aggregated under the same test run in grafana dashboard.
This is used when multiple iterations of tests are run against same release version to aggregate the results under same dashboard view.

#### CCIP.Groups.[testgroup].LoadProfile.MsgProfile

Specifies the message profile for the test. The message profile is used to set up multiple ccip message details during load test.

##### CCIP.Groups.[testgroup].LoadProfile.MsgProfile.Frequencies

Specifies the frequency of each message profile.
For example, if `Frequencies` is set to [1, 2, 3], the test will send 1st message profile 1 time, 2nd message profile 2 times, and 3rd message profile 3 times in each iteration. Each iteration will be defined by (1+2+3) = 6 requests.
Example Breakdown:

- Frequencies = [4, 12, 3, 1]
- Total Sum of Frequencies = 4 + 12 + 3 + 1 = 20
- Percentages:
  - Message Type 1: (4 / 20) * 100% = 20%
  - Message Type 2: (12 / 20) * 100% = 60%
  - Message Type 3: (3 / 20) * 100% = 15%
  - Message Type 4: (1 / 20) * 100% = 5%
 These percentages reflect how often each message type should appear in the total set of messages.
 Please note - if the total set of messages is not equal to the multiple of sum of frequencies, the percentages will not be accurate.

##### CCIP.Groups.[testgroup].LoadProfile.MsgProfile.MsgDetails

Specifies the message details for each message profile. The fields are the same as [CCIP.Groups.[testgroup].MsgDetails](#ccipgroupstestgroupmsgdetails).

example usage:

```toml
# to represent 20%, 60%, 15%, 5% of the total messages
[CCIP.Groups.load.LoadProfile.MsgProfile]
Frequencies = [4,12,3,1]

[[CCIP.Groups.load.LoadProfile.MsgProfile.MsgDetails]]
MsgType = 'Token'
DestGasLimit = 0
DataLength = 0
NoOfTokens = 5
AmountPerToken = 1

[[CCIP.Groups.load.LoadProfile.MsgProfile.MsgDetails]]
MsgType = 'DataWithToken'
DestGasLimit = 500000
DataLength = 5000
NoOfTokens = 5
AmountPerToken = 1

[[CCIP.Groups.load.LoadProfile.MsgProfile.MsgDetails]]
MsgType = 'Data'
DestGasLimit = 800000
DataLength = 10000

[[CCIP.Groups.load.LoadProfile.MsgProfile.MsgDetails]]
MsgType = 'Data'
DestGasLimit = 2500000
DataLength = 10000
```
