package types

import (
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfTestEnv "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
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
	return s.PublicHttpUrls()
}

func (s *EVMNetworkWithRPCs) PublicWsUrls() []string {
	return s.RpcProvider.PublicWsUrls()
}

type EnvironmentConfig struct {
	ChainConfig
	DONConfig
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
	Options          Options
	ChainlinkConfigs []*chainlink.Config
	NewDONHooks
}

type Options struct {
	Networks  []string
	LogStream *logstream.LogStream
}

type DONConfig struct {
	ExistingDON *ExistingDONConfig
	NewDON      *NewDockerDONConfig
}

type DON struct {
	ClClients []*client.ChainlinkK8sClient
	deployment.Mocks
}
