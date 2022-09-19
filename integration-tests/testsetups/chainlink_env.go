package testsetups

import (
	"os"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"golang.org/x/exp/maps"
)

// ChainlinkEnvFlags represents flags that will enable Chainlink features
type ChainlinkEnvFlags struct {
	PyroscopeEnabled bool
}

// NewChainlink wraps chainlink.New() allowing to enable Chainlink features by passing flags
func NewChainlink(index int, props map[string]interface{}, flags ChainlinkEnvFlags) environment.ConnectedChart {
	if props == nil {
		props = map[string]interface{}{}
	}
	if _, ok := props["env"]; !ok {
		props["env"] = map[string]interface{}{}
	}

	envOverrides := map[string]interface{}{}

	if flags.PyroscopeEnabled {
		pyroscopeOAuthToken := os.Getenv("PYROSCOPE_OAUTH_TOKEN")
		pyroscopeServerAddr := os.Getenv("PYROSCOPE_SERVER_ADDRESS")
		if pyroscopeServerAddr != "" {
			envOverrides["PYROSCOPE_OAUTH_TOKEN"] = pyroscopeOAuthToken
			envOverrides["PYROSCOPE_SERVER_ADDRESS"] = pyroscopeServerAddr
		}
	}

	if _, ok := props["env"].(map[string]interface{}); !ok {
		panic("failed to set env overrides")
		return nil
	}

	maps.Copy(props["env"].(map[string]interface{}), envOverrides)

	return chainlink.New(index, props)
}

// NewChainlinkWithPyroscope wraps NewChainlink() enabling Pyroscope
func NewChainlinkWithPyroscope(index int, props map[string]interface{}) environment.ConnectedChart {
	return NewChainlink(index, props, ChainlinkEnvFlags{PyroscopeEnabled: true})
}
