package v1_2_0

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
)

func TestCommitReportEncoding(t *testing.T) {
	t.Parallel()
	report := cciptypes.CommitStoreReport{
		TokenPrices: []cciptypes.TokenPrice{
			{
				Token: cciptypes.Address(utils.RandomAddress().String()),
				Value: big.NewInt(9e18),
			},
			{
				Token: cciptypes.Address(utils.RandomAddress().String()),
				Value: big.NewInt(1e18),
			},
		},
		GasPrices: []cciptypes.GasPrice{
			{
				DestChainSelector: rand.Uint64(),
				Value:             big.NewInt(2000e9),
			},
			{
				DestChainSelector: rand.Uint64(),
				Value:             big.NewInt(3000e9),
			},
		},
		MerkleRoot: [32]byte{123},
		Interval:   cciptypes.CommitStoreInterval{Min: 1, Max: 10},
	}

	c, err := NewCommitStore(logger.TestLogger(t), utils.RandomAddress(), nil, mocks.NewLogPoller(t), nil)
	assert.NoError(t, err)

	encodedReport, err := c.EncodeCommitReport(report)
	require.NoError(t, err)
	assert.Greater(t, len(encodedReport), 0)

	decodedReport, err := c.DecodeCommitReport(encodedReport)
	require.NoError(t, err)
	require.Equal(t, report, decodedReport)
}

func TestCommitStoreV120ffchainConfigEncoding(t *testing.T) {
	t.Parallel()
	validConfig := JSONCommitOffchainConfig{
		SourceFinalityDepth:      3,
		DestFinalityDepth:        4,
		MaxGasPrice:              200e9,
		GasPriceHeartBeat:        *config.MustNewDuration(1 * time.Minute),
		DAGasPriceDeviationPPB:   10,
		ExecGasPriceDeviationPPB: 11,
		TokenPriceHeartBeat:      *config.MustNewDuration(2 * time.Minute),
		TokenPriceDeviationPPB:   12,
		InflightCacheExpiry:      *config.MustNewDuration(3 * time.Minute),
	}

	require.NoError(t, validConfig.Validate())

	tests := []struct {
		name       string
		want       JSONCommitOffchainConfig
		errPattern string
	}{
		{
			name: "legacy offchain config format parses",
			want: validConfig,
		},
		{
			name: "can omit finality depth",
			want: modifyCopy(validConfig, func(c *JSONCommitOffchainConfig) {
				c.SourceFinalityDepth = 0
				c.DestFinalityDepth = 0
			}),
		},
		{
			name: "can set the SourceMaxGasPrice",
			want: modifyCopy(validConfig, func(c *JSONCommitOffchainConfig) {
				c.MaxGasPrice = 0
				c.SourceMaxGasPrice = 200e9
			}),
		},
		{
			name: "must set SourceMaxGasPrice",
			want: modifyCopy(validConfig, func(c *JSONCommitOffchainConfig) {
				c.MaxGasPrice = 0
				c.SourceMaxGasPrice = 0
			}),
			errPattern: "SourceMaxGasPrice",
		},
		{
			name: "cannot set both MaxGasPrice and SourceMaxGasPrice",
			want: modifyCopy(validConfig, func(c *JSONCommitOffchainConfig) {
				c.SourceMaxGasPrice = c.MaxGasPrice
			}),
			errPattern: "MaxGasPrice and SourceMaxGasPrice",
		},
		{
			name: "must set GasPriceHeartBeat",
			want: modifyCopy(validConfig, func(c *JSONCommitOffchainConfig) {
				c.GasPriceHeartBeat = *config.MustNewDuration(0)
			}),
			errPattern: "GasPriceHeartBeat",
		},
		{
			name: "must set ExecGasPriceDeviationPPB",
			want: modifyCopy(validConfig, func(c *JSONCommitOffchainConfig) {
				c.ExecGasPriceDeviationPPB = 0
			}),
			errPattern: "ExecGasPriceDeviationPPB",
		},
		{
			name: "must set TokenPriceHeartBeat",
			want: modifyCopy(validConfig, func(c *JSONCommitOffchainConfig) {
				c.TokenPriceHeartBeat = *config.MustNewDuration(0)
			}),
			errPattern: "TokenPriceHeartBeat",
		},
		{
			name: "must set TokenPriceDeviationPPB",
			want: modifyCopy(validConfig, func(c *JSONCommitOffchainConfig) {
				c.TokenPriceDeviationPPB = 0
			}),
			errPattern: "TokenPriceDeviationPPB",
		},
		{
			name: "must set InflightCacheExpiry",
			want: modifyCopy(validConfig, func(c *JSONCommitOffchainConfig) {
				c.InflightCacheExpiry = *config.MustNewDuration(0)
			}),
			errPattern: "InflightCacheExpiry",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			exp := tc.want
			encode, err := ccipconfig.EncodeOffchainConfig(&exp)
			require.NoError(t, err)
			got, err := ccipconfig.DecodeOffchainConfig[JSONCommitOffchainConfig](encode)

			if tc.errPattern != "" {
				require.ErrorContains(t, err, tc.errPattern)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}

func TestCommitStoreV120ComputesGasPrice(t *testing.T) {
	validConfig := JSONCommitOffchainConfig{
		SourceFinalityDepth:      3,
		DestFinalityDepth:        4,
		MaxGasPrice:              200e9,
		GasPriceHeartBeat:        *config.MustNewDuration(1 * time.Minute),
		DAGasPriceDeviationPPB:   10,
		ExecGasPriceDeviationPPB: 11,
		TokenPriceHeartBeat:      *config.MustNewDuration(2 * time.Minute),
		TokenPriceDeviationPPB:   12,
		InflightCacheExpiry:      *config.MustNewDuration(3 * time.Minute),
	}

	require.NoError(t, validConfig.Validate())
	require.Equal(t, uint64(200e9), validConfig.ComputeSourceMaxGasPrice())

	validConfig.MaxGasPrice = 0
	validConfig.SourceMaxGasPrice = 250e9
	require.NoError(t, validConfig.Validate())
	require.Equal(t, uint64(250e9), validConfig.ComputeSourceMaxGasPrice())
}
