package ccip

import (
	"fmt"
	"strconv"

	"github.com/AlekSi/pointer"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
)

const (
	E2E_JD_IMAGE       = "E2E_JD_IMAGE"
	E2E_JD_VERSION     = "E2E_JD_VERSION"
	E2E_JD_GRPC        = "E2E_JD_GRPC"
	E2E_JD_WSRPC       = "E2E_JD_WSRPC"
	DEFAULT_DB_NAME    = "JD_DB"
	DEFAULT_DB_VERSION = "14.1"
)

var (
	ErrInvalidHomeChainSelector = fmt.Errorf("invalid home chain selector")
	ErrInvalidFeedChainSelector = fmt.Errorf("invalid feed chain selector")
)

type Config struct {
	PrivateEthereumNetworks map[string]*ctfconfig.EthereumNetworkConfig `toml:",omitempty"`
	CLNode                  *NodeConfig                                 `toml:",omitempty"`
	JobDistributorConfig    JDConfig                                    `toml:",omitempty"`
	HomeChainSelector       *string                                     `toml:",omitempty"`
	FeedChainSelector       *string                                     `toml:",omitempty"`
}

type NodeConfig struct {
	NoOfPluginNodes *int                    `toml:",omitempty"`
	NoOfBootstraps  *int                    `toml:",omitempty"`
	ClientConfig    *client.ChainlinkConfig `toml:",omitempty"`
}

type JDConfig struct {
	Image     *string `toml:",omitempty"`
	Version   *string `toml:",omitempty"`
	DBName    *string `toml:",omitempty"`
	DBVersion *string `toml:",omitempty"`
	JDGRPC    *string `toml:",omitempty"`
	JDWSRPC   *string `toml:",omitempty"`
}

func (o *Config) Validate() error {
	return nil
}

// TODO: include all JD specific input in generic secret handling
func (o *Config) GetJDGRPC() string {
	grpc := pointer.GetString(o.JobDistributorConfig.JDGRPC)
	if grpc == "" {
		return ctfconfig.MustReadEnvVar_String(E2E_JD_GRPC)
	}
	return grpc
}

func (o *Config) GetJDWSRPC() string {
	wsrpc := pointer.GetString(o.JobDistributorConfig.JDWSRPC)
	if wsrpc == "" {
		return ctfconfig.MustReadEnvVar_String(E2E_JD_WSRPC)
	}
	return wsrpc
}

func (o *Config) GetJDImage() string {
	image := pointer.GetString(o.JobDistributorConfig.Image)
	if image == "" {
		return ctfconfig.MustReadEnvVar_String(E2E_JD_IMAGE)
	}
	return image
}

func (o *Config) GetJDVersion() string {
	version := pointer.GetString(o.JobDistributorConfig.Version)
	if version == "" {
		return ctfconfig.MustReadEnvVar_String(E2E_JD_VERSION)
	}
	return version
}

func (o *Config) GetJDDBName() string {
	dbname := pointer.GetString(o.JobDistributorConfig.DBName)
	if dbname == "" {
		return DEFAULT_DB_NAME
	}
	return dbname
}

func (o *Config) GetJDDBVersion() string {
	dbversion := pointer.GetString(o.JobDistributorConfig.DBVersion)
	if dbversion == "" {
		return DEFAULT_DB_VERSION
	}
	return dbversion
}

func (o *Config) GetHomeChainSelector(evmNetworks []blockchain.EVMNetwork) (uint64, error) {
	homeChainSelector, err := strconv.ParseUint(pointer.GetString(o.HomeChainSelector), 10, 64)
	if err != nil {
		return 0, err
	}
	isValid, err := IsSelectorValid(homeChainSelector, evmNetworks)
	if err != nil {
		return 0, err
	}
	if !isValid {
		return 0, ErrInvalidHomeChainSelector
	}
	return homeChainSelector, nil
}

func (o *Config) GetFeedChainSelector(evmNetworks []blockchain.EVMNetwork) (uint64, error) {
	feedChainSelector, err := strconv.ParseUint(pointer.GetString(o.FeedChainSelector), 10, 64)
	if err != nil {
		return 0, err
	}
	isValid, err := IsSelectorValid(feedChainSelector, evmNetworks)
	if err != nil {
		return 0, err
	}
	if !isValid {
		return 0, ErrInvalidFeedChainSelector
	}
	return feedChainSelector, nil
}

func IsSelectorValid(selector uint64, evmNetworks []blockchain.EVMNetwork) (bool, error) {
	chainId, err := chainselectors.ChainIdFromSelector(selector)
	if err != nil {
		return false, err
	}
	for _, net := range evmNetworks {
		if net.ChainID == int64(chainId) {
			return true, nil
		}
	}
	return false, nil
}
