package testsetups

import (
	"testing"

	e "github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type DonChain struct {
	conf              *DonChainConfig
	EVMClient         blockchain.EVMClient
	EVMNetwork        *blockchain.EVMNetwork
	ContractDeployer  contracts.ContractDeployer
	LinkTokenContract contracts.LinkToken
	ChainlinkNodes    []*client.Chainlink
	Mockserver        *ctfClient.MockserverClient
}

type DonChainConfig struct {
	T               *testing.T
	Env             *e.Environment
	EVMNetwork      *blockchain.EVMNetwork
	EthereumProps   *ethereum.Props
	ChainlinkValues map[string]interface{}
}

func NewDonChain(conf *DonChainConfig) *DonChain {
	return &DonChain{
		conf:       conf,
		EVMNetwork: conf.EVMNetwork,
	}
}

func (s *DonChain) Deploy() {
	var err error

	s.conf.Env.AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(s.conf.EthereumProps)).
		AddHelm(chainlink.New(0, s.conf.ChainlinkValues))

	err = s.conf.Env.Run()
	require.NoError(s.conf.T, err)

	s.initializeClients()
}

func (s *DonChain) initializeClients() {
	var err error
	network := *s.conf.EVMNetwork
	s.EVMClient, err = blockchain.NewEVMClient(network, s.conf.Env)
	require.NoError(s.conf.T, err, "Connecting to blockchain nodes shouldn't fail")

	s.ContractDeployer, err = contracts.NewContractDeployer(s.EVMClient)
	require.NoError(s.conf.T, err)

	s.ChainlinkNodes, err = client.ConnectChainlinkNodes(s.conf.Env)
	require.NoError(s.conf.T, err, "Connecting to chainlink nodes shouldn't fail")

	s.Mockserver, err = ctfClient.ConnectMockServer(s.conf.Env)
	require.NoError(s.conf.T, err, "Creating mockserver clients shouldn't fail")

	s.EVMClient.ParallelTransactions(true)

	s.LinkTokenContract, err = s.ContractDeployer.DeployLinkTokenContract()
	require.NoError(s.conf.T, err, "Deploying Link Token Contract shouldn't fail")
}
