package ccip

import (
	"fmt"
	"math"
	"strconv"

	"github.com/AlekSi/pointer"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
)

const (
	E2E_JD_IMAGE              = "E2E_JD_IMAGE"
	E2E_JD_VERSION            = "E2E_JD_VERSION"
	E2E_JD_GRPC               = "E2E_JD_GRPC"
	E2E_JD_WSRPC              = "E2E_JD_WSRPC"
	DEFAULT_DB_NAME           = "JD_DB"
	DEFAULT_DB_VERSION        = "14.1"
	E2E_RMN_RAGEPROXY_IMAGE   = "E2E_RMN_RAGEPROXY_IMAGE"
	E2E_RMN_RAGEPROXY_VERSION = "E2E_RMN_RAGEPROXY_VERSION"
	E2E_RMN_AFN2PROXY_IMAGE   = "E2E_RMN_AFN2PROXY_IMAGE"
	E2E_RMN_AFN2PROXY_VERSION = "E2E_RMN_AFN2PROXY_VERSION"
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
	RMNConfig               RMNConfig                                   `toml:",omitempty"`
	Load                    *LoadConfig                                 `toml:",omitempty"`
}

type RMNConfig struct {
	NoOfNodes    *int    `toml:",omitempty"`
	ProxyImage   *string `toml:",omitempty"`
	ProxyVersion *string `toml:",omitempty"`
	AFNImage     *string `toml:",omitempty"`
	AFNVersion   *string `toml:",omitempty"`
}

func (r *RMNConfig) GetProxyImage() string {
	image := pointer.GetString(r.ProxyImage)
	if image == "" {
		return ctfconfig.MustReadEnvVar_String(E2E_RMN_RAGEPROXY_IMAGE)
	}
	return image
}

func (r *RMNConfig) GetProxyVersion() string {
	version := pointer.GetString(r.ProxyVersion)
	if version == "" {
		return ctfconfig.MustReadEnvVar_String(E2E_RMN_RAGEPROXY_VERSION)
	}
	return version
}

func (r *RMNConfig) GetAFN2ProxyImage() string {
	image := pointer.GetString(r.AFNImage)
	if image == "" {
		return ctfconfig.MustReadEnvVar_String(E2E_RMN_AFN2PROXY_IMAGE)
	}
	return image
}

func (r *RMNConfig) GetAFN2ProxyVersion() string {
	version := pointer.GetString(r.AFNVersion)
	if version == "" {
		return ctfconfig.MustReadEnvVar_String(E2E_RMN_AFN2PROXY_VERSION)
	}
	return version
}

type NodeConfig struct {
	NoOfPluginNodes *int                        `toml:",omitempty"`
	NoOfBootstraps  *int                        `toml:",omitempty"`
	ClientConfig    *nodeclient.ChainlinkConfig `toml:",omitempty"`
}

type JDConfig struct {
	Image     *string `toml:",omitempty"`
	Version   *string `toml:",omitempty"`
	DBName    *string `toml:",omitempty"`
	DBVersion *string `toml:",omitempty"`
	JDGRPC    *string `toml:",omitempty"`
	JDWSRPC   *string `toml:",omitempty"`
}

// TODO: include all JD specific input in generic secret handling
func (o *JDConfig) GetJDGRPC() string {
	grpc := pointer.GetString(o.JDGRPC)
	if grpc == "" {
		return ctfconfig.MustReadEnvVar_String(E2E_JD_GRPC)
	}
	return grpc
}

func (o *JDConfig) GetJDWSRPC() string {
	wsrpc := pointer.GetString(o.JDWSRPC)
	if wsrpc == "" {
		return ctfconfig.MustReadEnvVar_String(E2E_JD_WSRPC)
	}
	return wsrpc
}

func (o *JDConfig) GetJDImage() string {
	image := pointer.GetString(o.Image)
	if image == "" {
		return ctfconfig.MustReadEnvVar_String(E2E_JD_IMAGE)
	}
	return image
}

func (o *JDConfig) GetJDVersion() string {
	version := pointer.GetString(o.Version)
	if version == "" {
		return ctfconfig.MustReadEnvVar_String(E2E_JD_VERSION)
	}
	return version
}

func (o *JDConfig) GetJDDBName() string {
	dbname := pointer.GetString(o.DBName)
	if dbname == "" {
		return DEFAULT_DB_NAME
	}
	return dbname
}

func (o *JDConfig) GetJDDBVersion() string {
	dbversion := pointer.GetString(o.DBVersion)
	if dbversion == "" {
		return DEFAULT_DB_VERSION
	}
	return dbversion
}

func (o *Config) Validate() error {
	var chainIds []int64
	for _, net := range o.PrivateEthereumNetworks {
		if net.EthereumChainConfig.ChainID < 0 {
			return fmt.Errorf("negative chain ID found for network %d", net.EthereumChainConfig.ChainID)
		}
		chainIds = append(chainIds, int64(net.EthereumChainConfig.ChainID))
	}
	homeChainSelector, err := strconv.ParseUint(pointer.GetString(o.HomeChainSelector), 10, 64)
	if err != nil {
		return err
	}
	isValid, err := IsSelectorValid(homeChainSelector, chainIds)
	if err != nil {
		return err
	}
	if !isValid {
		return ErrInvalidHomeChainSelector
	}
	feedChainSelector, err := strconv.ParseUint(pointer.GetString(o.FeedChainSelector), 10, 64)
	if err != nil {
		return err
	}
	isValid, err = IsSelectorValid(feedChainSelector, chainIds)
	if err != nil {
		return err
	}
	if !isValid {
		return ErrInvalidFeedChainSelector
	}
	return nil
}

func (o *Config) GetHomeChainSelector() uint64 {
	selector, _ := strconv.ParseUint(pointer.GetString(o.HomeChainSelector), 10, 64)
	return selector
}

func (o *Config) GetFeedChainSelector() uint64 {
	selector, _ := strconv.ParseUint(pointer.GetString(o.FeedChainSelector), 10, 64)
	return selector
}

func IsSelectorValid(selector uint64, chainIds []int64) (bool, error) {
	chainId, err := chainselectors.ChainIdFromSelector(selector)
	if err != nil {
		return false, err
	}

	for _, cID := range chainIds {
		if isEqualUint64AndInt64(chainId, cID) {
			return true, nil
		}
	}
	return false, nil
}

func isEqualUint64AndInt64(u uint64, i int64) bool {
	if i < 0 {
		return false // uint64 cannot be equal to a negative int64
	}
	if u > math.MaxInt64 {
		return false // uint64 cannot be equal to an int64 if it exceeds the maximum int64 value
	}
	return u == uint64(i)
}
