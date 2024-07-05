package crib

import (
	"fmt"
	"os"
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

// ConnectRemote connects to a local environment, see https://github.com/smartcontractkit/crib/tree/main/core
// connects to default CRIB network if simulated = true
func ConnectRemote(simulated bool) (
	*seth.Client,
	*msClient.MockserverClient,
	*client.ChainlinkK8sClient,
	[]*client.ChainlinkK8sClient,
	error,
) {
	ingressSuffix := os.Getenv("K8S_STAGING_INGRESS_SUFFIX")
	if ingressSuffix == "" {
		return nil, nil, nil, nil, errors.New("K8S_STAGING_INGRESS_SUFFIX must be set to connect to k8s ingresses")
	}
	config, err := tc.GetConfig([]string{"CRIB"}, tc.OCR)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if config.CRIB.CLNodesNum < 2 {
		return nil, nil, nil, nil, fmt.Errorf("not enough chainlink nodes, need at least 2, TOML key: [CRIB.nodes]")
	}
	cfg := config.CRIB
	mockserverURL := fmt.Sprintf(mockserverCRIBTemplate, cfg.Namespace, ingressSuffix)
	var sethClient *seth.Client
	if simulated {
		netWSURL := fmt.Sprintf(ingressNetworkWSURLTemplate, cfg.Namespace, ingressSuffix)
		netHTTPURL := fmt.Sprintf(ingressNetworkHTTPURLTemplate, cfg.Namespace, ingressSuffix)
		setSethConfig(config, netWSURL, netHTTPURL)
		net := blockchain.EVMNetwork{
			Name:                 cfg.NetworkName,
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
	}
	// bootstrap node
	clClients := make([]*client.ChainlinkK8sClient, 0)
	c, err := client.NewChainlinkK8sClient(&client.ChainlinkConfig{
		URL:        fmt.Sprintf("https://%s-node%d%s", cfg.Namespace, 1, ingressSuffix),
		Email:      client.CLNodeTestEmail,
		InternalIP: fmt.Sprintf(internalNodeDNSTemplate, 1),
		Password:   client.CLNodeTestPassword,
	}, fmt.Sprintf(internalNodeDNSTemplate, 1), cfg.Namespace)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	clClients = append(clClients, c)
	// all the other nodes, indices of nodes in CRIB starts with 1
	for i := 2; i <= cfg.CLNodesNum; i++ {
		cl, err := client.NewChainlinkK8sClient(&client.ChainlinkConfig{
			URL:        fmt.Sprintf("https://%s-node%d%s", cfg.Namespace, i, ingressSuffix),
			Email:      client.CLNodeTestEmail,
			InternalIP: fmt.Sprintf(internalNodeDNSTemplate, i),
			Password:   client.CLNodeTestPassword,
		}, fmt.Sprintf(internalNodeDNSTemplate, i), cfg.Namespace)
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
