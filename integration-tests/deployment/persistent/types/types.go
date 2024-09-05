package types

import (
	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctfTestEnv "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logstream"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type ChainConfig struct {
	// ExistingEVMChains are Chains that are already running in a separate process or machine.
	ExistingEVMChains []ExistingEVMChainProducer
	// NewEVMChains are Chains that will be started by the test environment.
	NewEVMChains []NewEVMChainProducer
}

type NewEVMChainProducer interface {
	Chain() (deployment.Chain, RpcProvider, error)
	Hooks() ctfTestEnv.EthereumNetworkHooks
}

type ExistingEVMChainProducer interface {
	Chain() (deployment.Chain, RpcProvider, error)
}

type RpcProvider interface {
	EVMNetwork() blockchain.EVMNetwork
	PrivateHttpUrls() []string
	PrivateWsUrls() []string
	PublicHttpUrls() []string
	PublicWsUrls() []string
}

func NewEVMNetworkWithRPCs(evmNetwork blockchain.EVMNetwork, rpcProvider ctfTestEnv.RpcProvider) RpcProvider {
	return &EVMNetworkWithRPCs{
		evmNetwork,
		rpcProvider,
	}
}

type EVMNetworkWithRPCs struct {
	evmNetwork blockchain.EVMNetwork
	ctfTestEnv.RpcProvider
}

func (s *EVMNetworkWithRPCs) EVMNetwork() blockchain.EVMNetwork {
	return s.evmNetwork
}

func (s *EVMNetworkWithRPCs) PrivateHttpUrls() []string {
	return s.RpcProvider.PrivateHttpUrls()
}

func (s *EVMNetworkWithRPCs) PrivateWsUrls() []string {
	return s.RpcProvider.PrivateWsUrsl()
}

func (s *EVMNetworkWithRPCs) PublicHttpUrls() []string {
	return s.RpcProvider.PublicHttpUrls()
}

func (s *EVMNetworkWithRPCs) PublicWsUrls() []string {
	return s.RpcProvider.PublicWsUrls()
}

type EnvironmentHooks interface {
	// PostChainStartupHooks is called after chains have been created, which means that each chain is already up and there's an on-chain client connected to it.
	PostChainStartupHooks(map[uint64]deployment.Chain, map[uint64]RpcProvider, *EnvironmentConfig) error
	// PostNodeStartupHooks is called after nodes have been created, which means that each node is already up and running, and we have an off-chain client connected to them.
	PostNodeStartupHooks(*DON, *EnvironmentConfig) error
	// PostMocksStartupHooks is called after mocks have been created, which means they are up and running, and we can interact with them.
	PostMocksStartupHooks(*deployment.Mocks, *EnvironmentConfig) error
}

type EnvironmentConfig struct {
	ChainConfig
	DONConfig
	EnvironmentHooks
}

type ExistingDONConfig struct {
	*testconfig.CLCluster
	MockServerURL *string `toml:",omitempty"`
}

type NewDONHooks interface {
	// PreStartupHook is called before the DON is started. No containers are running yet. For example, you can use this hook to modify configuration of each node.
	PreStartupHook([]*test_env.ClNode) error
	// PostStartupHook is called after the DON is started. All containers are running. For example, you can use this hook to interact with them using the API.
	PostStartupHook([]*test_env.ClNode) error
}

type NewDockerDONConfig struct {
	*testconfig.ChainlinkDeployment
	Options          DockerOptions
	ChainlinkConfigs []*chainlink.Config
	NewDONHooks
}

type DockerOptions struct {
	Networks  []string
	LogStream *logstream.LogStream
}

type DONConfig struct {
	ExistingDON *ExistingDONConfig
	NewDON      *NewDockerDONConfig
}

type DON struct {
	ChainlinkClients    []*client.ChainlinkK8sClient
	ChainlinkContainers []*test_env.ClNode
	deployment.Mocks
}
