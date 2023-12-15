package actions

import (
	"github.com/rs/zerolog"

	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func EthereumNetworkConfigFromConfig(l zerolog.Logger, config *tc.TestConfig) (network ctf_test_env.EthereumNetwork, err error) {
	if config.PrivateEthereumNetwork == nil {
		l.Warn().Msg("No TOML private ethereum network config found, will use old geth")
		ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
		network, err = ethBuilder.
			WithConsensusType(ctf_test_env.ConsensusType_PoW).
			WithExecutionLayer(ctf_test_env.ExecutionLayer_Geth).
			Build()

		return
	}

	ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
	network, err = ethBuilder.
		WithExistingConfig(*config.PrivateEthereumNetwork).
		Build()

	return
}
