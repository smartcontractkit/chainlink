package types

import (
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfTestEnv "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
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
