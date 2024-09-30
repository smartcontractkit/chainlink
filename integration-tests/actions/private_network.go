package actions

import (
	"github.com/rs/zerolog"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctf_config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"
	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
)

func EthereumNetworkConfigFromConfig(l zerolog.Logger, config ctf_config.GlobalTestConfig) (network ctf_test_env.EthereumNetwork, err error) {
	if config.GetPrivateEthereumNetworkConfig() == nil {
		l.Warn().Msg("No TOML private ethereum network config found, will use old geth")
		ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
		network, err = ethBuilder.
			WithEthereumVersion(ctf_config_types.EthereumVersion_Eth1).
			WithExecutionLayer(ctf_config_types.ExecutionLayer_Geth).
			Build()

		return
	}

	ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
	network, err = ethBuilder.
		WithExistingConfig(*config.GetPrivateEthereumNetworkConfig()).
		Build()

	return
}
