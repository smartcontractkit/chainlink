package actions

import (
	"errors"

	"github.com/rs/zerolog"

	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
)

func EthereumNetworkConfigFromEnvOrDefault(l zerolog.Logger) (network ctf_test_env.EthereumNetwork, err error) {
	chainConfig := ctf_test_env.EthereumChainConfig{
		SecondsPerSlot: 8,
		SlotsPerEpoch:  4,
	}

	ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
	network, err = ethBuilder.
		WithExecClientFromEnvVar().
		WithEthereumChainConfig(chainConfig).
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
