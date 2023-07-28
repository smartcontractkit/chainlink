package docker

import (
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
	These are the high level components you should reuse in your docker tests in other repos
	Interfaces from env.go should be moved to chainlink-env or CTF in the next stage
	Implementations must stay in the repos where they belong
*/

type Endpoints struct {
	Networks   []string
	Nodes      []string
	Mockserver string
}

type Clients struct {
	Networks         []blockchain.EVMClient
	NetworkDeployers []contracts.ContractDeployer
	Mockserver       *ctfClient.MockserverClient
	Chainlink        []*client.Chainlink
}

// ConnectClients connects network/CL nodes/mockserver clients
func ConnectClients(network blockchain.EVMNetwork, env *Environment) (*Clients, error) {
	endpoints := &Endpoints{
		Networks: make([]string, 0),
		Nodes:    make([]string, 0),
	}
	clients := &Clients{
		Chainlink: make([]*client.Chainlink, 0),
		Networks:  make([]blockchain.EVMClient, 0),
	}

	network.URLs = make([]string, 0)
	network.HTTPURLs = make([]string, 0)

	// networks
	for _, networkNode := range env.Get("geth") {
		url := networkNode.(*Geth).ExternalWsUrl
		network.Name = "geth"
		network.URLs = append(network.URLs, url)
		network.HTTPURLs = append(network.URLs, url)
		c, err := blockchain.NewDecoupledEVMClient(network)
		if err != nil {
			return nil, err
		}
		cd, err := contracts.NewContractDeployer(c)
		if err != nil {
			return nil, err
		}
		endpoints.Networks = append(endpoints.Networks, url)
		clients.Networks = append(clients.Networks, c)
		clients.NetworkDeployers = append(clients.NetworkDeployers, cd)
	}

	// cl nodes
	for _, n := range env.Get("chainlink") {
		endpoints.Nodes = append(endpoints.Nodes, n.(*Chainlink).Endpoint)
		clc, err := client.NewChainlink(&client.ChainlinkConfig{
			URL:      n.(*Chainlink).Endpoint,
			Email:    "local@local.com",
			Password: "localdevpassword",
		})
		if err != nil {
			return nil, err
		}
		clients.Chainlink = append(clients.Chainlink, clc)
	}

	// mockserver
	msComponent := env.Get("mockserver")[0].(*MockServer)
	endpoints.Mockserver = msComponent.Endpoint
	ms := ctfClient.NewMockserverClient(&ctfClient.MockserverConfig{
		LocalURL:   endpoints.Mockserver,
		ClusterURL: msComponent.InternalEndpoint,
	})
	clients.Mockserver = ms
	log.Info().Interface("Endpoints", endpoints).Msg("Connected to environment")
	return clients, nil
}

// NewCommonChainlinkCluster is a common N nodes cluster that contains:
// - N nodes of Chainlink
// - Geth dev network
// - Mockserver
func NewCommonChainlinkCluster(t *testing.T, nodes int) (*Environment, error) {
	// if pushToLoki set to true you can select logs for your test like
	// {type="log_watch", test="TestFluxBasic", container="..."}
	lw, err := logwatch.NewLogWatch(t, nil)
	require.NoError(t, err)
	env, err := NewEnvironment(lw, false).
		Add(NewGeth(nil)).
		Add(NewMockServer(nil)).
		Start(true)
	require.NoError(t, err)
	gethComponent := env.Get("geth")[0].(*Geth)
	for i := 0; i < nodes; i++ {
		env.Add(NewChainlink(NodeConfigOpts{
			EVM: NodeEVMSettings{
				HTTPURL: gethComponent.InternalHttpUrl,
				WSURL:   gethComponent.InternalWsUrl,
			}}))
	}
	env, err = env.Start(true)
	require.NoError(t, err)
	return env, nil
}
