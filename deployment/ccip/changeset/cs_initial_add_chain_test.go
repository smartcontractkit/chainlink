package changeset

import (
	"testing"
	"time"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestOverrideCCIPParams(t *testing.T) {
	params := DefaultOCRParams(chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector, nil, nil)
	overrides := CCIPOCRParams{
		ExecuteOffChainConfig: pluginconfig.ExecuteOffchainConfig{
			RelativeBoostPerWaitHour: 10,
		},
		CommitOffChainConfig: pluginconfig.CommitOffchainConfig{
			TokenPriceBatchWriteFrequency:     *config.MustNewDuration(1_000_000 * time.Hour),
			RemoteGasPriceBatchWriteFrequency: *config.MustNewDuration(1_000_000 * time.Hour),
		},
	}
	newParams, err := params.Override(overrides)
	require.NoError(t, err)
	require.Equal(t, overrides.ExecuteOffChainConfig.RelativeBoostPerWaitHour, newParams.ExecuteOffChainConfig.RelativeBoostPerWaitHour)
	require.Equal(t, overrides.CommitOffChainConfig.TokenPriceBatchWriteFrequency, newParams.CommitOffChainConfig.TokenPriceBatchWriteFrequency)
	require.Equal(t, overrides.CommitOffChainConfig.RemoteGasPriceBatchWriteFrequency, newParams.CommitOffChainConfig.RemoteGasPriceBatchWriteFrequency)
	require.Equal(t, params.OCRParameters, newParams.OCRParameters)
}
