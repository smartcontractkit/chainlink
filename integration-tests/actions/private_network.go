package actions

import (
	"errors"

	"github.com/rs/zerolog"

	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func EthereumNetworkConfigFromConfig(l zerolog.Logger, config *tc.TestConfig) (network ctf_test_env.EthereumNetwork, err error) {
	ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
	network, err = ethBuilder.
		WithExistingConfig(*config.PrivateEthereumNetwork).
		Build()

	if errors.Is(err, ctf_test_env.ErrMissingExecClientEnvVar) {
		l.Warn().Msg("No exec client env var set, will use old geth")
		ethBuilder = ctf_test_env.NewEthereumNetworkBuilder()
		network, err = ethBuilder.
			WithConsensusType(ctf_test_env.ConsensusType_PoW).
			WithExecutionLayer(ctf_test_env.ExecutionLayer_Geth).
			Build()
	}

	return
}
