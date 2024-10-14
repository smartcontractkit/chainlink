package test_env

import (
	"encoding/json"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"

	env "github.com/smartcontractkit/chainlink/integration-tests/types/envcommon"
)

type TestEnvConfig struct {
	Networks              []string                          `json:"networks"`
	Geth                  GethConfig                        `json:"geth"`
	MockAdapter           MockAdapterConfig                 `json:"mock_adapter"`
	ClCluster             *ClCluster                        `json:"clCluster"`
	EthereumNetworkConfig *ctf_config.EthereumNetworkConfig `json:"private_ethereum_config"`
}

type MockAdapterConfig struct {
	ContainerName string `json:"container_name"`
	ImpostersPath string `json:"imposters_path"`
}

type GethConfig struct {
	ContainerName string `json:"container_name"`
}

func NewTestEnvConfigFromFile(path string) (*TestEnvConfig, error) {
	c := &TestEnvConfig{}
	err := env.ParseJSONFile(path, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *TestEnvConfig) Json() string {
	b, _ := json.Marshal(c)
	return string(b)
}
