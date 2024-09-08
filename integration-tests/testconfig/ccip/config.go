package ccip

import (
	"github.com/AlekSi/pointer"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
)

type Config struct {
	PrivateEthereumNetworks map[string]*ctfconfig.EthereumNetworkConfig `toml:",omitempty"`
	CLNode                  *NodeConfig                                 `toml:",omitempty"`
	JobDistributorConfig    *JDConfig                                   `toml:",omitempty"`
}

type NodeConfig struct {
	NoOfPluginNodes *int `toml:",omitempty"`
	NoOfBootstraps  *int `toml:",omitempty"`
}

type JDConfig struct {
	JDImage *string `toml:",omitempty"`
}

func (o *Config) Validate() error {
	return nil
}

func (o *Config) GetJDImage() string {
	return pointer.GetString(o.JobDistributorConfig.JDImage)
}
