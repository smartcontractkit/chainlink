package k8s

import (
	"fmt"
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	client2 "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

const (
	DefaultConfigFilePath        = "../../../charts/chainlink-cluster/connect.toml"
	ErrReadConnectionConfig      = "failed to read TOML environment connection config"
	ErrUnmarshalConnectionConfig = "failed to unmarshal TOML environment connection config"
)

type ConnectionVars struct {
	Namespace                       string `toml:"namespace"`
	NetworkName                     string `toml:"network_name"`
	NetworkChainID                  int64  `toml:"network_chain_id"`
	NetworkPrivateKey               string `toml:"network_private_key"`
	NetworkWSURL                    string `toml:"network_ws_url"`
	NetworkHTTPURL                  string `toml:"network_http_url"`
	CLNodesNum                      int    `toml:"cl_nodes_num"`
	CLNodeURLTemplate               string `toml:"cl_node_url_template"`
	CLNodeInternalDNSRecordTemplate string `toml:"cl_node_internal_dns_record_template"`
	CLNodeUser                      string `toml:"cl_node_user"`
	CLNodePassword                  string `toml:"cl_node_password"`
	MockServerURL                   string `toml:"mockserver_url"`
}

// ConnectRemote connects to a local environment, see charts/chainlink-cluster
func ConnectRemote(l zerolog.Logger) (blockchain.EVMClient, *client2.MockserverClient, contracts.ContractDeployer, *client.ChainlinkK8sClient, []*client.ChainlinkK8sClient, error) {
	cfg, err := ReadConfig()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	net := &blockchain.EVMNetwork{
		Name:                 cfg.NetworkName,
		Simulated:            true,
		SupportsEIP1559:      true,
		ClientImplementation: blockchain.EthereumClientImplementation,
		ChainID:              1337,
		PrivateKeys: []string{
			cfg.NetworkPrivateKey,
		},
		URLs:                      []string{cfg.NetworkWSURL},
		HTTPURLs:                  []string{cfg.NetworkHTTPURL},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}
	cc, err := blockchain.NewEVMClientFromNetwork(*net, l)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	cd, err := contracts.NewContractDeployer(cc, l)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	clClients := make([]*client.ChainlinkK8sClient, 0)
	for i := 1; i <= cfg.CLNodesNum; i++ {
		c, err := client.NewChainlinkK8sClient(&client.ChainlinkConfig{
			URL:        fmt.Sprintf(cfg.CLNodeURLTemplate, i),
			Email:      cfg.CLNodeUser,
			InternalIP: fmt.Sprintf(cfg.CLNodeInternalDNSRecordTemplate, i),
			Password:   cfg.CLNodePassword,
		}, fmt.Sprintf(cfg.CLNodeInternalDNSRecordTemplate, i), cfg.Namespace)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
		clClients = append(clClients, c)
	}
	msClient := client2.NewMockserverClient(&client2.MockserverConfig{
		LocalURL:   cfg.MockServerURL,
		ClusterURL: cfg.MockServerURL,
	})
	return cc, msClient, cd, clClients[0], clClients[1:], nil
}

func ReadConfig() (*ConnectionVars, error) {
	var cfg *ConnectionVars
	var d []byte
	var err error
	d, err = os.ReadFile(DefaultConfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("%s, err: %w", ErrReadConnectionConfig, err)
	}
	err = toml.Unmarshal(d, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%s, err: %w", ErrUnmarshalConnectionConfig, err)
	}
	log.Info().Interface("Config", cfg).Msg("Connecting to environment from config")
	return cfg, nil
}
