package crib

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	seth_utils "github.com/smartcontractkit/chainlink-testing-framework/utils/seth"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	msClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
)

const (
	// these are constants for simulated CRIB that should never change
	// CRIB: https://github.com/smartcontractkit/crib/tree/main/core
	// Core Chart: https://github.com/smartcontractkit/infra-charts/tree/main/chainlink-cluster
	mockserverCRIBTemplate        = "https://%s-mockserver%s"
	internalNodeDNSTemplate       = "app-node%d"
	ingressNetworkWSURLTemplate   = "wss://%s-geth-1337-ws%s"
	ingressNetworkHTTPURLTemplate = "https://%s-geth-1337-http%s"
)

func setSethConfig(cfg tc.TestConfig, netWSURL string, netHTTPURL string) {
	netName := "CRIB_SIMULATED"
	cfg.Network.SelectedNetworks = []string{netName}
	cfg.Network.RpcHttpUrls = map[string][]string{}
	cfg.Network.RpcHttpUrls[netName] = []string{netHTTPURL}
	cfg.Network.RpcWsUrls = map[string][]string{}
	cfg.Network.RpcWsUrls[netName] = []string{netWSURL}
	cfg.Seth.EphemeralAddrs = ptr.Ptr(int64(0))
}

type ConnectionVars struct {
	IngressSuffix string
	Namespace     string
	Network       string
	Nodes         int
}

func ReadCRIBVars() (*ConnectionVars, error) {
	ingressSuffix := os.Getenv("K8S_STAGING_INGRESS_SUFFIX")
	if ingressSuffix == "" {
		return nil, errors.New("K8S_STAGING_INGRESS_SUFFIX must be set to connect to k8s ingresses")
	}
	cribNamespace := os.Getenv("CRIB_NAMESPACE")
	if cribNamespace == "" {
		return nil, errors.New("CRIB_NAMESPACE must be set to connect")
	}
	cribNetwork := os.Getenv("CRIB_NETWORK")
	if cribNetwork == "" {
		return nil, errors.New("CRIB_NETWORK must be set to connect, only 'geth' is supported for now")
	}
	cribNodes := os.Getenv("CRIB_NODES")
	nodes, err := strconv.Atoi(cribNodes)
	if err != nil {
		return nil, errors.New("CRIB_NODES must be a number, 5-19 nodes")
	}
	if nodes < 2 {
		return nil, fmt.Errorf("not enough chainlink nodes, need at least 2")
	}
	return &ConnectionVars{
		IngressSuffix: ingressSuffix,
		Namespace:     cribNamespace,
		Network:       cribNetwork,
		Nodes:         nodes,
	}, nil
}

// ConnectRemote connects to a local environment, see https://github.com/smartcontractkit/crib/tree/main/core
// connects to default CRIB network if simulated = true
func ConnectRemote() (
	*seth.Client,
	*msClient.MockserverClient,
	*client.ChainlinkK8sClient,
	[]*client.ChainlinkK8sClient,
	error,
) {
	vars, err := ReadCRIBVars()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	config, err := tc.GetConfig([]string{"CRIB"}, tc.OCR)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if vars.Nodes < 2 {
		return nil, nil, nil, nil, fmt.Errorf("not enough chainlink nodes, need at least 2")
	}
	mockserverURL := fmt.Sprintf(mockserverCRIBTemplate, vars.Namespace, vars.IngressSuffix)
	var sethClient *seth.Client
	switch vars.Network {
	case "geth":
		netWSURL := fmt.Sprintf(ingressNetworkWSURLTemplate, vars.Namespace, vars.IngressSuffix)
		netHTTPURL := fmt.Sprintf(ingressNetworkHTTPURLTemplate, vars.Namespace, vars.IngressSuffix)
		setSethConfig(config, netWSURL, netHTTPURL)
		net := blockchain.EVMNetwork{
			Name:                 vars.Network,
			Simulated:            true,
			SupportsEIP1559:      true,
			ClientImplementation: blockchain.EthereumClientImplementation,
			ChainID:              1337,
			PrivateKeys: []string{
				"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			},
			URLs:                      []string{netWSURL},
			HTTPURLs:                  []string{netHTTPURL},
			ChainlinkTransactionLimit: 500000,
			Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
			MinimumConfirmations:      1,
			GasEstimationBuffer:       10000,
		}
		sethClient, err = seth_utils.GetChainClient(config, net)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	default:
		return nil, nil, nil, nil, errors.New("CRIB network is not supported")
	}
	// bootstrap node
	clClients := make([]*client.ChainlinkK8sClient, 0)
	c, err := client.NewChainlinkK8sClient(&client.ChainlinkConfig{
		URL:        fmt.Sprintf("https://%s-node%d%s", vars.Namespace, 1, vars.IngressSuffix),
		Email:      client.CLNodeTestEmail,
		InternalIP: fmt.Sprintf(internalNodeDNSTemplate, 1),
		Password:   client.CLNodeTestPassword,
	}, fmt.Sprintf(internalNodeDNSTemplate, 1), vars.Namespace)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	clClients = append(clClients, c)
	// all the other nodes, indices of nodes in CRIB starts with 1
	for i := 2; i <= vars.Nodes; i++ {
		cl, err := client.NewChainlinkK8sClient(&client.ChainlinkConfig{
			URL:        fmt.Sprintf("https://%s-node%d%s", vars.Namespace, i, vars.IngressSuffix),
			Email:      client.CLNodeTestEmail,
			InternalIP: fmt.Sprintf(internalNodeDNSTemplate, i),
			Password:   client.CLNodeTestPassword,
		}, fmt.Sprintf(internalNodeDNSTemplate, i), vars.Namespace)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		clClients = append(clClients, cl)
	}
	mockServerClient := msClient.NewMockserverClient(&msClient.MockserverConfig{
		LocalURL:   mockserverURL,
		ClusterURL: mockserverURL,
	})

	//nolint:gosec // G602 - false positive https://github.com/securego/gosec/issues/1005
	return sethClient, mockServerClient, clClients[0], clClients[1:], nil
}
